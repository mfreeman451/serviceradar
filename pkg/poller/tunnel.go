// pkg/poller/tunnel.go

package poller

import (
	"context"
	"encoding/binary"
	"fmt"
	"log"
	"sync"

	"github.com/hashicorp/yamux"
	"github.com/mfreeman451/serviceradar/pkg/grpc"
	"github.com/mfreeman451/serviceradar/proto"
)

type TunnelManager struct {
	agentClient *grpc.ClientConn
	cloudClient *grpc.ClientConn
	sessions    map[string]*yamux.Session
	mu          sync.RWMutex
}

func (tm *TunnelManager) StartTunnel(ctx context.Context, nodeID string) error {
	// 1. Start agent capture stream
	agentStream, err := proto.NewAgentServiceClient(
		tm.agentClient.GetConnection(),
	).StartCapture(ctx, &proto.CaptureRequest{
		NodeId: nodeID,
	})
	if err != nil {
		return fmt.Errorf("failed to start agent capture: %w", err)
	}

	// 2. Start cloud forwarding stream
	cloudStream, err := proto.NewPollerServiceClient(
		tm.cloudClient.GetConnection(),
	).ForwardPackets(ctx)
	if err != nil {
		return fmt.Errorf("failed to start packet forwarding: %w", err)
	}

	// 3. Set up YAMUX session
	config := yamux.DefaultConfig()
	session, err := yamux.Client(grpc.StreamToConn(cloudStream), config)
	if err != nil {
		return fmt.Errorf("failed to create yamux session: %w", err)
	}

	tm.mu.Lock()
	tm.sessions[nodeID] = session
	tm.mu.Unlock()

	// 4. Start packet forwarding
	go func() {
		defer func() {
			tm.mu.Lock()
			delete(tm.sessions, nodeID)
			session.Close()
			tm.mu.Unlock()
		}()

		for {
			select {
			case <-ctx.Done():
				return
			default:
				packet, err := agentStream.Recv()
				if err != nil {
					log.Printf("Error receiving from agent: %v", err)
					return
				}

				// Forward to cloud
				if err := cloudStream.Send(packet); err != nil {
					log.Printf("Error forwarding to cloud: %v", err)
					return
				}

				// Get response
				resp, err := cloudStream.Recv()
				if err != nil {
					log.Printf("Error receiving cloud response: %v", err)
					return
				}

				if resp.Error != "" {
					log.Printf("Cloud reported error: %s", resp.Error)
				}
			}
		}
	}()

	return nil
}

func (tm *TunnelManager) writePacket(stream *yamux.Stream, packet *proto.PacketData) error {
	// Write packet size first
	size := uint32(len(packet.Data))
	if err := binary.Write(stream, binary.BigEndian, size); err != nil {
		return fmt.Errorf("failed to write size: %w", err)
	}

	// Write packet data
	if _, err := stream.Write(packet.Data); err != nil {
		return fmt.Errorf("failed to write packet: %w", err)
	}

	return nil
}
