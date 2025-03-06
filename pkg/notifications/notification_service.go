/*
 * Copyright 2025 Carver Automation Corporation.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package notifications

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"log"
	"sync"
	"time"
)

// Service implements the NotificationService interface.
type Service struct {
	store     NotificationStore
	handlers  map[TargetType]NotificationHandler
	handlerMu sync.RWMutex
}

// NewService creates a new notification service.
func NewService(store NotificationStore) *Service {
	return &Service{
		store:    store,
		handlers: make(map[TargetType]NotificationHandler),
	}
}

// RegisterHandler registers a notification handler for a specific target type.
func (s *Service) RegisterHandler(targetType TargetType, handler NotificationHandler) {
	s.handlerMu.Lock()
	defer s.handlerMu.Unlock()
	s.handlers[targetType] = handler
}

// GetHandler retrieves a notification handler for a specific target type.
func (s *Service) getHandler(targetType TargetType) (NotificationHandler, error) {
	s.handlerMu.RLock()
	defer s.handlerMu.RUnlock()

	handler, ok := s.handlers[targetType]
	if !ok {
		return nil, fmt.Errorf("no handler registered for target type: %s", targetType)
	}
	return handler, nil
}

// CreateNotification creates a new notification.
func (s *Service) CreateNotification(ctx context.Context, req NotificationRequest) (*Notification, error) {
	if req.AlertID == "" || req.Title == "" || req.Message == "" {
		return nil, fmt.Errorf("%w: missing required fields", ErrInvalidRequest)
	}

	// Start a transaction
	tx, err := s.store.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Printf("Error rolling back transaction: %v", rbErr)
			}
		}
	}()

	now := time.Now()
	notification := &Notification{
		AlertID:     req.AlertID,
		NodeID:      req.NodeID,
		ServiceName: req.ServiceName,
		Level:       req.Level,
		Title:       req.Title,
		Message:     req.Message,
		Status:      StatusPending,
		CreatedAt:   now,
		UpdatedAt:   now,
		Metadata:    req.Metadata,
	}

	// Set expiration time if provided
	if req.ExpireAfter != nil {
		expireAt := now.Add(*req.ExpireAfter)
		notification.ExpireAt = &expireAt
	}

	// Insert the notification
	notificationID, err := s.store.CreateNotification(ctx, tx, notification)
	if err != nil {
		return nil, fmt.Errorf("failed to create notification: %w", err)
	}
	notification.ID = notificationID

	// Add targets
	for _, targetReq := range req.Targets {
		target := &NotificationTarget{
			NotificationID: notificationID,
			TargetType:     targetReq.TargetType,
			TargetID:       targetReq.TargetID,
			Status:         TargetStatusPending,
		}

		targetID, err := s.store.AddNotificationTarget(ctx, tx, target)
		if err != nil {
			return nil, fmt.Errorf("failed to add target: %w", err)
		}
		target.ID = targetID
		notification.Targets = append(notification.Targets, *target)
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Send notifications to targets in background
	go s.processNotificationTargets(context.Background(), notification)

	return notification, nil
}

// processNotificationTargets sends the notification to all targets.
func (s *Service) processNotificationTargets(ctx context.Context, notification *Notification) {
	for i, target := range notification.Targets {
		// Get the handler for this target type
		handler, err := s.getHandler(target.TargetType)
		if err != nil {
			log.Printf("Error getting handler for target type %s: %v", target.TargetType, err)
			s.updateTargetStatus(ctx, target.ID, TargetStatusFailed, "", fmt.Sprintf("No handler for target type: %s", target.TargetType))
			continue
		}

		// Send the notification
		err = handler.SendNotification(ctx, notification, &target)
		if err != nil {
			log.Printf("Error sending notification to target %s: %v", target.TargetID, err)
			s.updateTargetStatus(ctx, target.ID, TargetStatusFailed, "", fmt.Sprintf("Failed to send: %v", err))
			continue
		}

		// Update target status to sent
		s.updateTargetStatus(ctx, target.ID, TargetStatusSent, target.ExternalID, target.ResponseData)

		// Update our local copy
		notification.Targets[i] = target
	}

	// Update notification status if all targets are processed
	allProcessed := true
	anySucceeded := false
	for _, target := range notification.Targets {
		if target.Status == TargetStatusPending {
			allProcessed = false
		}
		if target.Status == TargetStatusSent {
			anySucceeded = true
		}
	}

	if allProcessed {
		status := StatusSent
		if !anySucceeded {
			status = StatusPending // Keep pending if all sends failed
			log.Printf("All notification targets failed for notification %d", notification.ID)
		}

		tx, err := s.store.Begin()
		if err != nil {
			log.Printf("Error beginning transaction for notification status update: %v", err)
			return
		}

		if err := s.store.UpdateNotificationStatus(ctx, tx, notification.ID, status); err != nil {
			log.Printf("Error updating notification status: %v", err)
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Printf("Error rolling back transaction: %v", rbErr)
			}
			return
		}

		if err := tx.Commit(); err != nil {
			log.Printf("Error committing transaction: %v", err)
		}
	}
}

// updateTargetStatus updates the status of a notification target.
func (s *Service) updateTargetStatus(ctx context.Context, targetID int64, status TargetStatus, externalID, responseData string) {
	tx, err := s.store.Begin()
	if err != nil {
		log.Printf("Error beginning transaction for target status update: %v", err)
		return
	}

	if err := s.store.UpdateTargetStatus(ctx, tx, targetID, status, externalID, responseData); err != nil {
		log.Printf("Error updating target status: %v", err)
		if rbErr := tx.Rollback(); rbErr != nil {
			log.Printf("Error rolling back transaction: %v", rbErr)
		}
		return
	}

	if err := tx.Commit(); err != nil {
		log.Printf("Error committing transaction: %v", err)
	}
}

// GetNotification retrieves a notification by ID.
func (s *Service) GetNotification(ctx context.Context, id int64) (*Notification, error) {
	notification, err := s.store.GetNotification(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get notification: %w", err)
	}
	return notification, nil
}

// ListNotifications retrieves notifications based on filter criteria.
func (s *Service) ListNotifications(ctx context.Context, filter NotificationFilter) ([]Notification, error) {
	notifications, err := s.store.ListNotifications(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to list notifications: %w", err)
	}
	return notifications, nil
}

// AcknowledgeNotification acknowledges a notification.
func (s *Service) AcknowledgeNotification(ctx context.Context, req AcknowledgeRequest) error {
	if req.NotificationID == 0 || req.AcknowledgedBy == "" {
		return fmt.Errorf("%w: missing required fields", ErrInvalidRequest)
	}

	// Start a transaction
	tx, err := s.store.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Printf("Error rolling back transaction: %v", rbErr)
			}
		}
	}()

	// Get the notification first to check its status
	notification, err := s.store.GetNotification(ctx, req.NotificationID)
	if err != nil {
		return fmt.Errorf("failed to get notification: %w", err)
	}

	// Check if already acknowledged
	if notification.Status == StatusAcknowledged {
		return ErrAlreadyAcknowledged
	}

	// Create the acknowledgment
	ack := &Acknowledgment{
		NotificationID: req.NotificationID,
		TargetID:       req.TargetID,
		AcknowledgedBy: req.AcknowledgedBy,
		AcknowledgedAt: time.Now(),
		Method:         req.Method,
		Comment:        req.Comment,
	}

	// Insert the acknowledgment
	_, err = s.store.CreateAcknowledgment(ctx, tx, ack)
	if err != nil {
		return fmt.Errorf("failed to create acknowledgment: %w", err)
	}

	// Update notification status
	err = s.store.UpdateNotificationStatus(ctx, tx, req.NotificationID, StatusAcknowledged)
	if err != nil {
		return fmt.Errorf("failed to update notification status: %w", err)
	}

	// If a target is specified, update its status too
	if req.TargetID != nil {
		err = s.store.UpdateTargetStatus(ctx, tx, *req.TargetID, TargetStatusAcknowledged, "", "")
		if err != nil {
			return fmt.Errorf("failed to update target status: %w", err)
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// ResolveNotification marks a notification as resolved.
func (s *Service) ResolveNotification(ctx context.Context, id int64) error {
	tx, err := s.store.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Printf("Error rolling back transaction: %v", rbErr)
			}
		}
	}()

	// Get the notification first to check it exists
	_, err = s.store.GetNotification(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get notification: %w", err)
	}

	// Update notification status
	err = s.store.UpdateNotificationStatus(ctx, tx, id, StatusResolved)
	if err != nil {
		return fmt.Errorf("failed to update notification status: %w", err)
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// DeleteNotification deletes a notification.
func (s *Service) DeleteNotification(ctx context.Context, id int64) error {
	return s.store.DeleteNotification(ctx, id)
}

// CleanupExpiredNotifications removes or updates expired notifications.
func (s *Service) CleanupExpiredNotifications(ctx context.Context) (int, error) {
	expiredIDs, err := s.store.GetExpiredNotifications(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get expired notifications: %w", err)
	}

	count := 0
	for _, id := range expiredIDs {
		tx, err := s.store.Begin()
		if err != nil {
			log.Printf("Error beginning transaction for expired notification %d: %v", id, err)
			continue
		}

		err = s.store.UpdateNotificationStatus(ctx, tx, id, StatusExpired)
		if err != nil {
			log.Printf("Error updating status for expired notification %d: %v", id, err)
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Printf("Error rolling back transaction: %v", rbErr)
			}
			continue
		}

		if err := tx.Commit(); err != nil {
			log.Printf("Error committing transaction for expired notification %d: %v", id, err)
			continue
		}

		count++
	}

	return count, nil
}

// CreateAPIKey creates a new API key for a notification service.
func (s *Service) CreateAPIKey(ctx context.Context, serviceName string, permissions []string, expires *time.Time) (*APIKey, string, error) {
	// Generate a random API key
	keyBytes := make([]byte, 32)
	if _, err := rand.Read(keyBytes); err != nil {
		return nil, "", fmt.Errorf("failed to generate random key: %w", err)
	}
	keyValue := base64.URLEncoding.EncodeToString(keyBytes)

	// Generate a random key ID
	idBytes := make([]byte, 16)
	if _, err := rand.Read(idBytes); err != nil {
		return nil, "", fmt.Errorf("failed to generate random ID: %w", err)
	}
	keyID := hex.EncodeToString(idBytes)

	// Create the API key record
	key := &APIKey{
		KeyID:       keyID,
		ServiceName: serviceName,
		CreatedAt:   time.Now(),
		ExpiresAt:   expires,
		Permissions: permissions,
	}

	// Store the API key
	err := s.store.CreateAPIKey(ctx, key, keyValue)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create API key: %w", err)
	}

	return key, keyValue, nil
}

// ValidateAPIKey validates an API key and returns its permissions.
func (s *Service) ValidateAPIKey(ctx context.Context, keyID, keyValue string) (*APIKey, error) {
	key, err := s.store.GetAPIKeyByID(ctx, keyID)
	if err != nil {
		return nil, ErrInvalidAPIKey
	}

	// Check if expired
	if key.ExpiresAt != nil && key.ExpiresAt.Before(time.Now()) {
		return nil, ErrAPIKeyExpired
	}

	// Verify the key hash
	expectedHash := sha256.Sum256([]byte(keyValue))
	expectedHashHex := hex.EncodeToString(expectedHash[:])

	if subtle.ConstantTimeCompare([]byte(key.KeyHash), []byte(expectedHashHex)) != 1 {
		return nil, ErrInvalidAPIKey
	}

	// Update last used timestamp
	go func() {
		if err := s.store.UpdateAPIKeyLastUsed(context.Background(), keyID); err != nil {
			log.Printf("Error updating API key last used: %v", err)
		}
	}()

	return key, nil
}

// RevokeAPIKey revokes an API key.
func (s *Service) RevokeAPIKey(ctx context.Context, keyID string) error {
	return s.store.DeleteAPIKey(ctx, keyID)
}
