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

import "time"

// NotificationLevel represents the severity of a notification.
type NotificationLevel string

const (
	LevelInfo    NotificationLevel = "info"
	LevelWarning NotificationLevel = "warning"
	LevelError   NotificationLevel = "error"
)

// NotificationStatus represents the current status of a notification.
type NotificationStatus string

const (
	StatusPending      NotificationStatus = "pending"
	StatusSent         NotificationStatus = "sent"
	StatusAcknowledged NotificationStatus = "acknowledged"
	StatusResolved     NotificationStatus = "resolved"
	StatusExpired      NotificationStatus = "expired"
)

// TargetStatus represents the delivery status of a notification to a specific target.
type TargetStatus string

const (
	TargetStatusPending      TargetStatus = "pending"
	TargetStatusSent         TargetStatus = "sent"
	TargetStatusFailed       TargetStatus = "failed"
	TargetStatusAcknowledged TargetStatus = "acknowledged"
)

// TargetType represents the type of notification target.
type TargetType string

const (
	TargetTypeWebhook TargetType = "webhook"
	TargetTypeDiscord TargetType = "discord"
	TargetTypeSlack   TargetType = "slack"
	TargetTypeMsTeams TargetType = "msteams"
	TargetTypeEmail   TargetType = "email"
	TargetTypeSMS     TargetType = "sms"
)

// AcknowledgmentMethod represents how a notification was acknowledged.
type AcknowledgmentMethod string

const (
	AckMethodAPI     AcknowledgmentMethod = "api"
	AckMethodWebhook AcknowledgmentMethod = "webhook"
	AckMethodUI      AcknowledgmentMethod = "ui"
	AckMethodDiscord AcknowledgmentMethod = "discord"
	AckMethodSlack   AcknowledgmentMethod = "slack"
	AckMethodEmail   AcknowledgmentMethod = "email"
)

// Notification represents a notification sent to one or more targets.
type Notification struct {
	ID              int64                  `json:"id"`
	AlertID         string                 `json:"alert_id"`     // Unique ID for deduplication
	NodeID          string                 `json:"node_id"`      // Node that generated the notification
	ServiceName     string                 `json:"service_name"` // Service that generated the notification
	Level           NotificationLevel      `json:"level"`
	Title           string                 `json:"title"`
	Message         string                 `json:"message"`
	Status          NotificationStatus     `json:"status"`
	CreatedAt       time.Time              `json:"created_at"`
	UpdatedAt       time.Time              `json:"updated_at"`
	ExpireAt        *time.Time             `json:"expire_at,omitempty"`
	Metadata        map[string]interface{} `json:"metadata,omitempty"`
	Targets         []NotificationTarget   `json:"targets,omitempty"`
	Acknowledgments []Acknowledgment       `json:"acknowledgments,omitempty"`
}

// NotificationTarget represents a specific delivery target for a notification.
type NotificationTarget struct {
	ID             int64        `json:"id"`
	NotificationID int64        `json:"notification_id"`
	TargetType     TargetType   `json:"target_type"`
	TargetID       string       `json:"target_id"`
	Status         TargetStatus `json:"status"`
	SentAt         *time.Time   `json:"sent_at,omitempty"`
	ExternalID     string       `json:"external_id,omitempty"`
	ResponseData   string       `json:"response_data,omitempty"`
}

// Acknowledgment represents an acknowledgment of a notification.
type Acknowledgment struct {
	ID             int64                `json:"id"`
	NotificationID int64                `json:"notification_id"`
	TargetID       *int64               `json:"target_id,omitempty"`
	AcknowledgedBy string               `json:"acknowledged_by"`
	AcknowledgedAt time.Time            `json:"acknowledged_at"`
	Method         AcknowledgmentMethod `json:"method"`
	Comment        string               `json:"comment,omitempty"`
}

// APIKey represents an API key for notification services.
type APIKey struct {
	ID          int64      `json:"id"`
	KeyID       string     `json:"key_id"`
	KeyHash     string     `json:"-"` // Hash is never exposed in JSON
	ServiceName string     `json:"service_name"`
	CreatedAt   time.Time  `json:"created_at"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
	Permissions []string   `json:"permissions"`
	LastUsedAt  *time.Time `json:"last_used_at,omitempty"`
}

// NotificationRequest is used to create a new notification.
type NotificationRequest struct {
	AlertID     string                 `json:"alert_id"`
	NodeID      string                 `json:"node_id,omitempty"`
	ServiceName string                 `json:"service_name,omitempty"`
	Level       NotificationLevel      `json:"level"`
	Title       string                 `json:"title"`
	Message     string                 `json:"message"`
	ExpireAfter *time.Duration         `json:"expire_after,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	Targets     []TargetRequest        `json:"targets,omitempty"`
}

// TargetRequest specifies a notification target.
type TargetRequest struct {
	TargetType TargetType `json:"target_type"`
	TargetID   string     `json:"target_id"`
}

// AcknowledgeRequest is used to acknowledge a notification.
type AcknowledgeRequest struct {
	NotificationID int64                `json:"notification_id"`
	TargetID       *int64               `json:"target_id,omitempty"`
	AcknowledgedBy string               `json:"acknowledged_by"`
	Method         AcknowledgmentMethod `json:"method"`
	Comment        string               `json:"comment,omitempty"`
}

// NotificationFilter is used to filter notifications in queries.
type NotificationFilter struct {
	AlertID      string              `json:"alert_id,omitempty"`
	NodeID       string              `json:"node_id,omitempty"`
	ServiceName  string              `json:"service_name,omitempty"`
	Level        *NotificationLevel  `json:"level,omitempty"`
	Status       *NotificationStatus `json:"status,omitempty"`
	Acknowledged *bool               `json:"acknowledged,omitempty"`
	Since        *time.Time          `json:"since,omitempty"`
	Until        *time.Time          `json:"until,omitempty"`
	Limit        int                 `json:"limit,omitempty"`
	Offset       int                 `json:"offset,omitempty"`
}

// CleanupConfig defines configuration for the notification cleanup service.
type CleanupConfig struct {
	// How often to run the cleanup process
	Interval time.Duration

	// Maximum age of resolved notifications to keep
	ResolvedRetention time.Duration

	// Maximum age of acknowledged notifications to keep
	AcknowledgedRetention time.Duration

	// Maximum age of pending notifications to keep
	PendingRetention time.Duration

	// Whether to delete or just update expired notifications
	DeleteExpired bool
}

// CleanupService manages periodic cleanup of notifications.
type CleanupService struct {
	service NotificationService
	config  CleanupConfig
	stopCh  chan struct{}
}
