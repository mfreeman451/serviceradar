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

// Package notifications provides types and interfaces for the notification system.
package notifications

import (
	"context"
	"database/sql"
	"fmt"
)

// GetAcknowledgmentsForNotification retrieves all acknowledgments for a notification.
func (s *Store) GetAcknowledgmentsForNotification(ctx context.Context, notificationID int64) ([]Acknowledgment, error) {
	query := `
		SELECT id, notification_id, target_id, acknowledged_by, acknowledged_at, method, comment
		FROM acknowledgments
		WHERE notification_id = ?
		ORDER BY acknowledged_at DESC
	`

	rows, err := s.db.QueryContext(ctx, query, notificationID)
	if err != nil {
		return nil, fmt.Errorf("failed to query acknowledgments: %w", err)
	}
	defer rows.Close()

	var acknowledgments []Acknowledgment
	for rows.Next() {
		var ack Acknowledgment
		var targetID sql.NullInt64

		err := rows.Scan(
			&ack.ID,
			&ack.NotificationID,
			&targetID,
			&ack.AcknowledgedBy,
			&ack.AcknowledgedAt,
			&ack.Method,
			&ack.Comment,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan acknowledgment: %w", err)
		}

		if targetID.Valid {
			ack.TargetID = &targetID.Int64
		}

		acknowledgments = append(acknowledgments, ack)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating acknowledgments: %w", err)
	}

	return acknowledgments, nil
}

// Additional helper methods for acknowledgments that might be needed

// HasAcknowledgment checks if a notification has been acknowledged.
func (s *Store) HasAcknowledgment(ctx context.Context, notificationID int64) (bool, error) {
	var count int
	err := s.db.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM acknowledgments WHERE notification_id = ?",
		notificationID).Scan(&count)

	if err != nil {
		return false, fmt.Errorf("failed to check for acknowledgments: %w", err)
	}

	return count > 0, nil
}

// GetAcknowledgmentByTarget retrieves an acknowledgment for a specific target.
func (s *Store) GetAcknowledgmentByTarget(ctx context.Context, notificationID int64, targetID int64) (*Acknowledgment, error) {
	var ack Acknowledgment

	err := s.db.QueryRowContext(ctx, `
		SELECT id, notification_id, target_id, acknowledged_by, acknowledged_at, method, comment
		FROM acknowledgments
		WHERE notification_id = ? AND target_id = ?
		LIMIT 1`,
		notificationID, targetID).Scan(
		&ack.ID,
		&ack.NotificationID,
		&ack.TargetID,
		&ack.AcknowledgedBy,
		&ack.AcknowledgedAt,
		&ack.Method,
		&ack.Comment,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No acknowledgment found
		}
		return nil, fmt.Errorf("failed to get acknowledgment: %w", err)
	}

	return &ack, nil
}

// DeleteAcknowledgments deletes all acknowledgments for a notification.
func (s *Store) DeleteAcknowledgments(ctx context.Context, tx *sql.Tx, notificationID int64) error {
	var err error
	if tx != nil {
		_, err = tx.ExecContext(ctx,
			"DELETE FROM acknowledgments WHERE notification_id = ?",
			notificationID)
	} else {
		_, err = s.db.ExecContext(ctx,
			"DELETE FROM acknowledgments WHERE notification_id = ?",
			notificationID)
	}

	if err != nil {
		return fmt.Errorf("failed to delete acknowledgments: %w", err)
	}

	return nil
}
