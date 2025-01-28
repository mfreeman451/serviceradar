// pkg/bridge/wireshark.go

package bridge

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/mfreeman451/serviceradar/pkg/tunnel"
	"github.com/mfreeman451/serviceradar/proto"
)

type WiresharkBridge struct {
	pipePath string
	stats    *proto.CaptureStats
	done     chan struct{}
}

func NewWiresharkBridge(pipePath string) Bridge {
	return &WiresharkBridge{
		pipePath: pipePath,
		stats: &proto.CaptureStats{
			StartTime: time.Now().Unix(),
		},
		done: make(chan struct{}),
	}
}

func (b *WiresharkBridge) Start(ctx context.Context, t tunnel.Tunnel) error {
	// Create named pipe for Wireshark
	if err := syscall.Mkfifo(b.pipePath, 0666); err != nil {
		if !os.IsExist(err) {
			return fmt.Errorf("failed to create pipe: %w", err)
		}
	}

	// Open pipe for writing
	pipe, err := os.OpenFile(b.pipePath, os.O_WRONLY, os.ModeNamedPipe)
	if err != nil {
		return fmt.Errorf("failed to open pipe: %w", err)
	}
	defer pipe.Close()

	// Write pcap file header
	if err := b.writePcapHeader(pipe); err != nil {
		return fmt.Errorf("failed to write pcap header: %w", err)
	}

	// Get packet stream from tunnel
	stream, err := t.GetPacketStream(ctx)
	if err != nil {
		return fmt.Errorf("failed to get packet stream: %w", err)
	}
	defer stream.Close()

	buffer := make([]byte, 65536)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-b.done:
			return nil
		default:
			n, err := stream.Read(buffer)
			if err == io.EOF {
				return nil
			}
			if err != nil {
				log.Printf("Error reading from stream: %v", err)
				continue
			}

			if _, err := pipe.Write(buffer[:n]); err != nil {
				return fmt.Errorf("failed to write to pipe: %w", err)
			}

			atomic.AddUint64(&b.stats.PacketsReceived, 1)
			atomic.AddUint64(&b.stats.BytesReceived, uint64(n))
		}
	}
}

func (b *WiresharkBridge) Close() error {
	close(b.done)
	b.stats.EndTime = time.Now().Unix()
	return os.Remove(b.pipePath)
}

func (b *WiresharkBridge) GetStats() *proto.CaptureStats {
	stats := &proto.CaptureStats{
		PacketsReceived: atomic.LoadUint64(&b.stats.PacketsReceived),
		BytesReceived:   atomic.LoadUint64(&b.stats.BytesReceived),
		PacketsDropped:  atomic.LoadUint64(&b.stats.PacketsDropped),
		StartTime:       b.stats.StartTime,
	}
	if b.stats.EndTime > 0 {
		stats.EndTime = b.stats.EndTime
	} else {
		stats.EndTime = time.Now().Unix()
	}
	return stats
}

func (b *WiresharkBridge) writePcapHeader(w io.Writer) error {
	header := &pcapHeader{
		MagicNumber:  0xa1b2c3d4,
		VersionMajor: 2,
		VersionMinor: 4,
		ThisZone:     0,
		SigFigs:      0,
		SnapLen:      65535,
		Network:      1, // LINKTYPE_ETHERNET
	}

	return binary.Write(w, binary.LittleEndian, header)
}

type pcapHeader struct {
	MagicNumber  uint32
	VersionMajor uint16
	VersionMinor uint16
	ThisZone     int32
	SigFigs      uint32
	SnapLen      uint32
	Network      uint32
}
