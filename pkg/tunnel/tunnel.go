// pkg/tunnel/tunnel.go

package tunnel

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"github.com/hashicorp/yamux"
)

const (
	StreamControl = 1
	StreamPackets = 2
	StreamStats   = 3
)

// YamuxTunnel implements the Tunnel interface using YAMUX.
type YamuxTunnel struct {
	session *yamux.Session
	streams map[uint32]net.Conn
	mu      sync.RWMutex
	done    chan struct{}
}

// Stats provides statistics about the tunnel.
type Stats struct {
	PacketsForwarded uint64
	BytesForwarded   uint64
	Errors           uint64
}

// NewTunnel creates a new multiplexed tunnel.
func NewTunnel(conn net.Conn) (Tunnel, error) {
	// Configure YAMUX for high throughput
	config := yamux.DefaultConfig()
	config.EnableKeepAlive = true
	config.KeepAliveInterval = 30           // More aggressive keepalive
	config.MaxStreamWindowSize = 256 * 1024 // Larger windows for better throughput

	session, err := yamux.Client(conn, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create yamux session: %w", err)
	}

	return &YamuxTunnel{
		session: session,
		streams: make(map[uint32]net.Conn),
		done:    make(chan struct{}),
	}, nil
}

// OpenStream opens a new stream with the given ID
func (t *YamuxTunnel) OpenStream(id uint32) (net.Conn, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if stream, exists := t.streams[id]; exists {
		return stream, nil
	}

	stream, err := t.session.OpenStream()
	if err != nil {
		return nil, fmt.Errorf("failed to open stream %d: %w", id, err)
	}

	// Write stream ID as first byte
	if _, err := stream.Write([]byte{byte(id)}); err != nil {
		err := stream.Close()
		if err != nil {
			return nil, err
		}

		return nil, fmt.Errorf("failed to write stream ID: %w", err)
	}

	t.streams[id] = stream
	return stream, nil
}

// GetPacketStream returns a stream for receiving packet data
func (t *YamuxTunnel) GetPacketStream(ctx context.Context) (io.ReadCloser, error) {
	return t.OpenStream(StreamPackets)
}

// StartPacketCapture initiates packet capture on the specified interface
func (t *YamuxTunnel) StartPacketCapture(ctx context.Context, iface string) error {
	// Open packet stream
	packetStream, err := t.OpenStream(StreamPackets)
	if err != nil {
		return fmt.Errorf("failed to open packet stream: %w", err)
	}
	defer func(packetStream net.Conn) {
		err := packetStream.Close()
		if err != nil {
			log.Printf("failed to close packet stream: %v", err)
		}
	}(packetStream)

	// Open stats stream
	statsStream, err := t.OpenStream(StreamStats)
	if err != nil {
		return fmt.Errorf("failed to open stats stream: %w", err)
	}
	defer func(statsStream net.Conn) {
		err := statsStream.Close()
		if err != nil {
			log.Printf("failed to close stats stream: %v", err)
		}
	}(statsStream)

	// Start packet capture
	handle, err := pcap.OpenLive(iface, 65535, true, pcap.BlockForever)
	if err != nil {
		return fmt.Errorf("failed to open interface %s: %w", iface, err)
	}
	defer handle.Close()

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

	var stats struct {
		PacketsSent uint64
		BytesSent   uint64
	}

	// Start stats reporting goroutine
	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				statsJSON, _ := json.Marshal(stats)
				if _, err := statsStream.Write(statsJSON); err != nil {
					log.Printf("failed to write stats: %v", err)
					return
				}
			}
		}
	}()

	// Main packet processing loop
	for packet := range packetSource.Packets() {
		select {
		case <-ctx.Done():
			return nil
		default:
			data := packet.Data()
			if _, err := packetStream.Write(data); err != nil {
				return fmt.Errorf("failed to write packet: %w", err)
			}

			stats.PacketsSent++
			stats.BytesSent += uint64(len(data))
		}
	}

	return nil
}

// Close closes the tunnel and all streams
func (t *YamuxTunnel) Close() error {
	close(t.done)

	t.mu.Lock()
	defer t.mu.Unlock()

	for _, stream := range t.streams {
		err := stream.Close()
		if err != nil {
			return err
		}
	}

	return t.session.Close()
}
