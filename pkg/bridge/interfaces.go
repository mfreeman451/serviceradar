// Package bridge pkg/bridge/bridge.go
package bridge

import (
	"context"

	"github.com/mfreeman451/serviceradar/pkg/tunnel"
	"github.com/mfreeman451/serviceradar/proto"
)

// Bridge defines the interface for packet bridging.
type Bridge interface {
	// Start begins forwarding packets from the tunnel to the bridge.
	Start(ctx context.Context, t tunnel.Tunnel) error

	// Close cleans up resources and stops packet forwarding.
	Close() error

	// GetStats returns current capture statistics.
	GetStats() *proto.CaptureStats
}
