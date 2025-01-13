// pkg/poller/poller.go
package poller

import (
	"context"
	"log"
	"time"

	pb "github.com/mfreeman451/homemon/proto"
	"google.golang.org/grpc"
)

type Config struct {
	Agents       map[string]string
	CloudAddress string
	PollInterval time.Duration
	PollerID     string
}

type Poller struct {
	config      Config
	cloudClient pb.PollerServiceClient
}

func New(config Config) (*Poller, error) {
	conn, err := grpc.Dial(config.CloudAddress, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return &Poller{
		config:      config,
		cloudClient: pb.NewPollerServiceClient(conn),
	}, nil
}

func (p *Poller) pollAgent(ctx context.Context, agentAddr string) (*pb.StatusResponse, error) {
	conn, err := grpc.Dial(agentAddr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := pb.NewAgentServiceClient(conn)
	return client.GetStatus(ctx, &pb.StatusRequest{ServiceName: "nginx"})
}

func (p *Poller) Start(ctx context.Context) error {
	ticker := time.NewTicker(p.config.PollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if err := p.poll(ctx); err != nil {
				log.Printf("Error during poll: %v", err)
			}
		}
	}
}

func (p *Poller) poll(ctx context.Context) error {
	statuses := make([]*pb.ServiceStatus, 0)

	for name, addr := range p.config.Agents {
		status, err := p.pollAgent(ctx, addr)
		if err != nil {
			log.Printf("Error polling agent %s: %v", name, err)
			statuses = append(statuses, &pb.ServiceStatus{
				ServiceName: name,
				Available:   false,
				Message:     err.Error(),
			})
			continue
		}

		statuses = append(statuses, &pb.ServiceStatus{
			ServiceName: name,
			Available:   status.Available,
			Message:     status.Message,
		})
	}

	_, err := p.cloudClient.ReportStatus(ctx, &pb.PollerStatusRequest{
		Services:  statuses,
		PollerId:  p.config.PollerID,
		Timestamp: time.Now().Unix(),
	})
	return err
}
