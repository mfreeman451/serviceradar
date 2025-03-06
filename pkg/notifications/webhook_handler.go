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
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"text/template"
	"time"
)

// WebhookHandlerConfig defines the configuration for a webhook handler.
type WebhookHandlerConfig struct {
	URL             string            `json:"url"`                        // Webhook endpoint URL
	Method          string            `json:"method"`                     // HTTP method (default: POST)
	Headers         map[string]string `json:"headers,omitempty"`          // Custom headers
	Template        string            `json:"template,omitempty"`         // Custom template
	Secret          string            `json:"secret,omitempty"`           // Secret for HMAC signing
	Timeout         time.Duration     `json:"timeout,omitempty"`          // Request timeout
	MaxRetries      int               `json:"max_retries,omitempty"`      // Maximum number of retries
	RetryDelay      time.Duration     `json:"retry_delay,omitempty"`      // Delay between retries
	VerificationURL string            `json:"verification_url,omitempty"` // URL for acknowledgment verification
	SignatureHeader string            `json:"signature_header,omitempty"` // Header for signature
}

// WebhookHandler implements the NotificationHandler interface for webhooks.
type WebhookHandler struct {
	configs map[string]WebhookHandlerConfig // Map of configs by targetID
	client  *http.Client
}

// NewWebhookHandler creates a new webhook handler.
func NewWebhookHandler(configs map[string]WebhookHandlerConfig) *WebhookHandler {
	return &WebhookHandler{
		configs: configs,
		client: &http.Client{
			Timeout: 10 * time.Second, // Default timeout
		},
	}
}

// AddConfig adds or updates a webhook configuration.
func (h *WebhookHandler) AddConfig(targetID string, config WebhookHandlerConfig) {
	// Set defaults
	if config.Method == "" {
		config.Method = http.MethodPost
	}
	if config.Timeout == 0 {
		config.Timeout = 10 * time.Second
	}
	if config.MaxRetries == 0 {
		config.MaxRetries = 3
	}
	if config.RetryDelay == 0 {
		config.RetryDelay = 5 * time.Second
	}
	if config.SignatureHeader == "" {
		config.SignatureHeader = "X-Signature"
	}

	h.configs[targetID] = config
}

// GetConfig retrieves a webhook configuration by target ID.
func (h *WebhookHandler) GetConfig(targetID string) (WebhookHandlerConfig, error) {
	config, ok := h.configs[targetID]
	if !ok {
		return WebhookHandlerConfig{}, fmt.Errorf("no webhook configuration found for target ID: %s", targetID)
	}
	return config, nil
}

// SendNotification sends a notification to a webhook.
func (h *WebhookHandler) SendNotification(ctx context.Context, notification *Notification, target *NotificationTarget) error {
	// Get config for this target
	config, err := h.GetConfig(target.TargetID)
	if err != nil {
		return err
	}

	// Prepare the payload
	payload, err := h.preparePayload(notification, target, config)
	if err != nil {
		return fmt.Errorf("failed to prepare payload: %w", err)
	}

	// Create request
	req, err := http.NewRequestWithContext(ctx, config.Method, config.URL, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set default headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "ServiceRadar-Notifications/1.0")

	// Set custom headers
	for key, value := range config.Headers {
		req.Header.Set(key, value)
	}

	// Add signature if secret is provided
	if config.Secret != "" {
		signature := h.signPayload(payload, config.Secret)
		req.Header.Set(config.SignatureHeader, signature)
	}

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

			// Set a unique ID from response if available (for acknowledgment)
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

	return fmt.Errorf("webhook request failed after %d retries: %w", config.MaxRetries, lastErr)
}

// preparePayload prepares the payload for a webhook request.
func (h *WebhookHandler) preparePayload(notification *Notification, target *NotificationTarget, config WebhookHandlerConfig) ([]byte, error) {
	// If custom template is provided, use it
	if config.Template != "" {
		tmpl, err := template.New("webhook").Parse(config.Template)
		if err != nil {
			return nil, fmt.Errorf("failed to parse template: %w", err)
		}

		var buf bytes.Buffer
		if err := tmpl.Execute(&buf, map[string]interface{}{
			"notification": notification,
			"target":       target,
		}); err != nil {
			return nil, fmt.Errorf("failed to execute template: %w", err)
		}

		return buf.Bytes(), nil
	}

	// Default payload structure
	payload := map[string]interface{}{
		"id":           notification.ID,
		"alert_id":     notification.AlertID,
		"level":        notification.Level,
		"title":        notification.Title,
		"message":      notification.Message,
		"created_at":   notification.CreatedAt.Format(time.RFC3339),
		"node_id":      notification.NodeID,
		"service_name": notification.ServiceName,
	}

	// Include metadata if present
	if notification.Metadata != nil {
		payload["metadata"] = notification.Metadata
	}

	// Include verification URL if configured
	if config.VerificationURL != "" {
		payload["verification_url"] = fmt.Sprintf("%s/%d", config.VerificationURL, notification.ID)
	}

	return json.Marshal(payload)
}

// signPayload creates an HMAC signature for the payload.
func (h *WebhookHandler) signPayload(payload []byte, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	return hex.EncodeToString(mac.Sum(nil))
}

// ParseAcknowledgment parses an acknowledgment from a webhook callback.
func (h *WebhookHandler) ParseAcknowledgment(ctx context.Context, data []byte) (*AcknowledgeRequest, error) {
	var ackData struct {
		NotificationID int64  `json:"notification_id"`
		TargetID       *int64 `json:"target_id,omitempty"`
		AcknowledgedBy string `json:"acknowledged_by"`
		Comment        string `json:"comment,omitempty"`
	}

	if err := json.Unmarshal(data, &ackData); err != nil {
		return nil, fmt.Errorf("failed to parse acknowledgment data: %w", err)
	}

	if ackData.NotificationID == 0 || ackData.AcknowledgedBy == "" {
		return nil, errors.New("invalid acknowledgment data: missing required fields")
	}

	return &AcknowledgeRequest{
		NotificationID: ackData.NotificationID,
		TargetID:       ackData.TargetID,
		AcknowledgedBy: ackData.AcknowledgedBy,
		Method:         AckMethodWebhook,
		Comment:        ackData.Comment,
	}, nil
}

// VerifySignature verifies the signature of an incoming webhook callback.
func (h *WebhookHandler) VerifySignature(data []byte, signature, targetID string) bool {
	config, err := h.GetConfig(targetID)
	if err != nil {
		log.Printf("Error getting config for signature verification: %v", err)
		return false
	}

	if config.Secret == "" {
		// No signature verification configured
		return true
	}

	expectedSignature := h.signPayload(data, config.Secret)
	return hmac.Equal([]byte(signature), []byte(expectedSignature))
}
