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

// Package notifications provides types and interfaces for the notification system
package notifications

import (
	"context"
	"database/sql"
	"time"
)

// NotificationService defines the interface for the notification service.
type NotificationService interface {
	// CreateNotification creates a new notification.
	CreateNotification(ctx context.Context, req NotificationRequest) (*Notification, error)

	// GetNotification retrieves a notification by ID.
	GetNotification(ctx context.Context, id int64) (*Notification, error)

	// ListNotifications retrieves notifications based on filter criteria.
	ListNotifications(ctx context.Context, filter NotificationFilter) ([]Notification, error)

	// AcknowledgeNotification acknowledges a notification.
	AcknowledgeNotification(ctx context.Context, req AcknowledgeRequest) error

	// ResolveNotification marks a notification as resolved.
	ResolveNotification(ctx context.Context, id int64) error

	// DeleteNotification deletes a notification.
	DeleteNotification(ctx context.Context, id int64) error

	// CleanupExpiredNotifications removes or updates expired notifications.
	CleanupExpiredNotifications(ctx context.Context) (int, error)

	// CreateAPIKey creates a new API key for a notification service.
	CreateAPIKey(ctx context.Context, serviceName string, permissions []string, expires *time.Time) (*APIKey, string, error)

	// ValidateAPIKey validates an API key and returns its permissions.
	ValidateAPIKey(ctx context.Context, keyID, keyValue string) (*APIKey, error)

	// RevokeAPIKey revokes an API key.
	RevokeAPIKey(ctx context.Context, keyID string) error
}

// NotificationStore defines the interface for storage operations.
type NotificationStore interface {
	// Notification operations

	CreateNotification(ctx context.Context, tx *sql.Tx, notification *Notification) (int64, error)
	GetNotification(ctx context.Context, id int64) (*Notification, error)
	ListNotifications(ctx context.Context, filter NotificationFilter) ([]Notification, error)
	UpdateNotificationStatus(ctx context.Context, tx *sql.Tx, id int64, status NotificationStatus) error
	DeleteNotification(ctx context.Context, id int64) error

	// Target operations

	AddNotificationTarget(ctx context.Context, tx *sql.Tx, target *NotificationTarget) (int64, error)
	UpdateTargetStatus(ctx context.Context, tx *sql.Tx, id int64, status TargetStatus, externalID, responseData string) error
	GetTargetsForNotification(ctx context.Context, notificationID int64) ([]NotificationTarget, error)

	// Acknowledgment operations

	CreateAcknowledgment(ctx context.Context, tx *sql.Tx, ack *Acknowledgment) (int64, error)
	GetAcknowledgmentsForNotification(ctx context.Context, notificationID int64) ([]Acknowledgment, error)

	// API key operations

	CreateAPIKey(ctx context.Context, key *APIKey, keyValue string) error
	GetAPIKeyByID(ctx context.Context, keyID string) (*APIKey, error)
	UpdateAPIKeyLastUsed(ctx context.Context, keyID string) error
	DeleteAPIKey(ctx context.Context, keyID string) error

	// Maintenance operations

	GetExpiredNotifications(ctx context.Context) ([]int64, error)

	// Transaction helpers

	Begin() (*sql.Tx, error)
}

// NotificationHandler is the interface that notification targets must implement.
type NotificationHandler interface {
	// SendNotification sends a notification to the target.
	SendNotification(ctx context.Context, notification *Notification, target *NotificationTarget) error

	// ParseAcknowledgment parses an acknowledgment from the target.
	ParseAcknowledgment(ctx context.Context, data []byte) (*AcknowledgeRequest, error)
}
