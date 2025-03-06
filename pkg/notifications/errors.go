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

import "errors"

// Common errors that can be returned by notification operations
var (
	// ErrNotificationNotFound is returned when a notification is not found.
	ErrNotificationNotFound = errors.New("notification not found")

	// ErrInvalidRequest is returned when a request is invalid.
	ErrInvalidRequest = errors.New("invalid request")

	// ErrInvalidAPIKey is returned when an API key is invalid.
	ErrInvalidAPIKey = errors.New("invalid API key")

	// ErrAPIKeyExpired is returned when an API key has expired.
	ErrAPIKeyExpired = errors.New("API key expired")

	// ErrPermissionDenied is returned when an API key lacks the required permissions.
	ErrPermissionDenied = errors.New("permission denied")

	// ErrAlreadyAcknowledged is returned when trying to acknowledge an already acknowledged notification.
	ErrAlreadyAcknowledged = errors.New("notification already acknowledged")

	// ErrServiceUnavailable is returned when the notification service is unavailable.
	ErrServiceUnavailable = errors.New("notification service unavailable")

	// ErrDatabaseError is returned when a database operation fails.
	ErrDatabaseError = errors.New("database error")

	// ErrNoHandler is returned when no handler is registered for a target type.
	ErrNoHandler = errors.New("no handler registered for target type")

	// ErrHandlerError is returned when a handler fails to process a notification.
	ErrHandlerError = errors.New("handler failed to process notification")

	// ErrInvalidSignature is returned when a webhook signature is invalid.
	ErrInvalidSignature = errors.New("invalid signature")

	// ErrConfigurationError is returned when there's an issue with configuration.
	ErrConfigurationError = errors.New("configuration error")
)
