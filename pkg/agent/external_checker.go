package agent

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/metadata"
)

// ExternalChecker implements checker.Checker for external checker processes
type ExternalChecker struct {
	name       string
	address    string
	client     grpc_health_v1.HealthClient
	connection *grpc.ClientConn
}

// NewExternalChecker creates a new checker that connects to an external process
func NewExternalChecker(name, address string) (*ExternalChecker, error) {
	log.Printf("Creating new external checker %s at %s", name, address)

	// Create gRPC connection with proper options
	conn, err := grpc.Dial(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
		grpc.WithTimeout(5*time.Second),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to external checker: %w", err)
	}

	checker := &ExternalChecker{
		name:       name,
		address:    address,
		client:     grpc_health_v1.NewHealthClient(conn),
		connection: conn,
	}

	// Test the connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = checker.client.Check(ctx, &grpc_health_v1.HealthCheckRequest{})
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("initial health check failed: %w", err)
	}

	log.Printf("Successfully created external checker %s", name)
	return checker, nil
}

func (e *ExternalChecker) Check(ctx context.Context) (bool, string) {
	// Create a new context with metadata
	ctx = metadata.NewOutgoingContext(ctx, metadata.New(nil))

	var header metadata.MD
	resp, err := e.client.Check(ctx, &grpc_health_v1.HealthCheckRequest{},
		grpc.Header(&header)) // Add this to capture headers
	if err != nil {
		return false, fmt.Sprintf("Failed to check %s: %v", e.name, err)
	}

	healthy := resp.Status == grpc_health_v1.HealthCheckResponse_SERVING

	// Extract block details from headers
	if details := header.Get("block-details"); len(details) > 0 {
		log.Printf("Received block details: %s", details[0])
		return healthy, details[0]
	}

	if !healthy {
		return false, fmt.Sprintf("%s is not healthy", e.name)
	}

	return true, fmt.Sprintf("%s is healthy", e.name)
}

// Close cleans up the checker's resources
func (e *ExternalChecker) Close() error {
	if e.connection != nil {
		return e.connection.Close()
	}
	return nil
}
