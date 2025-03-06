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
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

// APIHandler handles HTTP requests for the notification API.
type APIHandler struct {
	service  NotificationService
	handlers map[TargetType]NotificationHandler
}

// NewAPIHandler creates a new notification API handler.
func NewAPIHandler(service NotificationService, handlers map[TargetType]NotificationHandler) *APIHandler {
	return &APIHandler{
		service:  service,
		handlers: handlers,
	}
}

// RegisterRoutes registers the notification API routes with the router.
func (h *APIHandler) RegisterRoutes(router *mux.Router) {
	// Notifications routes
	router.HandleFunc("/api/notifications", h.CreateNotification).Methods(http.MethodPost)
	router.HandleFunc("/api/notifications", h.ListNotifications).Methods(http.MethodGet)
	router.HandleFunc("/api/notifications/{id:[0-9]+}", h.GetNotification).Methods(http.MethodGet)
	router.HandleFunc("/api/notifications/{id:[0-9]+}/acknowledge", h.AcknowledgeNotification).Methods(http.MethodPost)
	router.HandleFunc("/api/notifications/{id:[0-9]+}/resolve", h.ResolveNotification).Methods(http.MethodPost)
	router.HandleFunc("/api/notifications/{id:[0-9]+}", h.DeleteNotification).Methods(http.MethodDelete)

	// Webhook callback routes
	router.HandleFunc("/api/notifications/callbacks/{target_type}/{target_id}", h.HandleCallback).Methods(http.MethodPost)

	// API key management routes
	router.HandleFunc("/api/notifications/api-keys", h.CreateAPIKey).Methods(http.MethodPost)
	router.HandleFunc("/api/notifications/api-keys/{key_id}", h.RevokeAPIKey).Methods(http.MethodDelete)
}

// respondJSON sends a JSON response.
func respondJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshaling JSON response: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// respondError sends an error response.
func respondError(w http.ResponseWriter, code int, message string) {
	respondJSON(w, code, map[string]string{"error": message})
}

// validateAPIKey validates the API key from the request.
func (h *APIHandler) validateAPIKey(r *http.Request, requiredPermissions ...string) (*APIKey, error) {
	// Get API key header
	keyID := r.Header.Get("X-API-Key-ID")
	keyValue := r.Header.Get("X-API-Key")

	if keyID == "" || keyValue == "" {
		return nil, ErrInvalidAPIKey
	}

	// Validate the key
	key, err := h.service.ValidateAPIKey(r.Context(), keyID, keyValue)
	if err != nil {
		return nil, err
	}

	// Check permissions if required
	if len(requiredPermissions) > 0 {
		hasPermission := false
		for _, required := range requiredPermissions {
			for _, permission := range key.Permissions {
				if permission == required || permission == "*" {
					hasPermission = true
					break
				}
			}
			if hasPermission {
				break
			}
		}

		if !hasPermission {
			return nil, ErrPermissionDenied
		}
	}

	return key, nil
}

// parseNotificationIDFromURL parses the notification ID from the URL.
func parseNotificationIDFromURL(r *http.Request) (int64, error) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	if idStr == "" {
		return 0, errors.New("notification ID is required")
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid notification ID: %w", err)
	}

	return id, nil
}

// CreateNotification handles the creation of a new notification.
func (h *APIHandler) CreateNotification(w http.ResponseWriter, r *http.Request) {
	// Validate API key
	key, err := h.validateAPIKey(r, "notifications:write")
	if err != nil {
		if errors.Is(err, ErrInvalidAPIKey) || errors.Is(err, ErrAPIKeyExpired) {
			respondError(w, http.StatusUnauthorized, err.Error())
			return
		}
		if errors.Is(err, ErrPermissionDenied) {
			respondError(w, http.StatusForbidden, err.Error())
			return
		}
		respondError(w, http.StatusInternalServerError, "Failed to validate API key")
		return
	}

	// Parse request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Failed to read request body")
		return
	}

	var req NotificationRequest
	if err := json.Unmarshal(body, &req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid JSON body")
		return
	}

	// Create the notification
	notification, err := h.service.CreateNotification(r.Context(), req)
	if err != nil {
		log.Printf("Error creating notification: %v", err)

		if errors.Is(err, ErrInvalidRequest) {
			respondError(w, http.StatusBadRequest, err.Error())
			return
		}

		respondError(w, http.StatusInternalServerError, "Failed to create notification")
		return
	}

	log.Printf("Created notification %d from service %s", notification.ID, key.ServiceName)

	// Return the created notification
	respondJSON(w, http.StatusCreated, notification)
}

// GetNotification handles retrieving a notification by ID.
func (h *APIHandler) GetNotification(w http.ResponseWriter, r *http.Request) {
	// Validate API key
	_, err := h.validateAPIKey(r, "notifications:read")
	if err != nil {
		if errors.Is(err, ErrInvalidAPIKey) || errors.Is(err, ErrAPIKeyExpired) {
			respondError(w, http.StatusUnauthorized, err.Error())
			return
		}
		if errors.Is(err, ErrPermissionDenied) {
			respondError(w, http.StatusForbidden, err.Error())
			return
		}
		respondError(w, http.StatusInternalServerError, "Failed to validate API key")
		return
	}

	// Parse notification ID
	id, err := parseNotificationIDFromURL(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Get the notification
	notification, err := h.service.GetNotification(r.Context(), id)
	if err != nil {
		log.Printf("Error getting notification %d: %v", id, err)

		if errors.Is(err, ErrNotificationNotFound) {
			respondError(w, http.StatusNotFound, "Notification not found")
			return
		}

		respondError(w, http.StatusInternalServerError, "Failed to get notification")
		return
	}

	// Return the notification
	respondJSON(w, http.StatusOK, notification)
}

// ListNotifications handles retrieving a list of notifications based on filters.
func (h *APIHandler) ListNotifications(w http.ResponseWriter, r *http.Request) {
	// Validate API key
	_, err := h.validateAPIKey(r, "notifications:read")
	if err != nil {
		if errors.Is(err, ErrInvalidAPIKey) || errors.Is(err, ErrAPIKeyExpired) {
			respondError(w, http.StatusUnauthorized, err.Error())
			return
		}
		if errors.Is(err, ErrPermissionDenied) {
			respondError(w, http.StatusForbidden, err.Error())
			return
		}
		respondError(w, http.StatusInternalServerError, "Failed to validate API key")
		return
	}

	// Parse query parameters for filtering
	filter := NotificationFilter{}

	// Basic string filters
	filter.AlertID = r.URL.Query().Get("alert_id")
	filter.NodeID = r.URL.Query().Get("node_id")
	filter.ServiceName = r.URL.Query().Get("service_name")

	// Parse level if provided
	if levelStr := r.URL.Query().Get("level"); levelStr != "" {
		level := NotificationLevel(levelStr)
		filter.Level = &level
	}

	// Parse status if provided
	if statusStr := r.URL.Query().Get("status"); statusStr != "" {
		status := NotificationStatus(statusStr)
		filter.Status = &status
	}

	// Parse acknowledged flag if provided
	if ackStr := r.URL.Query().Get("acknowledged"); ackStr != "" {
		acknowledged := strings.ToLower(ackStr) == "true"
		filter.Acknowledged = &acknowledged
	}

	// Parse since timestamp if provided
	if sinceStr := r.URL.Query().Get("since"); sinceStr != "" {
		since, err := time.Parse(time.RFC3339, sinceStr)
		if err == nil {
			filter.Since = &since
		}
	}

	// Parse until timestamp if provided
	if untilStr := r.URL.Query().Get("until"); untilStr != "" {
		until, err := time.Parse(time.RFC3339, untilStr)
		if err == nil {
			filter.Until = &until
		}
	}

	// Parse pagination parameters
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 {
			filter.Limit = limit
		}
	}

	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil && offset >= 0 {
			filter.Offset = offset
		}
	}

	// Get the notifications
	notifications, err := h.service.ListNotifications(r.Context(), filter)
	if err != nil {
		log.Printf("Error listing notifications: %v", err)
		respondError(w, http.StatusInternalServerError, "Failed to list notifications")
		return
	}

	// Return the notifications
	respondJSON(w, http.StatusOK, notifications)
}

// AcknowledgeNotification handles acknowledging a notification.
func (h *APIHandler) AcknowledgeNotification(w http.ResponseWriter, r *http.Request) {
	// Validate API key
	key, err := h.validateAPIKey(r, "notifications:acknowledge")
	if err != nil {
		if errors.Is(err, ErrInvalidAPIKey) || errors.Is(err, ErrAPIKeyExpired) {
			respondError(w, http.StatusUnauthorized, err.Error())
			return
		}
		if errors.Is(err, ErrPermissionDenied) {
			respondError(w, http.StatusForbidden, err.Error())
			return
		}
		respondError(w, http.StatusInternalServerError, "Failed to validate API key")
		return
	}

	// Parse notification ID
	id, err := parseNotificationIDFromURL(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Parse request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Failed to read request body")
		return
	}

	var ackRequest struct {
		AcknowledgedBy string `json:"acknowledged_by"`
		Comment        string `json:"comment,omitempty"`
		TargetID       *int64 `json:"target_id,omitempty"`
	}

	if err := json.Unmarshal(body, &ackRequest); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid JSON body")
		return
	}

	// Use provided acknowledger or default to API key service name
	acknowledgedBy := ackRequest.AcknowledgedBy
	if acknowledgedBy == "" {
		acknowledgedBy = key.ServiceName
	}

	// Create acknowledgment request
	req := AcknowledgeRequest{
		NotificationID: id,
		TargetID:       ackRequest.TargetID,
		AcknowledgedBy: acknowledgedBy,
		Method:         AckMethodAPI,
		Comment:        ackRequest.Comment,
	}

	// Acknowledge the notification
	err = h.service.AcknowledgeNotification(r.Context(), req)
	if err != nil {
		log.Printf("Error acknowledging notification %d: %v", id, err)

		if errors.Is(err, ErrNotificationNotFound) {
			respondError(w, http.StatusNotFound, "Notification not found")
			return
		}

		if errors.Is(err, ErrAlreadyAcknowledged) {
			respondError(w, http.StatusConflict, "Notification already acknowledged")
			return
		}

		respondError(w, http.StatusInternalServerError, "Failed to acknowledge notification")
		return
	}

	// Return success
	respondJSON(w, http.StatusOK, map[string]string{
		"status":  "acknowledged",
		"message": "Notification acknowledged successfully",
	})
}

// ResolveNotification handles resolving a notification.
func (h *APIHandler) ResolveNotification(w http.ResponseWriter, r *http.Request) {
	// Validate API key
	_, err := h.validateAPIKey(r, "notifications:write")
	if err != nil {
		if errors.Is(err, ErrInvalidAPIKey) || errors.Is(err, ErrAPIKeyExpired) {
			respondError(w, http.StatusUnauthorized, err.Error())
			return
		}
		if errors.Is(err, ErrPermissionDenied) {
			respondError(w, http.StatusForbidden, err.Error())
			return
		}
		respondError(w, http.StatusInternalServerError, "Failed to validate API key")
		return
	}

	// Parse notification ID
	id, err := parseNotificationIDFromURL(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Resolve the notification
	err = h.service.ResolveNotification(r.Context(), id)
	if err != nil {
		log.Printf("Error resolving notification %d: %v", id, err)

		if errors.Is(err, ErrNotificationNotFound) {
			respondError(w, http.StatusNotFound, "Notification not found")
			return
		}

		respondError(w, http.StatusInternalServerError, "Failed to resolve notification")
		return
	}

	// Return success
	respondJSON(w, http.StatusOK, map[string]string{
		"status":  "resolved",
		"message": "Notification resolved successfully",
	})
}

// DeleteNotification handles deleting a notification.
func (h *APIHandler) DeleteNotification(w http.ResponseWriter, r *http.Request) {
	// Validate API key
	_, err := h.validateAPIKey(r, "notifications:delete")
	if err != nil {
		if errors.Is(err, ErrInvalidAPIKey) || errors.Is(err, ErrAPIKeyExpired) {
			respondError(w, http.StatusUnauthorized, err.Error())
			return
		}
		if errors.Is(err, ErrPermissionDenied) {
			respondError(w, http.StatusForbidden, err.Error())
			return
		}
		respondError(w, http.StatusInternalServerError, "Failed to validate API key")
		return
	}

	// Parse notification ID
	id, err := parseNotificationIDFromURL(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Delete the notification
	err = h.service.DeleteNotification(r.Context(), id)
	if err != nil {
		log.Printf("Error deleting notification %d: %v", id, err)
		respondError(w, http.StatusInternalServerError, "Failed to delete notification")
		return
	}

	// Return success
	w.WriteHeader(http.StatusNoContent)
}

// HandleCallback handles webhook callbacks from notification targets.
func (h *APIHandler) HandleCallback(w http.ResponseWriter, r *http.Request) {
	// Extract target type and ID from URL
	vars := mux.Vars(r)
	targetTypeStr := vars["target_type"]
	targetID := vars["target_id"]

	targetType := TargetType(targetTypeStr)

	// Get the handler for this target type
	handler, ok := h.handlers[targetType]
	if !ok {
		respondError(w, http.StatusBadRequest, fmt.Sprintf("Unsupported target type: %s", targetType))
		return
	}

	// Read request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Failed to read request body")
		return
	}

	// Verify signature if it's a webhook
	if targetType == TargetTypeWebhook {
		if webhookHandler, ok := handler.(*WebhookHandler); ok {
			signature := r.Header.Get("X-Signature")
			if signature != "" && !webhookHandler.VerifySignature(body, signature, targetID) {
				respondError(w, http.StatusUnauthorized, "Invalid signature")
				return
			}
		}
	}

	// Parse acknowledgment
	ackReq, err := handler.ParseAcknowledgment(r.Context(), body)
	if err != nil {
		respondError(w, http.StatusBadRequest, fmt.Sprintf("Failed to parse acknowledgment: %v", err))
		return
	}

	// Acknowledge the notification
	err = h.service.AcknowledgeNotification(r.Context(), *ackReq)
	if err != nil {
		log.Printf("Error acknowledging notification %d from callback: %v", ackReq.NotificationID, err)

		if errors.Is(err, ErrNotificationNotFound) {
			respondError(w, http.StatusNotFound, "Notification not found")
			return
		}

		if errors.Is(err, ErrAlreadyAcknowledged) {
			respondJSON(w, http.StatusOK, map[string]string{
				"status":  "already_acknowledged",
				"message": "Notification was already acknowledged",
			})
			return
		}

		respondError(w, http.StatusInternalServerError, "Failed to acknowledge notification")
		return
	}

	// Return success
	respondJSON(w, http.StatusOK, map[string]string{
		"status":  "acknowledged",
		"message": "Notification acknowledged successfully",
	})
}

// CreateAPIKey handles creating a new API key.
func (h *APIHandler) CreateAPIKey(w http.ResponseWriter, r *http.Request) {
	// Validate existing API key - must have admin permission
	_, err := h.validateAPIKey(r, "notifications:admin")
	if err != nil {
		if errors.Is(err, ErrInvalidAPIKey) || errors.Is(err, ErrAPIKeyExpired) {
			respondError(w, http.StatusUnauthorized, err.Error())
			return
		}
		if errors.Is(err, ErrPermissionDenied) {
			respondError(w, http.StatusForbidden, err.Error())
			return
		}
		respondError(w, http.StatusInternalServerError, "Failed to validate API key")
		return
	}

	// Parse request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Failed to read request body")
		return
	}

	var req struct {
		ServiceName  string   `json:"service_name"`
		Permissions  []string `json:"permissions"`
		ExpiresAfter *string  `json:"expires_after,omitempty"`
	}

	if err := json.Unmarshal(body, &req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid JSON body")
		return
	}

	// Validate required fields
	if req.ServiceName == "" || len(req.Permissions) == 0 {
		respondError(w, http.StatusBadRequest, "Service name and permissions are required")
		return
	}

	// Parse expiration if provided
	var expiresAt *time.Time
	if req.ExpiresAfter != nil {
		duration, err := time.ParseDuration(*req.ExpiresAfter)
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid expiration duration format")
			return
		}

		expTime := time.Now().Add(duration)
		expiresAt = &expTime
	}

	// Create the API key
	key, keyValue, err := h.service.CreateAPIKey(r.Context(), req.ServiceName, req.Permissions, expiresAt)
	if err != nil {
		log.Printf("Error creating API key: %v", err)
		respondError(w, http.StatusInternalServerError, "Failed to create API key")
		return
	}

	// Return the created key
	respondJSON(w, http.StatusCreated, map[string]interface{}{
		"key_id":       key.KeyID,
		"key":          keyValue, // This is the only time the key value is returned
		"service_name": key.ServiceName,
		"permissions":  key.Permissions,
		"created_at":   key.CreatedAt,
		"expires_at":   key.ExpiresAt,
		"message":      "Store this key securely as it will not be shown again",
	})
}

// RevokeAPIKey handles revoking an API key.
func (h *APIHandler) RevokeAPIKey(w http.ResponseWriter, r *http.Request) {
	// Validate existing API key - must have admin permission
	_, err := h.validateAPIKey(r, "notifications:admin")
	if err != nil {
		if errors.Is(err, ErrInvalidAPIKey) || errors.Is(err, ErrAPIKeyExpired) {
			respondError(w, http.StatusUnauthorized, err.Error())
			return
		}
		if errors.Is(err, ErrPermissionDenied) {
			respondError(w, http.StatusForbidden, err.Error())
			return
		}
		respondError(w, http.StatusInternalServerError, "Failed to validate API key")
		return
	}

	// Get key ID from URL
	vars := mux.Vars(r)
	keyID := vars["key_id"]
	if keyID == "" {
		respondError(w, http.StatusBadRequest, "Key ID is required")
		return
	}

	// Revoke the key
	err = h.service.RevokeAPIKey(r.Context(), keyID)
	if err != nil {
		log.Printf("Error revoking API key %s: %v", keyID, err)
		respondError(w, http.StatusInternalServerError, "Failed to revoke API key")
		return
	}

	// Return success
	w.WriteHeader(http.StatusNoContent)
}
