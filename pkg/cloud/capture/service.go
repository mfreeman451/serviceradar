// Package capture pkg/cloud/capture/service.go
package capture

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/mfreeman451/serviceradar/pkg/bridge"
	"github.com/mfreeman451/serviceradar/pkg/grpc"
	"github.com/mfreeman451/serviceradar/pkg/metrics"
	"github.com/mfreeman451/serviceradar/pkg/tunnel"
	"github.com/mfreeman451/serviceradar/proto"
)

type CaptureStatus struct {
	NodeID      string
	Interface   string
	StartTime   time.Time
	PacketsRead uint64
	BytesRead   uint64
	IsActive    bool
	Error       string
}

type CaptureService struct {
	mu         sync.RWMutex
	tunnels    map[string]tunnel.Tunnel
	bridges    map[string]bridge.Bridge
	pollerConn map[string]*grpc.ClientConn
	status     map[string]*CaptureStatus
	metrics    metrics.MetricCollector
}

func NewCaptureService(metrics metrics.MetricCollector) *CaptureService {
	return &CaptureService{
		tunnels:    make(map[string]tunnel.Tunnel),
		bridges:    make(map[string]bridge.Bridge),
		pollerConn: make(map[string]*grpc.ClientConn),
		status:     make(map[string]*CaptureStatus),
		metrics:    metrics,
	}
}

// RegisterPoller registers a poller's connection.
func (s *CaptureService) RegisterPoller(pollerID string, conn *grpc.ClientConn) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.pollerConn[pollerID] = conn

	// Initialize status for this poller
	s.status[pollerID] = &CaptureStatus{
		NodeID:    pollerID,
		IsActive:  false,
		StartTime: time.Now(),
	}
}

// CleanupPoller removes a poller and stops any active captures.
func (s *CaptureService) CleanupPoller(pollerID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Stop any active capture
	if t, exists := s.tunnels[pollerID]; exists {
		t.Close()
		delete(s.tunnels, pollerID)
	}

	// Cleanup bridge
	if b, exists := s.bridges[pollerID]; exists {
		b.Close()
		delete(s.bridges, pollerID)
	}

	// Cleanup status
	delete(s.status, pollerID)
	delete(s.pollerConn, pollerID)

	// Record metric
	if s.metrics != nil {
		s.metrics.AddMetric(pollerID, time.Now(), 0, "capture_cleanup")
	}
}

// ListActiveCaptures returns all active captures.
func (s *CaptureService) ListActiveCaptures() []*CaptureStatus {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var active []*CaptureStatus

	for _, status := range s.status {
		if status.IsActive {
			active = append(active, status)
		}
	}

	return active
}

// GetCaptureStatus returns status for a specific capture.
func (s *CaptureService) GetCaptureStatus(nodeID string) (*CaptureStatus, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	status, exists := s.status[nodeID]
	if !exists {
		return nil, fmt.Errorf("no capture status found for node %s", nodeID)
	}

	return status, nil
}

// StartCapture starts a remote packet capture session.
func (s *CaptureService) StartCapture(ctx context.Context, nodeID, iface string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Update status
	status, exists := s.status[nodeID]
	if !exists {
		status = &CaptureStatus{
			NodeID:    nodeID,
			StartTime: time.Now(),
		}
		s.status[nodeID] = status
	}

	// Get the poller connection
	pollerConn, exists := s.pollerConn[nodeID]
	if !exists {
		status.Error = "no connection available"

		return fmt.Errorf("no connection available for node %s", nodeID)
	}

	// Start a bidirectional stream
	client := proto.NewPollerServiceClient(pollerConn.GetConnection())
	stream, err := client.ForwardPackets(ctx)
	if err != nil {
		status.Error = fmt.Sprintf("stream error: %v", err)

		return fmt.Errorf("failed to start packet forwarding: %w", err)
	}

	// Convert stream to net.Conn
	conn := grpc.StreamToConn(stream)

	// Create tunnel
	t, err := tunnel.NewTunnel(conn)
	if err != nil {
		status.Error = fmt.Sprintf("tunnel error: %v", err)

		return fmt.Errorf("failed to create tunnel: %w", err)
	}
	s.tunnels[nodeID] = t

	// Create bridge
	b := bridge.NewWiresharkBridge(fmt.Sprintf("/tmp/pcap_%s.pipe", nodeID))
	s.bridges[nodeID] = b

	// Update status
	status.Interface = iface
	status.IsActive = true
	status.Error = ""

	// Start metrics collection
	if s.metrics != nil {
		err := s.metrics.AddMetric(nodeID, time.Now(), 0, "capture_start")
		if err != nil {
			return err
		}
	}

	// Start bridge
	go func() {
		if err := b.Start(ctx, t); err != nil {
			log.Printf("Bridge error for node %s: %v", nodeID, err)
			s.mu.Lock()
			status.Error = fmt.Sprintf("bridge error: %v", err)
			status.IsActive = false
			s.mu.Unlock()
		}
	}()

	// Start capture
	go func() {
		if err := t.StartPacketCapture(ctx, iface); err != nil {
			log.Printf("Capture error for node %s: %v", nodeID, err)
			s.mu.Lock()
			status.Error = fmt.Sprintf("capture error: %v", err)
			status.IsActive = false
			s.mu.Unlock()
		}
	}()

	return nil
}

// StopCapture stops a remote packet capture session
func (s *CaptureService) StopCapture(nodeID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	status := s.status[nodeID]
	if status != nil {
		status.IsActive = false
	}

	// Close tunnel
	if t, exists := s.tunnels[nodeID]; exists {
		if err := t.Close(); err != nil {
			log.Printf("Error closing tunnel for node %s: %v", nodeID, err)
		}
		delete(s.tunnels, nodeID)
	}

	// Close bridge
	if b, exists := s.bridges[nodeID]; exists {
		if err := b.Close(); err != nil {
			log.Printf("Error closing bridge for node %s: %v", nodeID, err)
		}
		delete(s.bridges, nodeID)
	}

	// Record metric
	if s.metrics != nil {
		err := s.metrics.AddMetric(nodeID, time.Now(), 0, "capture_stop")
		if err != nil {
			return err
		}
	}

	return nil
}

// GetCaptureStats returns statistics for an active capture session
func (s *CaptureService) GetCaptureStats(nodeID string) (*proto.CaptureStats, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	b, exists := s.bridges[nodeID]
	if !exists {
		return nil, fmt.Errorf("no active capture for node %s", nodeID)
	}

	stats := b.GetStats()

	// Update status
	if status := s.status[nodeID]; status != nil {
		status.PacketsRead = stats.PacketsReceived
		status.BytesRead = stats.BytesReceived
	}

	// Record metrics
	if s.metrics != nil {
		err := s.metrics.AddMetric(nodeID, time.Now(), int64(stats.PacketsReceived), "capture_packets")
		if err != nil {
			return nil, err
		}
		err = s.metrics.AddMetric(nodeID, time.Now(), int64(stats.BytesReceived), "capture_bytes")
		if err != nil {
			return nil, err
		}
	}

	return stats, nil
}
