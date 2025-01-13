package agent

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
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
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to external checker: %w", err)
	}

	return &ExternalChecker{
		name:       name,
		address:    address,
		client:     grpc_health_v1.NewHealthClient(conn),
		connection: conn,
	}, nil
}

// Check implements the checker.Checker interface
func (e *ExternalChecker) Check(ctx context.Context) (bool, string) {
	// Add timeout to context
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := e.client.Check(ctx, &grpc_health_v1.HealthCheckRequest{})
	if err != nil {
		return false, fmt.Sprintf("Failed to check %s: %v", e.name, err)
	}

	healthy := resp.Status == grpc_health_v1.HealthCheckResponse_SERVING
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
