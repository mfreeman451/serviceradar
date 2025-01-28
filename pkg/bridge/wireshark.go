// pkg/bridge/wireshark.go

package bridge

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
	"syscall"
	"time"

	"github.com/hashicorp/yamux"
	"github.com/mfreeman451/serviceradar/pkg/tunnel"
)

type WiresharkBridge struct {
	pipePath string
	session  *yamux.Session
}

func NewWiresharkBridge(pipePath string) *WiresharkBridge {
	return &WiresharkBridge{
		pipePath: pipePath,
	}
}

func (b *WiresharkBridge) Start(ctx context.Context, t tunnel.Tunnel) error {
	// setup a context with a timeout
	_, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// Create named pipe for Wireshark
	if err := syscall.Mkfifo(b.pipePath, 0666); err != nil {
		return fmt.Errorf("failed to create pipe: %w", err)
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

	// Read packets from YAMUX session and write to pipe
	stream, err := b.session.AcceptStream()
	if err != nil {
		return fmt.Errorf("failed to accept stream: %w", err)
	}

	buffer := make([]byte, 65536)
	for {
		n, err := stream.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("Error reading from stream: %v", err)
			continue
		}

		if _, err := pipe.Write(buffer[:n]); err != nil {
			return fmt.Errorf("failed to write to pipe: %w", err)
		}
	}

	return nil
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
