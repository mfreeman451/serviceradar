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
	"log"
	"time"
)

// NewCleanupService creates a new notification cleanup service.
func NewCleanupService(service NotificationService, config CleanupConfig) *CleanupService {
	// Set default values if not provided
	if config.Interval == 0 {
		config.Interval = 1 * time.Hour
	}
	if config.ResolvedRetention == 0 {
		config.ResolvedRetention = 30 * 24 * time.Hour // 30 days
	}
	if config.AcknowledgedRetention == 0 {
		config.AcknowledgedRetention = 7 * 24 * time.Hour // 7 days
	}
	if config.PendingRetention == 0 {
		config.PendingRetention = 3 * 24 * time.Hour // 3 days
	}

	return &CleanupService{
		service: service,
		config:  config,
		stopCh:  make(chan struct{}),
	}
}

// Start starts the cleanup service.
func (s *CleanupService) Start() {
	go s.run()
}

// Stop stops the cleanup service.
func (s *CleanupService) Stop() {
	close(s.stopCh)
}

// run executes the cleanup loop.
func (s *CleanupService) run() {
	ticker := time.NewTicker(s.config.Interval)
	defer ticker.Stop()

	// Run cleanup immediately on start
	s.cleanup()

	for {
		select {
		case <-s.stopCh:
			log.Println("Notification cleanup service stopped")
			return
		case <-ticker.C:
			s.cleanup()
		}
	}
}

// cleanup performs the actual cleanup operations.
func (s *CleanupService) cleanup() {
	log.Println("Starting notification cleanup process")

	// Use a background context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// 1. Clean up expired notifications
	expiredCount, err := s.service.CleanupExpiredNotifications(ctx)
	if err != nil {
		log.Printf("Error cleaning up expired notifications: %v", err)
	} else {
		log.Printf("Marked %d expired notifications", expiredCount)
	}

	// 2. Clean up old notifications by status
	cleanupByStatus := func(status NotificationStatus, retention time.Duration) {
		since := time.Now().Add(-retention)
		filter := NotificationFilter{
			Status: &status,
			Until:  &since,
		}

		notifications, err := s.service.ListNotifications(ctx, filter)
		if err != nil {
			log.Printf("Error listing old %s notifications: %v", status, err)
			return
		}

		for _, notification := range notifications {
			if err := s.service.DeleteNotification(ctx, notification.ID); err != nil {
				log.Printf("Error deleting old notification %d: %v", notification.ID, err)
			}
		}

		log.Printf("Deleted %d old %s notifications", len(notifications), status)
	}

	// Clean up by status with appropriate retention periods
	cleanupByStatus(StatusResolved, s.config.ResolvedRetention)
	cleanupByStatus(StatusAcknowledged, s.config.AcknowledgedRetention)
	cleanupByStatus(StatusExpired, s.config.PendingRetention)

	// 3. Delete old pending notifications that have been around too long
	pendingSince := time.Now().Add(-s.config.PendingRetention)
	pendingFilter := NotificationFilter{
		Status: &StatusPending,
		Until:  &pendingSince,
	}

	pendingNotifications, err := s.service.ListNotifications(ctx, pendingFilter)
	if err != nil {
		log.Printf("Error listing old pending notifications: %v", err)
	} else {
		for _, notification := range pendingNotifications {
			if err := s.service.DeleteNotification(ctx, notification.ID); err != nil {
				log.Printf("Error deleting old pending notification %d: %v", notification.ID, err)
			}
		}
		log.Printf("Deleted %d old pending notifications", len(pendingNotifications))
	}

	log.Println("Notification cleanup process completed")
}
