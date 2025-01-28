// Package agent pkg/agent/capture.go

package agent

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"github.com/mfreeman451/serviceradar/proto"
)

type CaptureService struct {
	mu      sync.RWMutex
	handles map[string]*pcap.Handle
	stats   map[string]*proto.CaptureStats
}

func NewCaptureService() *CaptureService {
	return &CaptureService{
		handles: make(map[string]*pcap.Handle),
		stats:   make(map[string]*proto.CaptureStats),
	}
}

func (s *CaptureService) StartCapture(req *proto.CaptureRequest, stream proto.AgentService_StartCaptureServer) error {
	s.mu.Lock()
	// Check if we're already capturing on this interface
	if _, exists := s.handles[req.Interface]; exists {
		s.mu.Unlock()

		return fmt.Errorf("already capturing on interface %s", req.Interface)
	}

	// Initialize stats
	s.stats[req.Interface] = &proto.CaptureStats{}
	s.mu.Unlock()

	// Open the capture device
	handle, err := pcap.OpenLive(req.Interface, int32(req.SnapLen), req.Promiscuous, pcap.BlockForever)
	if err != nil {
		return fmt.Errorf("failed to open interface: %w", err)
	}

	// Apply BPF filter if provided
	if req.Filter != "" {
		if err := handle.SetBPFFilter(req.Filter); err != nil {
			handle.Close()
			return fmt.Errorf("failed to set BPF filter: %w", err)
		}
	}

	s.mu.Lock()
	s.handles[req.Interface] = handle
	s.mu.Unlock()

	defer func() {
		s.mu.Lock()
		delete(s.handles, req.Interface)
		s.mu.Unlock()
		handle.Close()
	}()

	// Start capturing packets
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	var sequence uint32

	log.Printf("Starting packet capture on interface %s", req.Interface)

	for packet := range packetSource.Packets() {
		select {
		case <-stream.Context().Done():
			return nil
		default:
			// Send packet data through the gRPC stream
			if err := stream.Send(&proto.PacketData{
				Data:      packet.Data(),
				Timestamp: packet.Metadata().Timestamp.UnixNano(),
				Interface: req.Interface,
				Sequence:  sequence,
			}); err != nil {
				return fmt.Errorf("failed to send packet: %w", err)
			}

			// Update stats
			s.mu.Lock()
			if stats := s.stats[req.Interface]; stats != nil {
				stats.PacketsReceived++
				stats.BytesReceived += uint64(len(packet.Data()))
			}
			s.mu.Unlock()

			sequence++
		}
	}

	return nil
}

func (s *CaptureService) StopCapture(ctx context.Context, req *proto.StopCaptureRequest) (*proto.CaptureStats, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	// Find and close the handle
	for iface, handle := range s.handles {
		handle.Close()
		delete(s.handles, iface)
	}

	// Get the stats
	stats := s.stats[req.NodeId]
	if stats == nil {
		stats = &proto.CaptureStats{}
	}
	delete(s.stats, req.NodeId)

	return stats, nil
}
