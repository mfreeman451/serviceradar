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

package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

// Migration_MigrateAlertsToNotifications migrates data from the old alerts system to the new notification system
func Migration_MigrateAlertsToNotifications(db *sql.DB) error {
	log.Println("Running migration: Migrating alerts to notifications system")

	// Check if alerts table exists
	var tableExists bool
	err := db.QueryRow(`
		SELECT COUNT(*) > 0 FROM sqlite_master 
		WHERE type='table' AND name='alerts'
	`).Scan(&tableExists)

	if err != nil {
		return fmt.Errorf("failed to check if alerts table exists: %w", err)
	}

	if !tableExists {
		log.Println("Alerts table does not exist, skipping migration")
		return nil
	}

	// Start a transaction
	tx, err := db.BeginTx(context.Background(), nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Rollback in case of error
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Printf("Error rolling back transaction: %v", rbErr)
			}
		}
	}()

	// Query all alerts from old system
	rows, err := tx.Query(`
		SELECT 
			id, 
			node_id, 
			service_name, 
			level, 
			title, 
			message, 
			acknowledged, 
			created_at,
			acknowledged_at,
			acknowledged_by,
			details
		FROM alerts
		ORDER BY created_at DESC
	`)

	if err != nil {
		return fmt.Errorf("failed to query alerts: %w", err)
	}
	defer rows.Close()

	// Migrate each alert to a notification
	var migratedCount int
	for rows.Next() {
		var (
			id, nodeID, serviceName, level, title, message string
			acknowledged                                   bool
			createdAt                                      time.Time
			acknowledgedAt                                 sql.NullTime
			acknowledgedBy                                 sql.NullString
			detailsJSON                                    sql.NullString
		)

		if err := rows.Scan(
			&id,
			&nodeID,
			&serviceName,
			&level,
			&title,
			&message,
			&acknowledged,
			&createdAt,
			&acknowledgedAt,
			&acknowledgedBy,
			&detailsJSON,
		); err != nil {
			return fmt.Errorf("failed to scan alert row: %w", err)
		}

		// Map alert level to notification level
		var notificationLevel string
		switch level {
		case "info":
			notificationLevel = "info"
		case "warning":
			notificationLevel = "warning"
		case "error":
			notificationLevel = "error"
		default:
			notificationLevel = "info"
		}

		// Determine notification status based on acknowledgment
		var notificationStatus string
		if acknowledged {
			notificationStatus = "acknowledged"
		} else {
			notificationStatus = "sent" // Assume it was sent successfully
		}

		// Parse details JSON
		var details map[string]interface{}
		if detailsJSON.Valid && detailsJSON.String != "" {
			if err := json.Unmarshal([]byte(detailsJSON.String), &details); err != nil {
				log.Printf("Warning: Failed to parse details JSON for alert %s: %v", id, err)
				details = map[string]interface{}{
					"raw_details": detailsJSON.String,
				}
			}
		}

		// Serialize details to JSON string
		var metadataStr string
		if len(details) > 0 {
			metadata, err := json.Marshal(details)
			if err != nil {
				log.Printf("Warning: Failed to marshal details for alert %s: %v", id, err)
			} else {
				metadataStr = string(metadata)
			}
		}

		// Create notification
		notificationResult, err := tx.Exec(`
			INSERT INTO notifications (
				alert_id, 
				node_id, 
				service_name, 
				level, 
				title, 
				message, 
				status, 
				created_at, 
				updated_at,
				metadata
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`,
			id, // Use the old alert ID as the alert_id
			nodeID,
			serviceName,
			notificationLevel,
			title,
			message,
			notificationStatus,
			createdAt,
			time.Now(), // updated_at
			metadataStr,
		)

		if err != nil {
			return fmt.Errorf("failed to insert notification for alert %s: %w", id, err)
		}

		notificationID, err := notificationResult.LastInsertId()
		if err != nil {
			return fmt.Errorf("failed to get notification ID for alert %s: %w", id, err)
		}

		// If acknowledged, create an acknowledgment record
		if acknowledged && acknowledgedAt.Valid {
			ackBy := "unknown"
			if acknowledgedBy.Valid && acknowledgedBy.String != "" {
				ackBy = acknowledgedBy.String
			}

			_, err := tx.Exec(`
				INSERT INTO acknowledgments (
					notification_id,
					acknowledged_by,
					acknowledged_at,
					method,
					comment
				) VALUES (?, ?, ?, ?, ?)
			`,
				notificationID,
				ackBy,
				acknowledgedAt.Time,
				"migrated", // Use a special method to indicate this was migrated
				"Migrated from legacy alert system",
			)

			if err != nil {
				log.Printf("Warning: Failed to create acknowledgment for alert %s: %v", id, err)
			}
		}

		migratedCount++
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("error iterating alerts: %w", err)
	}

	// Rename the old alerts table to keep the data but prevent further use
	_, err = tx.Exec("ALTER TABLE alerts RENAME TO alerts_legacy")
	if err != nil {
		return fmt.Errorf("failed to rename alerts table: %w", err)
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	log.Printf("Successfully migrated %d alerts to notifications", migratedCount)
	return nil
}

// Add this migration to the migrations slice in the Initialize function
func init() {
	// Add this migration after adding the notification tables
	Migrations = append(Migrations, Migration_MigrateAlertsToNotifications)
}
