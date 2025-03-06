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
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

// Store implements the NotificationStore interface with SQLite.
type Store struct {
	db *sql.DB
}

// NewStore creates a new notification store with the given database.
func NewStore(database *sql.DB) (*Store, error) {
	return &Store{
		db: database,
	}, nil
}

// Begin starts a new transaction.
func (s *Store) Begin() (*sql.Tx, error) {
	return s.db.Begin()
}

// serializeMetadata converts a map to a JSON string.
func serializeMetadata(metadata map[string]interface{}) (string, error) {
	if metadata == nil {
		return "", nil
	}
	data, err := json.Marshal(metadata)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// deserializeMetadata converts a JSON string to a map.
func deserializeMetadata(data string) (map[string]interface{}, error) {
	if data == "" {
		return nil, nil
	}
	var metadata map[string]interface{}
	err := json.Unmarshal([]byte(data), &metadata)
	if err != nil {
		return nil, err
	}
	return metadata, nil
}

// CreateNotification adds a new notification to the database.
func (s *Store) CreateNotification(ctx context.Context, tx *sql.Tx, notification *Notification) (int64, error) {
	metadataJSON, err := serializeMetadata(notification.Metadata)
	if err != nil {
		return 0, fmt.Errorf("failed to serialize metadata: %w", err)
	}

	query := `
		INSERT INTO notifications 
		(alert_id, node_id, service_name, level, title, message, status, created_at, updated_at, expire_at, metadata)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	var result sql.Result
	if tx != nil {
		result, err = tx.ExecContext(ctx, query,
			notification.AlertID,
			notification.NodeID,
			notification.ServiceName,
			notification.Level,
			notification.Title,
			notification.Message,
			notification.Status,
			notification.CreatedAt,
			notification.UpdatedAt,
			notification.ExpireAt,
			metadataJSON,
		)
	} else {
		result, err = s.db.ExecContext(ctx, query,
			notification.AlertID,
			notification.NodeID,
			notification.ServiceName,
			notification.Level,
			notification.Title,
			notification.Message,
			notification.Status,
			notification.CreatedAt,
			notification.UpdatedAt,
			notification.ExpireAt,
			metadataJSON,
		)
	}

	if err != nil {
		return 0, fmt.Errorf("failed to create notification: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert ID: %w", err)
	}

	return id, nil
}

// GetNotification retrieves a notification by ID.
func (s *Store) GetNotification(ctx context.Context, id int64) (*Notification, error) {
	query := `
		SELECT id, alert_id, node_id, service_name, level, title, message, status, 
			created_at, updated_at, expire_at, metadata
		FROM notifications
		WHERE id = ?
	`

	var notification Notification
	var metadataJSON sql.NullString
	var expireAt sql.NullTime

	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&notification.ID,
		&notification.AlertID,
		&notification.NodeID,
		&notification.ServiceName,
		&notification.Level,
		&notification.Title,
		&notification.Message,
		&notification.Status,
		&notification.CreatedAt,
		&notification.UpdatedAt,
		&expireAt,
		&metadataJSON,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("notification not found: %d", id)
		}
		return nil, fmt.Errorf("failed to get notification: %w", err)
	}

	if expireAt.Valid {
		notification.ExpireAt = &expireAt.Time
	}

	if metadataJSON.Valid {
		metadata, err := deserializeMetadata(metadataJSON.String)
		if err != nil {
			return nil, fmt.Errorf("failed to deserialize metadata: %w", err)
		}
		notification.Metadata = metadata
	}

	// Get targets
	targets, err := s.GetTargetsForNotification(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get targets: %w", err)
	}
	notification.Targets = targets

	// Get acknowledgments
	acks, err := s.GetAcknowledgmentsForNotification(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get acknowledgments: %w", err)
	}
	notification.Acknowledgments = acks

	return &notification, nil
}

// ListNotifications retrieves notifications based on filter criteria.
func (s *Store) ListNotifications(ctx context.Context, filter NotificationFilter) ([]Notification, error) {
	var conditions []string
	var args []interface{}

	if filter.AlertID != "" {
		conditions = append(conditions, "alert_id = ?")
		args = append(args, filter.AlertID)
	}

	if filter.NodeID != "" {
		conditions = append(conditions, "node_id = ?")
		args = append(args, filter.NodeID)
	}

	if filter.ServiceName != "" {
		conditions = append(conditions, "service_name = ?")
		args = append(args, filter.ServiceName)
	}

	if filter.Level != nil {
		conditions = append(conditions, "level = ?")
		args = append(args, *filter.Level)
	}

	if filter.Status != nil {
		conditions = append(conditions, "status = ?")
		args = append(args, *filter.Status)
	}

	if filter.Acknowledged != nil {
		if *filter.Acknowledged {
			conditions = append(conditions, "status = ?")
			args = append(args, StatusAcknowledged)
		} else {
			conditions = append(conditions, "status != ?")
			args = append(args, StatusAcknowledged)
		}
	}

	if filter.Since != nil {
		conditions = append(conditions, "created_at >= ?")
		args = append(args, *filter.Since)
	}

	if filter.Until != nil {
		conditions = append(conditions, "created_at <= ?")
		args = append(args, *filter.Until)
	}

	query := "SELECT id, alert_id, node_id, service_name, level, title, message, status, created_at, updated_at, expire_at, metadata FROM notifications"
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}
	query += " ORDER BY created_at DESC"

	if filter.Limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", filter.Limit)
		if filter.Offset > 0 {
			query += fmt.Sprintf(" OFFSET %d", filter.Offset)
		}
	}

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query notifications: %w", err)
	}
	defer rows.Close()

	var notifications []Notification
	for rows.Next() {
		var notification Notification
		var metadataJSON sql.NullString
		var expireAt sql.NullTime

		err := rows.Scan(
			&notification.ID,
			&notification.AlertID,
			&notification.NodeID,
			&notification.ServiceName,
			&notification.Level,
			&notification.Title,
			&notification.Message,
			&notification.Status,
			&notification.CreatedAt,
			&notification.UpdatedAt,
			&expireAt,
			&metadataJSON,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan notification: %w", err)
		}

		if expireAt.Valid {
			notification.ExpireAt = &expireAt.Time
		}

		if metadataJSON.Valid {
			metadata, err := deserializeMetadata(metadataJSON.String)
			if err != nil {
				return nil, fmt.Errorf("failed to deserialize metadata: %w", err)
			}
			notification.Metadata = metadata
		}

		notifications = append(notifications, notification)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating notifications: %w", err)
	}

	// For each notification, get the targets and acknowledgments
	for i := range notifications {
		targets, err := s.GetTargetsForNotification(ctx, notifications[i].ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get targets: %w", err)
		}
		notifications[i].Targets = targets

		acks, err := s.GetAcknowledgmentsForNotification(ctx, notifications[i].ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get acknowledgments: %w", err)
		}
		notifications[i].Acknowledgments = acks
	}

	return notifications, nil
}

// UpdateNotificationStatus updates the status of a notification.
func (s *Store) UpdateNotificationStatus(ctx context.Context, tx *sql.Tx, id int64, status NotificationStatus) error {
	var err error
	if tx != nil {
		_, err = tx.ExecContext(ctx, "UPDATE notifications SET status = ? WHERE id = ?", status, id)
	} else {
		_, err = s.db.ExecContext(ctx, "UPDATE notifications SET status = ? WHERE id = ?", status, id)
	}
	if err != nil {
		return fmt.Errorf("failed to update notification status: %w", err)
	}
	return nil
}

// DeleteNotification deletes a notification by ID.
func (s *Store) DeleteNotification(ctx context.Context, id int64) error {
	_, err := s.db.ExecContext(ctx, "DELETE FROM notifications WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("failed to delete notification: %w", err)
	}
	return nil
}

// AddNotificationTarget adds a target to a notification.
func (s *Store) AddNotificationTarget(ctx context.Context, tx *sql.Tx, target *NotificationTarget) (int64, error) {
	query := `
		INSERT INTO notification_targets 
		(notification_id, target_type, target_id, status, sent_at, external_id, response_data)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	var result sql.Result
	var err error
	if tx != nil {
		result, err = tx.ExecContext(ctx, query,
			target.NotificationID,
			target.TargetType,
			target.TargetID,
			target.Status,
			target.SentAt,
			target.ExternalID,
			target.ResponseData,
		)
	} else {
		result, err = s.db.ExecContext(ctx, query,
			target.NotificationID,
			target.TargetType,
			target.TargetID,
			target.Status,
			target.SentAt,
			target.ExternalID,
			target.ResponseData,
		)
	}

	if err != nil {
		return 0, fmt.Errorf("failed to add notification target: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert ID: %w", err)
	}

	return id, nil
}

// UpdateTargetStatus updates the status of a notification target.
func (s *Store) UpdateTargetStatus(ctx context.Context, tx *sql.Tx, id int64, status TargetStatus, externalID, responseData string) error {
	query := `
		UPDATE notification_targets
		SET status = ?, external_id = ?, response_data = ?, sent_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	var err error
	if tx != nil {
		_, err = tx.ExecContext(ctx, query, status, externalID, responseData, id)
	} else {
		_, err = s.db.ExecContext(ctx, query, status, externalID, responseData, id)
	}

	if err != nil {
		return fmt.Errorf("failed to update target status: %w", err)
	}

	return nil
}

// GetTargetsForNotification retrieves all targets for a notification.
func (s *Store) GetTargetsForNotification(ctx context.Context, notificationID int64) ([]NotificationTarget, error) {
	query := `
		SELECT id, notification_id, target_type, target_id, status, sent_at, external_id, response_data
		FROM notification_targets
		WHERE notification_id = ?
	`

	rows, err := s.db.QueryContext(ctx, query, notificationID)
	if err != nil {
		return nil, fmt.Errorf("failed to query notification targets: %w", err)
	}
	defer rows.Close()

	var targets []NotificationTarget
	for rows.Next() {
		var target NotificationTarget
		var sentAt sql.NullTime

		err := rows.Scan(
			&target.ID,
			&target.NotificationID,
			&target.TargetType,
			&target.TargetID,
			&target.Status,
			&sentAt,
			&target.ExternalID,
			&target.ResponseData,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan notification target: %w", err)
		}

		if sentAt.Valid {
			target.SentAt = &sentAt.Time
		}

		targets = append(targets, target)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating notification targets: %w", err)
	}

	return targets, nil
}
