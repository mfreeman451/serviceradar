// pkg/cloud/capture/service.go

package capture

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"sync"

	"github.com/hashicorp/yamux"
	"github.com/mfreeman451/serviceradar/pkg/bridge"
	"github.com/mfreeman451/serviceradar/pkg/tunnel"
)

type CaptureService struct {
	mu      sync.RWMutex
	tunnels map[string]*tunnel.Tunnel          // Map of nodeID to tunnel
	bridges map[string]*bridge.WiresharkBridge // Map of nodeID to Wireshark bridge
}

func NewCaptureService() *CaptureService {
	return &CaptureService{
		tunnels: make(map[string]*tunnel.Tunnel),
		bridges: make(map[string]*bridge.WiresharkBridge),
	}
}

// StartCapture starts a remote packet capture session
func (s *CaptureService) StartCapture(ctx context.Context, nodeID, iface string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check if we already have a tunnel for this node
	t, exists := s.tunnels[nodeID]
	if !exists {
		// Create new tunnel if needed
		conn, err := s.dialNode(ctx, nodeID)
		if err != nil {
			return fmt.Errorf("failed to connect to node: %w", err)
		}

		t, err = tunnel.NewTunnel(conn)
		if err != nil {
			return fmt.Errorf("failed to create tunnel: %w", err)
		}
		s.tunnels[nodeID] = t
	}

	// Create Wireshark bridge
	bridge := bridge.NewWiresharkBridge(fmt.Sprintf("/tmp/pcap_%s.pipe", nodeID))
	s.bridges[nodeID] = bridge

	// Start bridge in background
	go func() {
		if err := bridge.Start(ctx, t); err != nil {
			log.Printf("Bridge error for node %s: %v", nodeID, err)
		}
	}()

	// Start packet capture on the remote node
	go func() {
		if err := t.StartPacketCapture(ctx, iface); err != nil {
			log.Printf("Capture error for node %s: %v", nodeID, err)
		}
	}()

	return nil
}

// StopCapture stops a remote packet capture session
func (s *CaptureService) StopCapture(nodeID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

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

	return nil
}

// dialNode establishes connection to a remote node
func (s *CaptureService) dialNode(ctx context.Context, nodeID string) (net.Conn, error) {
	// This would integrate with your existing agent connection mechanism
	// For now, just a placeholder
	return nil, fmt.Errorf("not implemented")
}

// GetCaptureStats returns statistics for an active capture session
func (s *CaptureService) GetCaptureStats(nodeID string) (*CaptureStats, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	bridge, exists := s.bridges[nodeID]
	if !exists {
		return nil, fmt.Errorf("no active capture for node %s", nodeID)
	}

	return bridge.GetStats(), nil
}
