package tunnel

import (
	"context"
	"io"
	"net"
)

// Tunnel defines the interface for packet tunneling.
type Tunnel interface {
	// OpenStream opens a new stream with the given ID
	OpenStream(id uint32) (net.Conn, error)

	// GetPacketStream returns a stream for receiving packet data
	GetPacketStream(ctx context.Context) (io.ReadCloser, error)

	// StartPacketCapture starts capturing packets on the specified interface
	StartPacketCapture(ctx context.Context, iface string) error

	// Close closes the tunnel and all its streams
	Close() error
}

// YamuxStream represents a YAMUX stream that implements net.Conn.
type YamuxStream interface {
	net.Conn
	Write(p []byte) (n int, err error)
}
