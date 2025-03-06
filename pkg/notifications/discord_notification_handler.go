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
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

const (
	// Discord embed colors
	DiscordColorRed    = 15158332 // Error
	DiscordColorYellow = 16776960 // Warning
	DiscordColorBlue   = 3447003  // Info
)

// Discord webhook URL pattern
var discordWebhookPattern = regexp.MustCompile(`^https://discord.com/api/webhooks/\d+/[\w-]+$`)

// DiscordHandlerConfig defines the configuration for a Discord webhook.
type DiscordHandlerConfig struct {
	WebhookURL       string        `json:"webhook_url"`
	Username         string        `json:"username,omitempty"`
	AvatarURL        string        `json:"avatar_url,omitempty"`
	Timeout          time.Duration `json:"timeout,omitempty"`
	MaxRetries       int           `json:"max_retries,omitempty"`
	RetryDelay       time.Duration `json:"retry_delay,omitempty"`
	MentionRoles     []string      `json:"mention_roles,omitempty"`
	MentionUsers     []string      `json:"mention_users,omitempty"`
	MentionEveryone  bool          `json:"mention_everyone,omitempty"`
	IncludeTimestamp bool          `json:"include_timestamp,omitempty"`
}

// DiscordHandler implements the NotificationHandler interface for Discord webhooks.
type DiscordHandler struct {
	configs map[string]DiscordHandlerConfig // Map of configs by targetID
	client  *http.Client
}

// DiscordEmbed represents a Discord message embed.
type DiscordEmbed struct {
	Title       string              `json:"title"`
	Description string              `json:"description"`
	Color       int                 `json:"color"`
	Fields      []DiscordEmbedField `json:"fields,omitempty"`
	Timestamp   string              `json:"timestamp,omitempty"`
	Footer      *DiscordEmbedFooter `json:"footer,omitempty"`
}

// DiscordEmbedField represents a field in a Discord embed.
type DiscordEmbedField struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline,omitempty"`
}

// DiscordEmbedFooter represents the footer in a Discord embed.
type DiscordEmbedFooter struct {
	Text    string `json:"text"`
	IconURL string `json:"icon_url,omitempty"`
}

// DiscordWebhookPayload represents the payload for a Discord webhook.
type DiscordWebhookPayload struct {
	Content    string         `json:"content,omitempty"`
	Username   string         `json:"username,omitempty"`
	AvatarURL  string         `json:"avatar_url,omitempty"`
	Embeds     []DiscordEmbed `json:"embeds,omitempty"`
	Components []interface{}  `json:"components,omitempty"`
	TTS        bool           `json:"tts,omitempty"`
}

// DiscordButton represents a Discord button component.
type DiscordButton struct {
	Type     int    `json:"type"`
	Style    int    `json:"style"`
	Label    string `json:"label"`
	CustomID string `json:"custom_id"`
	URL      string `json:"url,omitempty"`
	Emoji    *struct {
		Name string `json:"name"`
		ID   string `json:"id,omitempty"`
	} `json:"emoji,omitempty"`
	Disabled bool `json:"disabled,omitempty"`
}

// DiscordActionRow represents a row of components in a Discord message.
type DiscordActionRow struct {
	Type       int             `json:"type"`
	Components []DiscordButton `json:"components"`
}

// NewDiscordHandler creates a new Discord webhook handler.
func NewDiscordHandler(configs map[string]DiscordHandlerConfig) *DiscordHandler {
	return &DiscordHandler{
		configs: configs,
		client: &http.Client{
			Timeout: 10 * time.Second, // Default timeout
		},
	}
}

// AddConfig adds or updates a Discord webhook configuration.
func (h *DiscordHandler) AddConfig(targetID string, config DiscordHandlerConfig) error {
	// Validate the webhook URL
	if !discordWebhookPattern.MatchString(config.WebhookURL) {
		return fmt.Errorf("invalid Discord webhook URL format")
	}

	// Set defaults
	if config.Timeout == 0 {
		config.Timeout = 10 * time.Second
	}
	if config.MaxRetries == 0 {
		config.MaxRetries = 3
	}
	if config.RetryDelay == 0 {
		config.RetryDelay = 5 * time.Second
	}

	h.configs[targetID] = config
	return nil
}

// GetConfig retrieves a Discord webhook configuration by target ID.
func (h *DiscordHandler) GetConfig(targetID string) (DiscordHandlerConfig, error) {
	config, ok := h.configs[targetID]
	if !ok {
		return DiscordHandlerConfig{}, fmt.Errorf("no Discord configuration found for target ID: %s", targetID)
	}
	return config, nil
}

// SendNotification sends a notification to a Discord webhook.
func (h *DiscordHandler) SendNotification(ctx context.Context, notification *Notification, target *NotificationTarget) error {
	// Get config for this target
	config, err := h.GetConfig(target.TargetID)
	if err != nil {
		return err
	}

	// Prepare the payload
	payload, err := h.preparePayload(notification, config)
	if err != nil {
		return fmt.Errorf("failed to prepare payload: %w", err)
	}

	// Create request
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, config.WebhookURL, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "ServiceRadar-Notifications/1.0")

	// Set client timeout
	client := h.client
	if config.Timeout > 0 {
		client = &http.Client{
			Timeout: config.Timeout,
		}
	}

	// Send the request with retries
	var resp *http.Response
	var lastErr error

	for i := 0; i <= config.MaxRetries; i++ {
		resp, lastErr = client.Do(req)
		if lastErr == nil && resp.StatusCode >= 200 && resp.StatusCode < 300 {
			// Success
			responseBody, _ := io.ReadAll(resp.Body)
			resp.Body.Close()

			// Store response data
			target.ResponseData = string(responseBody)
			target.Status = TargetStatusSent

			// Parse message ID from response for acknowledgment
			var respMap map[string]interface{}
			if json.Unmarshal(responseBody, &respMap) == nil {
				if id, ok := respMap["id"]; ok {
					target.ExternalID = fmt.Sprintf("%v", id)
				}
			}

			return nil
		}

		if resp != nil {
			resp.Body.Close()
		}

		// If this was the last retry, we'll exit the loop
		if i == config.MaxRetries {
			break
		}

		// Wait before retrying
		select {
		case <-time.After(config.RetryDelay):
			// Continue with retry
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	// If we got here, all retries failed
	responseBody := "No response"
	if resp != nil {
		body, _ := io.ReadAll(resp.Body)
		responseBody = string(body)
	}

	target.Status = TargetStatusFailed
	target.ResponseData = fmt.Sprintf("Failed after %d retries. Last error: %v. Last response: %s",
		config.MaxRetries, lastErr, responseBody)

	return fmt.Errorf("Discord webhook request failed after %d retries: %w", config.MaxRetries, lastErr)
}

// preparePayload prepares the payload for a Discord webhook.
func (h *DiscordHandler) preparePayload(notification *Notification, config DiscordHandlerConfig) ([]byte, error) {
	// Determine embed color based on notification level
	var color int
	switch notification.Level {
	case LevelError:
		color = DiscordColorRed
	case LevelWarning:
		color = DiscordColorYellow
	default:
		color = DiscordColorBlue
	}

	// Create basic embed
	embed := DiscordEmbed{
		Title:       notification.Title,
		Description: notification.Message,
		Color:       color,
		Fields:      []DiscordEmbedField{},
	}

	// Add timestamp if configured
	if config.IncludeTimestamp {
		embed.Timestamp = notification.CreatedAt.Format(time.RFC3339)
	}

	// Add footer with ID for reference
	embed.Footer = &DiscordEmbedFooter{
		Text: fmt.Sprintf("Notification ID: %d", notification.ID),
	}

	// Add fields for additional context
	if notification.NodeID != "" {
		embed.Fields = append(embed.Fields, DiscordEmbedField{
			Name:   "Node",
			Value:  notification.NodeID,
			Inline: true,
		})
	}

	if notification.ServiceName != "" {
		embed.Fields = append(embed.Fields, DiscordEmbedField{
			Name:   "Service",
			Value:  notification.ServiceName,
			Inline: true,
		})
	}

	// Add metadata fields if present
	if notification.Metadata != nil {
		for key, value := range notification.Metadata {
			// Skip large values
			valueStr := fmt.Sprintf("%v", value)
			if len(valueStr) > 100 {
				valueStr = valueStr[:97] + "..."
			}

			embed.Fields = append(embed.Fields, DiscordEmbedField{
				Name:   key,
				Value:  valueStr,
				Inline: true,
			})
		}
	}

	// Build message content with mentions
	var content string
	if len(config.MentionRoles) > 0 {
		for _, role := range config.MentionRoles {
			content += fmt.Sprintf("<@&%s> ", role)
		}
	}

	if len(config.MentionUsers) > 0 {
		for _, user := range config.MentionUsers {
			content += fmt.Sprintf("<@%s> ", user)
		}
	}

	if config.MentionEveryone {
		content += "@everyone "
	}

	// Create payload with buttons for acknowledgment
	payload := DiscordWebhookPayload{
		Content:   content,
		Username:  config.Username,
		AvatarURL: config.AvatarURL,
		Embeds:    []DiscordEmbed{embed},
	}

	// Add acknowledge button
	ackButton := DiscordButton{
		Type:     2, // Button
		Style:    3, // Success (green)
		Label:    "Acknowledge",
		CustomID: fmt.Sprintf("ack_%d", notification.ID),
	}

	viewButton := DiscordButton{
		Type:  2, // Button
		Style: 5, // Link
		Label: "View Details",
		URL:   fmt.Sprintf("https://serviceradar.carverauto.com/notifications/%d", notification.ID),
	}

	// Create action row with buttons
	actionRow := DiscordActionRow{
		Type:       1, // Action Row
		Components: []DiscordButton{ackButton, viewButton},
	}

	payload.Components = []interface{}{actionRow}

	return json.Marshal(payload)
}

// ParseAcknowledgment parses an acknowledgment from a Discord interaction.
func (h *DiscordHandler) ParseAcknowledgment(ctx context.Context, data []byte) (*AcknowledgeRequest, error) {
	var interaction struct {
		Type int `json:"type"`
		Data struct {
			CustomID      string `json:"custom_id"`
			ComponentType int    `json:"component_type"`
		} `json:"data"`
		Member struct {
			User struct {
				ID       string `json:"id"`
				Username string `json:"username"`
			} `json:"user"`
		} `json:"member"`
	}

	if err := json.Unmarshal(data, &interaction); err != nil {
		return nil, fmt.Errorf("failed to parse Discord interaction: %w", err)
	}

	// Verify it's a button interaction for acknowledgment
	if interaction.Type != 3 || interaction.Data.ComponentType != 2 {
		return nil, fmt.Errorf("not a button interaction")
	}

	if interaction.Data.CustomID == "" || !regexp.MustCompile(`^ack_\d+$`).MatchString(interaction.Data.CustomID) {
		return nil, fmt.Errorf("not an acknowledgment interaction")
	}

	// Extract notification ID from custom_id
	idStr := interaction.Data.CustomID[4:] // Skip "ack_" prefix
	notificationID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid notification ID in interaction: %w", err)
	}

	// Create acknowledgment request
	ackBy := interaction.Member.User.Username
	if ackBy == "" {
		ackBy = fmt.Sprintf("discord-user-%s", interaction.Member.User.ID)
	}

	return &AcknowledgeRequest{
		NotificationID: notificationID,
		AcknowledgedBy: ackBy,
		Method:         AckMethodDiscord,
		Comment:        fmt.Sprintf("Acknowledged via Discord by user %s", ackBy),
	}, nil
}

// RespondToInteraction sends a response to a Discord interaction.
func (h *DiscordHandler) RespondToInteraction(interactionID, interactionToken string, message string) error {
	url := fmt.Sprintf("https://discord.com/api/v10/interactions/%s/%s/callback", interactionID, interactionToken)

	payload := map[string]interface{}{
		"type": 4, // Channel message with source
		"data": map[string]interface{}{
			"content": message,
			"flags":   64, // Ephemeral (only visible to the user who clicked)
		},
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(jsonPayload))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := h.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("Discord API error: %d %s", resp.StatusCode, string(body))
	}

	return nil
}
