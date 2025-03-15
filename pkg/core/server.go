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

package core

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/carverauto/serviceradar/pkg/checker/snmp"
	"github.com/carverauto/serviceradar/pkg/core/alerts"
	"github.com/carverauto/serviceradar/pkg/core/api"
	"github.com/carverauto/serviceradar/pkg/db"
	"github.com/carverauto/serviceradar/pkg/metrics"
	"github.com/carverauto/serviceradar/pkg/models"
	"github.com/carverauto/serviceradar/proto"
)

const (
	shutdownTimeout          = 10 * time.Second
	oneDay                   = 24 * time.Hour
	oneWeek                  = 7 * oneDay
	serviceradarDirPerms     = 0700
	nodeHistoryLimit         = 1000
	nodeDiscoveryTimeout     = 30 * time.Second
	nodeNeverReportedTimeout = 30 * time.Second
	defaultDBPath            = "/var/lib/serviceradar/serviceradar.db"
	statusUnknown            = "unknown"
	sweepService             = "sweep"
	dailyCleanupInterval     = 24 * time.Hour
	monitorInterval          = 30 * time.Second
)

func NewServer(_ context.Context, config *Config) (*Server, error) {
	if config.Metrics.Retention == 0 {
		config.Metrics.Retention = 100
	}

	if config.Metrics.MaxNodes == 0 {
		config.Metrics.MaxNodes = 10000
	}

	metricsManager := metrics.NewManager(models.MetricsConfig{
		Enabled:   config.Metrics.Enabled,
		Retention: config.Metrics.Retention,
		MaxNodes:  config.Metrics.MaxNodes,
	})

	// Use default DB path if not specified
	dbPath := config.DBPath
	if dbPath == "" {
		dbPath = defaultDBPath
	}

	// Ensure the directory exists
	if err := os.MkdirAll("/var/lib/serviceradar", serviceradarDirPerms); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %w", err)
	}

	// Initialize database
	database, err := db.New(dbPath)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", errDatabaseError, err)
	}

	server := &Server{
		db:             database,
		alertThreshold: config.AlertThreshold,
		webhooks:       make([]alerts.AlertService, 0),
		ShutdownChan:   make(chan struct{}),
		pollerPatterns: config.PollerPatterns,
		metrics:        metricsManager,
		snmpManager:    snmp.NewSNMPManager(database),
		config:         config,
	}

	// Initialize webhooks
	server.initializeWebhooks(config.Webhooks)

	return server, nil
}

func (s *Server) initializeWebhooks(configs []alerts.WebhookConfig) {
	for i, config := range configs {
		log.Printf("Processing webhook config %d: enabled=%v", i, config.Enabled)

		if config.Enabled {
			alerter := alerts.NewWebhookAlerter(config)
			s.webhooks = append(s.webhooks, alerter)

			log.Printf("Added webhook alerter: %+v", config.URL)
		}
	}
}

// Start implements the lifecycle.Service interface.

func (s *Server) Start(ctx context.Context) error {
	log.Printf("Starting core service...")

	// Clean up any unknown pollers first
	if err := s.cleanupUnknownPollers(ctx); err != nil {
		log.Printf("Warning: Failed to clean up unknown pollers: %v", err)
	}

	if s.grpcServer != nil {
		errCh := make(chan error, 1)

		go func() {
			if err := s.grpcServer.Start(); err != nil {
				select {
				case errCh <- err:
				default:
					log.Printf("gRPC server error: %v", err)
				}
			}
		}()
	}

	if err := s.sendStartupNotification(ctx); err != nil {
		log.Printf("Failed to send startup notification: %v", err)
	}

	go s.periodicCleanup(ctx)

	go s.runMetricsCleanup(ctx)

	go s.monitorNodes(ctx)

	return nil
}

// Stop gracefully shuts down the server.
func (s *Server) Stop(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, shutdownTimeout)
	defer cancel()

	// Send shutdown notification
	if err := s.sendShutdownNotification(ctx); err != nil {
		log.Printf("Failed to send shutdown notification: %v", err)
	}

	// Stop GRPC server if it exists
	if s.grpcServer != nil {
		// Stop no longer returns an error, just call it
		s.grpcServer.Stop(ctx)
	}

	// Close database
	if err := s.db.Close(); err != nil {
		log.Printf("Error closing database: %v", err)
	}

	// Signal all background tasks to stop
	close(s.ShutdownChan)

	return nil
}

// monitorNodes runs the main node monitoring loop.
func (s *Server) monitorNodes(ctx context.Context) {
	log.Printf("Starting node monitoring...")

	time.Sleep(nodeDiscoveryTimeout)

	// Initial checks
	s.checkInitialStates()

	time.Sleep(nodeNeverReportedTimeout)
	s.checkNeverReportedPollers(ctx)

	// Start monitoring loop
	s.MonitorPollers(ctx)
}

func (s *Server) GetMetricsManager() metrics.MetricCollector {
	return s.metrics
}

func (s *Server) GetSNMPManager() snmp.SNMPManager {
	return s.snmpManager
}

func (s *Server) runMetricsCleanup(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if s.metrics != nil {
				if manager, ok := s.metrics.(*metrics.Manager); ok {
					manager.CleanupStaleNodes(oneWeek)
				} else {
					log.Printf("Error: s.metrics is not of type *metrics.Manager")
				}
			}
		}
	}
}

func (s *Server) isKnownPoller(pollerID string) bool {
	// Use the config values directly where needed
	for _, known := range s.config.KnownPollers {
		if known == pollerID {
			return true
		}
	}

	return false
}

func (s *Server) cleanupUnknownPollers(ctx context.Context) error {
	if len(s.config.KnownPollers) == 0 {
		return nil // No filtering if no known pollers specified
	}

	// set a timer on the context to ensure we don't run indefinitely
	_, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()

	// Build a query with placeholders for known pollers
	placeholders := make([]string, len(s.config.KnownPollers))

	args := make([]interface{}, len(s.config.KnownPollers))

	for i, poller := range s.config.KnownPollers {
		placeholders[i] = "?"
		args[i] = poller
	}

	// Delete all nodes not in known_pollers
	query := fmt.Sprintf("DELETE FROM nodes WHERE node_id NOT IN (%s)",
		strings.Join(placeholders, ","))

	result, err := s.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to cleanup unknown pollers: %w", err)
	}

	if rows, err := result.RowsAffected(); err == nil && rows > 0 {
		log.Printf("Cleaned up %d unknown poller(s) from database", rows)
	}

	return nil
}

func (s *Server) sendStartupNotification(ctx context.Context) error {
	if len(s.webhooks) == 0 {
		return nil
	}

	alert := &alerts.WebhookAlert{
		Level:     alerts.Info,
		Title:     "Core Service Started",
		Message:   fmt.Sprintf("ServiceRadar core service initialized at %s", time.Now().Format(time.RFC3339)),
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		NodeID:    "core",
		Details: map[string]any{
			"version":  "1.0.25",
			"hostname": getHostname(),
		},
	}

	return s.sendAlert(ctx, alert)
}

func (s *Server) sendShutdownNotification(ctx context.Context) error {
	if len(s.webhooks) == 0 {
		return nil
	}

	alert := &alerts.WebhookAlert{
		Level:     alerts.Warning,
		Title:     "Core Service Stopping",
		Message:   fmt.Sprintf("ServiceRadar core service shutting down at %s", time.Now().Format(time.RFC3339)),
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		NodeID:    "core",
		Details: map[string]any{
			"hostname": getHostname(),
		},
	}

	return s.sendAlert(ctx, alert)
}

func (s *Server) Shutdown(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, shutdownTimeout)
	defer cancel()

	if err := s.db.Close(); err != nil {
		log.Printf("Error closing database: %v", err)
	}

	if len(s.webhooks) > 0 {
		alert := alerts.WebhookAlert{
			Level:     alerts.Warning,
			Title:     "Core Service Stopping",
			Message:   fmt.Sprintf("ServiceRadar core service shutting down at %s", time.Now().Format(time.RFC3339)),
			Timestamp: time.Now().UTC().Format(time.RFC3339),
			NodeID:    "core",
			Details: map[string]any{
				"hostname": getHostname(),
				"pid":      os.Getpid(),
			},
		}

		err := s.sendAlert(ctx, &alert)
		if err != nil {
			log.Printf("Error sending shutdown alert: %v", err)

			return
		}
	}

	close(s.ShutdownChan)
}

func (s *Server) SetAPIServer(apiServer api.Service) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.apiServer = apiServer
	apiServer.SetKnownPollers(s.config.KnownPollers)

	apiServer.SetNodeHistoryHandler(func(nodeID string) ([]api.NodeHistoryPoint, error) {
		points, err := s.db.GetNodeHistoryPoints(nodeID, nodeHistoryLimit)
		if err != nil {
			return nil, fmt.Errorf("failed to get node history: %w", err)
		}

		apiPoints := make([]api.NodeHistoryPoint, len(points))
		for i, p := range points {
			apiPoints[i] = api.NodeHistoryPoint{
				Timestamp: p.Timestamp,
				IsHealthy: p.IsHealthy,
			}
		}

		return apiPoints, nil
	})
}

func (s *Server) checkInitialStates() {
	likeConditions := make([]string, 0, len(s.pollerPatterns))
	args := make([]interface{}, 0, len(s.pollerPatterns))

	// Construct the WHERE clause with multiple LIKE conditions
	for _, pattern := range s.pollerPatterns {
		likeConditions = append(likeConditions, "node_id LIKE ?")
		args = append(args, pattern)
	}

	// Base query without WHERE clause
	query := `
        SELECT node_id, is_healthy, last_seen 
        FROM nodes 
    `

	// Add WHERE clause only if there are conditions
	if len(likeConditions) > 0 {
		query += fmt.Sprintf("WHERE %s ", strings.Join(likeConditions, " OR "))
	}

	// Add ORDER BY clause
	query += "ORDER BY last_seen DESC"

	rows, err := s.db.Query(query, args...)
	if err != nil {
		log.Printf("Error querying nodes: %v", err)

		return
	}
	defer db.CloseRows(rows)

	for rows.Next() {
		var nodeID string

		var isHealthy bool

		var lastSeen time.Time

		if err := rows.Scan(&nodeID, &isHealthy, &lastSeen); err != nil {
			log.Printf("Error scanning node row: %v", err)
			continue
		}

		duration := time.Since(lastSeen)
		if duration > s.alertThreshold {
			log.Printf("Node %s found offline during initial check (last seen: %v ago)",
				nodeID, duration.Round(time.Second))
		}
	}
}

// updateAPIState updates the API server with the latest node status.
func (s *Server) updateAPIState(pollerID string, apiStatus *api.NodeStatus) {
	if s.apiServer == nil {
		log.Printf("Warning: API server not initialized, state not updated")

		return
	}

	s.apiServer.UpdateNodeStatus(pollerID, apiStatus)

	log.Printf("Updated API server state for node: %s", pollerID)
}

// getNodeHealthState retrieves the current health state of a node.
func (s *Server) getNodeHealthState(pollerID string) (bool, error) {
	var currentState bool

	err := s.db.QueryRow("SELECT is_healthy FROM nodes WHERE node_id = ?", pollerID).Scan(&currentState)

	return currentState, err
}

func (s *Server) processStatusReport(
	ctx context.Context, req *proto.PollerStatusRequest, now time.Time) (*api.NodeStatus, error) {
	currentState, err := s.getNodeHealthState(req.PollerId)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Printf("Error checking node state: %v", err)
	}

	apiStatus := s.createNodeStatus(req, now)

	s.processServices(req.PollerId, apiStatus, req.Services, now)

	if err := s.updateNodeState(ctx, req.PollerId, apiStatus, currentState, now); err != nil {
		return nil, err
	}

	return apiStatus, nil
}

func (*Server) createNodeStatus(req *proto.PollerStatusRequest, now time.Time) *api.NodeStatus {
	return &api.NodeStatus{
		NodeID:     req.PollerId,
		LastUpdate: now,
		IsHealthy:  true,
		Services:   make([]api.ServiceStatus, 0, len(req.Services)),
	}
}

func (s *Server) processServices(pollerID string, apiStatus *api.NodeStatus, services []*proto.ServiceStatus, now time.Time) {
	allServicesAvailable := true

	for _, svc := range services {
		log.Printf("Processing service %s for node %s", svc.ServiceName, pollerID)
		log.Printf("Service type: %s, Message length: %d", svc.ServiceType, len(svc.Message))

		apiService := api.ServiceStatus{
			Name:      svc.ServiceName,
			Type:      svc.ServiceType,
			Available: svc.Available,
			Message:   svc.Message,
		}

		if !svc.Available {
			allServicesAvailable = false
		}

		if svc.Message == "" {
			log.Printf("No message content for service %s", svc.ServiceName)

			if err := s.handleService(pollerID, &apiService, now); err != nil {
				log.Printf("Error handling service %s: %v", svc.ServiceName, err)
			}

			apiStatus.Services = append(apiStatus.Services, apiService)

			continue // Skip to the next service
		}

		var details json.RawMessage // Declare details here, outside the if block

		if err := json.Unmarshal([]byte(svc.Message), &details); err != nil {
			log.Printf("Error unmarshaling service details for %s: %v", svc.ServiceName, err)
			log.Printf("Raw message: %s", svc.Message)

			if err := s.handleService(pollerID, &apiService, now); err != nil {
				log.Printf("Error handling service %s: %v", svc.ServiceName, err)
			}

			apiStatus.Services = append(apiStatus.Services, apiService)

			continue // Skip to the next service
		}

		apiService.Details = details // Now details is in scope

		if svc.ServiceType == "snmp" {
			log.Printf("Found SNMP service, attempting to process metrics for node %s", pollerID)

			if err := s.processSNMPMetrics(pollerID, details, now); err != nil { // details is also available here
				log.Printf("Error processing SNMP metrics for node %s: %v", pollerID, err)
			}
		}

		if err := s.handleService(pollerID, &apiService, now); err != nil {
			log.Printf("Error handling service %s: %v", svc.ServiceName, err)
		}

		apiStatus.Services = append(apiStatus.Services, apiService)
	}

	apiStatus.IsHealthy = allServicesAvailable
}

// processSNMPMetrics extracts and stores SNMP metrics from service details.
func (s *Server) processSNMPMetrics(nodeID string, details json.RawMessage, timestamp time.Time) error {
	log.Printf("Processing SNMP metrics for node %s", nodeID)

	// Parse the outer structure which contains target-specific data
	var snmpData map[string]struct {
		Available bool                     `json:"available"`
		LastPoll  string                   `json:"last_poll"`
		OIDStatus map[string]OIDStatusData `json:"oid_status"`
	}

	if err := json.Unmarshal(details, &snmpData); err != nil {
		return fmt.Errorf("failed to parse SNMP data: %w", err)
	}

	// Process each target's data
	for targetName, targetData := range snmpData {
		log.Printf("Processing target %s with %d OIDs", targetName, len(targetData.OIDStatus))

		// Process each OID's data
		for oidName, oidStatus := range targetData.OIDStatus {
			// Create metadata
			metadata := map[string]interface{}{
				"target_name": targetName,
				"last_poll":   targetData.LastPoll,
			}

			// Convert the value to string for storage
			valueStr := fmt.Sprintf("%v", oidStatus.LastValue)

			// Create metric
			metric := &db.TimeseriesMetric{
				Name:      oidName,
				Value:     valueStr,
				Type:      "snmp",
				Timestamp: timestamp,
				Metadata:  metadata,
			}

			// Store in database
			if err := s.db.StoreMetric(nodeID, metric); err != nil {
				log.Printf("Error storing SNMP metric %s for node %s: %v", oidName, nodeID, err)

				continue
			}
		}
	}

	return nil
}

func (s *Server) handleService(pollerID string, svc *api.ServiceStatus, now time.Time) error {
	if svc.Type == sweepService {
		if err := s.processSweepData(svc, now); err != nil {
			return fmt.Errorf("failed to process sweep data: %w", err)
		}
	}

	return s.saveServiceStatus(pollerID, svc, now)
}

func (*Server) processSweepData(svc *api.ServiceStatus, now time.Time) error {
	var sweepData proto.SweepServiceStatus

	if err := json.Unmarshal([]byte(svc.Message), &sweepData); err != nil {
		return fmt.Errorf("%w: %w", errInvalidSweepData, err)
	}

	// If LastSweep is not set or is invalid (0 or negative), use current time
	if sweepData.LastSweep > now.Add(oneDay).Unix() {
		log.Printf("Invalid or missing LastSweep timestamp (%d), using current time", sweepData.LastSweep)
		sweepData.LastSweep = now.Unix()

		// Update the message with corrected timestamp
		updatedData := proto.SweepServiceStatus{
			Network:        sweepData.Network,
			TotalHosts:     sweepData.TotalHosts,
			AvailableHosts: sweepData.AvailableHosts,
			LastSweep:      now.Unix(),
		}

		updatedMessage, err := json.Marshal(&updatedData)
		if err != nil {
			return fmt.Errorf("failed to marshal updated sweep data: %w", err)
		}

		svc.Message = string(updatedMessage)

		log.Printf("Updated sweep data with current timestamp: %v", now.Format(time.RFC3339))
	} else {
		// Log the existing timestamp for debugging
		log.Printf("Processing sweep data with timestamp: %v",
			time.Unix(sweepData.LastSweep, 0).Format(time.RFC3339))
	}

	return nil
}

func (s *Server) saveServiceStatus(pollerID string, svc *api.ServiceStatus, now time.Time) error {
	status := &db.ServiceStatus{
		NodeID:      pollerID,
		ServiceName: svc.Name,
		ServiceType: svc.Type,
		Available:   svc.Available,
		Details:     svc.Message,
		Timestamp:   now,
	}

	if err := s.db.UpdateServiceStatus(status); err != nil {
		return fmt.Errorf("%w: failed to update service status", errDatabaseError)
	}

	return nil
}

// storeNodeStatus updates the node status in the database.
func (s *Server) storeNodeStatus(pollerID string, isHealthy bool, now time.Time) error {
	nodeStatus := &db.NodeStatus{
		NodeID:    pollerID,
		IsHealthy: isHealthy,
		LastSeen:  now,
	}

	if err := s.db.UpdateNodeStatus(nodeStatus); err != nil {
		return fmt.Errorf("failed to store node status: %w", err)
	}

	return nil
}

func (s *Server) updateNodeState(ctx context.Context, pollerID string, apiStatus *api.NodeStatus, wasHealthy bool, now time.Time) error {
	if err := s.storeNodeStatus(pollerID, apiStatus.IsHealthy, now); err != nil {
		return err
	}

	// Check for recovery
	if !wasHealthy && apiStatus.IsHealthy {
		s.handleNodeRecovery(ctx, pollerID, apiStatus, now)
	}

	return nil
}

// periodicCleanup runs regular maintenance tasks on the database.
func (s *Server) periodicCleanup(_ context.Context) {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-s.ShutdownChan:
			return
		case <-ticker.C:
			// Clean up old data (keep 7 days by default)
			if err := s.db.CleanOldData(7 * 24 * time.Hour); err != nil {
				log.Printf("Error during periodic cleanup: %v", err)
			}

			// Vacuum the database every 24 hours to reclaim space
			if time.Now().Hour() == 0 { // Run at midnight
				if _, err := s.db.Exec("VACUUM"); err != nil {
					log.Printf("Error vacuuming database: %v", err)
				}
			}
		}
	}
}

// checkNeverReportedNodes checks for and alerts on nodes that exist but have never reported.
func (s *Server) checkNeverReportedNodes(ctx context.Context) error {
	conditions := make([]string, 0, len(s.pollerPatterns))
	args := make([]interface{}, 0, len(s.pollerPatterns))

	// Build LIKE conditions for each pattern
	for _, pattern := range s.pollerPatterns {
		conditions = append(conditions, "node_id LIKE ?")
		args = append(args, pattern)
	}

	// Construct query with LIKE conditions
	query := `SELECT node_id FROM nodes WHERE last_seen = ''`
	if len(conditions) > 0 {
		query += " AND (" + strings.Join(conditions, " OR ") + ")"
	}

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return fmt.Errorf("error querying unreported nodes: %w", err)
	}
	defer db.CloseRows(rows)

	var unreportedNodes []string

	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			log.Printf("Error scanning node ID: %v", err)
			continue
		}

		unreportedNodes = append(unreportedNodes, id)
	}

	// Check for any errors encountered during iteration
	if err := rows.Err(); err != nil {
		return fmt.Errorf("error iterating rows: %w", err)
	}

	if len(unreportedNodes) > 0 {
		alert := &alerts.WebhookAlert{
			Level:     alerts.Warning,
			Title:     "Pollers Never Reported",
			Message:   fmt.Sprintf("%d poller(s) have not reported since startup", len(unreportedNodes)),
			NodeID:    "core",
			Timestamp: time.Now().UTC().Format(time.RFC3339),
			Details: map[string]any{
				"hostname":     getHostname(),
				"poller_ids":   unreportedNodes,
				"poller_count": len(unreportedNodes),
			},
		}

		if err := s.sendAlert(ctx, alert); err != nil {
			log.Printf("Error sending unreported nodes alert: %v", err)
			return err
		}
	}

	return nil
}

func (s *Server) checkNeverReportedPollers(ctx context.Context) {
	log.Printf("Checking for unreported nodes matching patterns: %v", s.pollerPatterns)

	// Build SQL pattern for REGEXP
	combinedPattern := strings.Join(s.pollerPatterns, "|")
	if combinedPattern == "" {
		return
	}

	var unreportedNodes []string

	rows, err := s.db.Query(`
        SELECT node_id 
        FROM nodes 
        WHERE node_id REGEXP ? AND last_seen = ''`,
		combinedPattern)
	if err != nil {
		log.Printf("Error querying unreported nodes: %v", err)
		return
	}
	defer db.CloseRows(rows)

	for rows.Next() {
		var nodeID string
		if err := rows.Scan(&nodeID); err != nil {
			log.Printf("Error scanning node ID: %v", err)
			continue
		}

		unreportedNodes = append(unreportedNodes, nodeID)
	}

	if len(unreportedNodes) > 0 {
		s.sendUnreportedNodesAlert(ctx, unreportedNodes)
	}
}

func (s *Server) sendUnreportedNodesAlert(ctx context.Context, nodeIDs []string) {
	alert := &alerts.WebhookAlert{
		Level:     alerts.Warning,
		Title:     "Pollers Never Reported",
		Message:   fmt.Sprintf("%d poller(s) have not reported since startup", len(nodeIDs)),
		NodeID:    "core",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Details: map[string]any{
			"hostname":     getHostname(),
			"poller_ids":   nodeIDs,
			"poller_count": len(nodeIDs),
		},
	}

	if err := s.sendAlert(ctx, alert); err != nil {
		log.Printf("Error sending unreported nodes alert: %v", err)
	}
}

func (c *Config) UnmarshalJSON(data []byte) error {
	type Alias Config

	aux := &struct {
		AlertThreshold string `json:"alert_threshold"`
		*Alias
	}{
		Alias: (*Alias)(c),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Parse the alert threshold
	if aux.AlertThreshold != "" {
		duration, err := time.ParseDuration(aux.AlertThreshold)
		if err != nil {
			return fmt.Errorf("invalid alert threshold format: %w", err)
		}

		c.AlertThreshold = duration
	}

	return nil
}

func LoadConfig(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("failed to read config: %w", err)
	}

	var config Config

	if err := json.Unmarshal(data, &config); err != nil {
		return Config{}, fmt.Errorf("failed to parse config: %w", err)
	}

	if config.Security != nil {
		log.Printf("Security config: Mode=%s, CertDir=%s, Role=%s",
			config.Security.Mode, config.Security.CertDir, config.Security.Role)
	}

	return config, nil
}

func (s *Server) MonitorPollers(ctx context.Context) {
	ticker := time.NewTicker(monitorInterval) // Check every 30 seconds
	defer ticker.Stop()

	cleanupTicker := time.NewTicker(dailyCleanupInterval)
	defer cleanupTicker.Stop()

	// Initial checks
	if err := s.checkNodeStates(ctx); err != nil {
		log.Printf("Initial state check failed: %v", err)
	}

	if err := s.checkNeverReportedNodes(ctx); err != nil {
		log.Printf("Initial never-reported check failed: %v", err)
	}

	for {
		select {
		case <-ctx.Done():
			return
		case <-s.ShutdownChan:
			return
		case <-ticker.C:
			s.handleMonitorTick(ctx)
		case <-cleanupTicker.C:
			s.handleCleanupTick()
		}
	}
}

// handleMonitorTick handles the logic for the monitor ticker.
func (s *Server) handleMonitorTick(ctx context.Context) {
	if err := s.checkNodeStates(ctx); err != nil {
		log.Printf("Node state check failed: %v", err)
	}

	if err := s.checkNeverReportedNodes(ctx); err != nil {
		log.Printf("Never-reported check failed: %v", err)
	}
}

// handleCleanupTick handles the logic for the cleanup ticker.
func (s *Server) handleCleanupTick() {
	if err := s.performDailyCleanup(); err != nil {
		log.Printf("Daily cleanup failed: %v", err)
	}
}

// performDailyCleanup performs the daily cleanup task.
func (s *Server) performDailyCleanup() error {
	log.Println("Performing daily cleanup...")

	if err := s.db.CleanOldData(oneWeek); err != nil {
		log.Printf("Error cleaning old data: %v", err)

		return err
	}

	return nil
}

func (s *Server) checkNodeStates(ctx context.Context) error {
	// Pre-allocate slices
	conditions := make([]string, 0, len(s.pollerPatterns))
	args := make([]interface{}, 0, len(s.pollerPatterns))

	// Build LIKE conditions for each pattern
	for _, pattern := range s.pollerPatterns {
		conditions = append(conditions, "node_id LIKE ?")
		args = append(args, pattern)
	}

	// Construct query with LIKE conditions
	query := `SELECT node_id, last_seen, is_healthy FROM nodes`
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " OR ")
	}

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return fmt.Errorf("failed to query nodes: %w", err)
	}
	defer db.CloseRows(rows)

	threshold := time.Now().Add(-s.alertThreshold)

	for rows.Next() {
		var nodeID string

		var lastSeen time.Time

		var isHealthy bool

		if err := rows.Scan(&nodeID, &lastSeen, &isHealthy); err != nil {
			log.Printf("Error scanning node row: %v", err)

			continue
		}

		err := s.evaluateNodeHealth(ctx, nodeID, lastSeen, isHealthy, threshold)
		if err != nil {
			// Only log errors, don't propagate service-related issues
			log.Printf("Error evaluating node %s health: %v", nodeID, err)
		}
	}

	return rows.Err()
}

func (s *Server) evaluateNodeHealth(
	ctx context.Context, nodeID string, lastSeen time.Time, isHealthy bool, threshold time.Time) error {
	log.Printf("Evaluating node health: id=%s lastSeen=%v isHealthy=%v threshold=%v",
		nodeID, lastSeen.Format(time.RFC3339), isHealthy, threshold.Format(time.RFC3339))

	// Case 1: Node was healthy but hasn't been seen recently (went down)
	if isHealthy && lastSeen.Before(threshold) {
		duration := time.Since(lastSeen).Round(time.Second)
		log.Printf("Node %s appears to be offline (last seen: %v ago)", nodeID, duration)

		return s.handleNodeDown(ctx, nodeID, lastSeen)
	}

	// Case 2: Node is healthy and reporting within threshold - DO NOTHING
	if isHealthy && !lastSeen.Before(threshold) {
		return nil
	}

	// Case 3: Node is reporting but its status might have changed
	if !lastSeen.Before(threshold) {
		// Get the current health status
		currentHealth, err := s.getNodeHealthState(nodeID)
		if err != nil {
			log.Printf("Error getting current health state for node %s: %v", nodeID, err)

			return fmt.Errorf("failed to get current health state: %w", err)
		}

		// ONLY handle potential recovery - do not send service alerts here
		if !isHealthy && currentHealth {
			return s.handlePotentialRecovery(ctx, nodeID, lastSeen)
		}
	}

	return nil
}

func (s *Server) handlePotentialRecovery(ctx context.Context, nodeID string, lastSeen time.Time) error {
	apiStatus := &api.NodeStatus{
		NodeID:     nodeID,
		LastUpdate: lastSeen,
		Services:   make([]api.ServiceStatus, 0),
	}

	s.handleNodeRecovery(ctx, nodeID, apiStatus, lastSeen)

	return nil
}

func (s *Server) handleNodeDown(ctx context.Context, nodeID string, lastSeen time.Time) error {
	if err := s.updateNodeStatus(nodeID, false, lastSeen); err != nil {
		return fmt.Errorf("failed to update node status: %w", err)
	}

	// Send alert
	alert := &alerts.WebhookAlert{
		Level:     alerts.Error,
		Title:     "Node Offline",
		Message:   fmt.Sprintf("Node '%s' is offline", nodeID),
		NodeID:    nodeID,
		Timestamp: lastSeen.UTC().Format(time.RFC3339),
		Details: map[string]any{
			"hostname": getHostname(),
			"duration": time.Since(lastSeen).String(),
		},
	}

	if err := s.sendAlert(ctx, alert); err != nil {
		log.Printf("Failed to send down alert: %v", err)
	}

	// Update API state
	if s.apiServer != nil {
		s.apiServer.UpdateNodeStatus(nodeID, &api.NodeStatus{
			NodeID:     nodeID,
			IsHealthy:  false,
			LastUpdate: lastSeen,
		})
	}

	return nil
}
func (s *Server) updateNodeStatus(nodeID string, isHealthy bool, timestamp time.Time) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func(tx db.Transaction) {
		err = tx.Rollback()
		if err != nil {
			log.Printf("Error rolling back transaction: %v", err)
		}
	}(tx)

	sqlTx, err := db.ToTx(tx)
	if err != nil {
		return fmt.Errorf("invalid transaction: %w", err)
	}

	// Update node status
	if err := s.updateNodeInTx(sqlTx, nodeID, isHealthy, timestamp); err != nil {
		return err
	}

	// Add history entry
	if _, err := tx.Exec(`
        INSERT INTO node_history (node_id, timestamp, is_healthy)
        VALUES (?, ?, ?)
    `, nodeID, timestamp, isHealthy); err != nil {
		return fmt.Errorf("failed to insert history: %w", err)
	}

	return tx.Commit()
}

func (*Server) updateNodeInTx(tx *sql.Tx, nodeID string, isHealthy bool, timestamp time.Time) error {
	// Check if node exists
	var exists bool
	if err := tx.QueryRow("SELECT EXISTS(SELECT 1 FROM nodes WHERE node_id = ?)", nodeID).Scan(&exists); err != nil {
		return fmt.Errorf("failed to check node existence: %w", err)
	}

	if exists {
		_, err := tx.Exec(`
            UPDATE nodes 
            SET is_healthy = ?,
                last_seen = ?
            WHERE node_id = ?
        `, isHealthy, timestamp, nodeID)

		return err
	}

	// Insert new node
	_, err := tx.Exec(`
        INSERT INTO nodes (node_id, first_seen, last_seen, is_healthy)
        VALUES (?, ?, ?, ?)
    `, nodeID, timestamp, timestamp, isHealthy)

	return err
}

func (s *Server) handleNodeRecovery(ctx context.Context, nodeID string, apiStatus *api.NodeStatus, timestamp time.Time) {
	// Reset the "down" state in the alerter *before* sending the alert.
	for _, webhook := range s.webhooks {
		if alerter, ok := webhook.(*alerts.WebhookAlerter); ok {
			alerter.MarkNodeAsRecovered(nodeID)
			alerter.MarkServiceAsRecovered(nodeID)
		}
	}

	alert := &alerts.WebhookAlert{
		Level:       alerts.Info,
		Title:       "Node Recovered",
		Message:     fmt.Sprintf("Node '%s' is back online", nodeID),
		NodeID:      nodeID,
		Timestamp:   timestamp.UTC().Format(time.RFC3339),
		ServiceName: "", // Ensure ServiceName is empty for node-level alerts
		Details: map[string]any{
			"hostname":      getHostname(),
			"recovery_time": timestamp.Format(time.RFC3339),
			"services":      len(apiStatus.Services), //  This might be 0, which is fine.
		},
	}

	if err := s.sendAlert(ctx, alert); err != nil {
		log.Printf("Failed to send recovery alert: %v", err)
	}
}

func (s *Server) sendAlert(ctx context.Context, alert *alerts.WebhookAlert) error {
	var errs []error

	log.Printf("Sending alert: %s", alert.Message)

	for _, webhook := range s.webhooks {
		if err := webhook.Alert(ctx, alert); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("%w: %v", errFailedToSendAlerts, errs)
	}

	return nil
}

// ReportStatus implements the PollerServiceServer interface. It processes status reports from pollers.
func (s *Server) ReportStatus(ctx context.Context, req *proto.PollerStatusRequest) (*proto.PollerStatusResponse, error) {
	if req.PollerId == "" {
		return nil, errEmptyPollerID
	}

	// Add check for known pollers
	if !s.isKnownPoller(req.PollerId) {
		log.Printf("Ignoring status report from unknown poller: %s", req.PollerId)

		return &proto.PollerStatusResponse{Received: true}, nil
	}

	now := time.Unix(req.Timestamp, 0)
	timestamp := time.Now()
	responseTime := timestamp.Sub(now).Nanoseconds()

	log.Printf("Response time for %s: %d ns (%.2f ms)",
		req.PollerId,
		responseTime,
		float64(responseTime)/float64(time.Millisecond))

	apiStatus, err := s.processStatusReport(ctx, req, now)
	if err != nil {
		return nil, fmt.Errorf("failed to process status report: %w", err)
	}

	if s.metrics != nil {
		for _, service := range req.Services {
			if service.ServiceType != "icmp" {
				continue
			}

			// Parse the ping response
			var pingResult struct {
				Host         string  `json:"host"`
				ResponseTime int64   `json:"response_time"`
				PacketLoss   float64 `json:"packet_loss"`
				Available    bool    `json:"available"`
			}

			if err := json.Unmarshal([]byte(service.Message), &pingResult); err != nil {
				log.Printf("Failed to parse ICMP response for service %s: %v", service.ServiceName, err)

				continue
			}

			// Add metric with the actual response time
			if err := s.metrics.AddMetric(
				req.PollerId,
				time.Now(),
				pingResult.ResponseTime,
				service.ServiceName,
			); err != nil {
				log.Printf("Failed to add ICMP metric for %s: %v", service.ServiceName, err)

				continue
			}

			log.Printf("Added ICMP metric for %s: time=%v response_time=%.2fms",
				service.ServiceName,
				time.Now().Format(time.RFC3339),
				float64(pingResult.ResponseTime)/float64(time.Millisecond))
		}
	}

	s.updateAPIState(req.PollerId, apiStatus)

	return &proto.PollerStatusResponse{Received: true}, nil
}

func getHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		return statusUnknown
	}

	return hostname
}
