package tunnel

import (
	"context"
	"io"
)

// Tunnel defines the interface for packet tunneling
type Tunnel interface {
	// GetPacketStream returns a stream for receiving packet data
	GetPacketStream(ctx context.Context) (io.ReadCloser, error)
	// StartPacketCapture starts capturing packets on the specified interface
	StartPacketCapture(ctx context.Context, iface string) error
	// Close closes the tunnel and all its streams
	Close() error
}
