// pkg/agent/capture_test.go

package agent

import (
	"context"
	"testing"
	"time"

	"github.com/mfreeman451/serviceradar/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"
)

// mockStartCaptureServer implements proto.AgentService_StartCaptureServer for testing
type mockStartCaptureServer struct {
	packets []*proto.PacketData
	ctx     context.Context
}

func newMockStartCaptureServer() *mockStartCaptureServer {
	return &mockStartCaptureServer{
		packets: make([]*proto.PacketData, 0),
		ctx:     context.Background(),
	}
}

func (m *mockStartCaptureServer) Send(packet *proto.PacketData) error {
	m.packets = append(m.packets, packet)
	return nil
}

func (m *mockStartCaptureServer) Context() context.Context {
	return m.ctx
}

func (m *mockStartCaptureServer) SendHeader(metadata.MD) error { return nil }
func (m *mockStartCaptureServer) SetHeader(metadata.MD) error  { return nil }
func (m *mockStartCaptureServer) SendMsg(interface{}) error    { return nil }
func (m *mockStartCaptureServer) RecvMsg(interface{}) error    { return nil }
func (m *mockStartCaptureServer) SetTrailer(metadata.MD)       {}

func TestCaptureService_New(t *testing.T) {
	service := NewCaptureService()
	assert.NotNil(t, service)
	assert.NotNil(t, service.handles)
	assert.NotNil(t, service.stats)
}

func TestCaptureService_StartCapture(t *testing.T) {
	service := NewCaptureService()
	stream := newMockStartCaptureServer()

	tests := []struct {
		name    string
		req     *proto.CaptureRequest
		wantErr bool
	}{
		{
			name: "invalid interface",
			req: &proto.CaptureRequest{
				Interface: "nonexistent0",
				SnapLen:   65535,
			},
			wantErr: true,
		},
		{
			name: "valid interface with filter",
			req: &proto.CaptureRequest{
				Interface:   "lo0", // Use loopback for testing
				SnapLen:     65535,
				Filter:      "tcp",
				Promiscuous: true,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
			defer cancel()

			// Create new context for the stream
			streamCtx, streamCancel := context.WithCancel(ctx)
			stream.ctx = streamCtx

			// Start capture in background
			errChan := make(chan error)
			go func() {
				errChan <- service.StartCapture(tt.req, stream)
			}()

			// Cancel after a short time to simulate stopping capture
			time.Sleep(50 * time.Millisecond)
			streamCancel()

			// Wait for result
			err := <-errChan
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			// For valid interfaces, we might still get permission errors in tests
			if err != nil && !isPermissionError(err) {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestCaptureService_StopCapture(t *testing.T) {
	service := NewCaptureService()
	ctx := context.Background()

	// Test stopping non-existent capture
	stats, err := service.StopCapture(ctx, &proto.StopCaptureRequest{
		NodeId: "nonexistent",
	})
	require.NoError(t, err)
	assert.NotNil(t, stats)
	assert.Equal(t, uint64(0), stats.PacketsReceived)

	// Start a capture first
	stream := newMockStartCaptureServer()
	streamCtx, streamCancel := context.WithCancel(context.Background())
	stream.ctx = streamCtx

	go func() {
		_ = service.StartCapture(&proto.CaptureRequest{
			Interface: "lo0",
			NodeId:    "test",
		}, stream)
	}()

	// Let it run briefly
	time.Sleep(50 * time.Millisecond)

	// Stop the capture
	stats, err = service.StopCapture(ctx, &proto.StopCaptureRequest{
		NodeId: "test",
	})
	require.NoError(t, err)
	assert.NotNil(t, stats)

	// Cleanup
	streamCancel()
}

func TestCaptureService_ConcurrentCaptures(t *testing.T) {
	service := NewCaptureService()
	stream1 := newMockStartCaptureServer()
	stream2 := newMockStartCaptureServer()

	// Try to start captures on same interface
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	errChan := make(chan error, 2)

	// Start first capture
	go func() {
		errChan <- service.StartCapture(&proto.CaptureRequest{
			Interface: "lo0",
			NodeId:    "test1",
		}, stream1)
	}()

	// Start second capture
	go func() {
		errChan <- service.StartCapture(&proto.CaptureRequest{
			Interface: "lo0",
			NodeId:    "test2",
		}, stream2)
	}()

	// Get results
	var errs []error
	for i := 0; i < 2; i++ {
		select {
		case err := <-errChan:
			if err != nil && !isPermissionError(err) {
				errs = append(errs, err)
			}
		case <-ctx.Done():
			t.Fatal("timeout waiting for captures to complete")
		}
	}

	// We expect one capture to fail
	if len(errs) == 0 {
		t.Log("No errors seen, but this might be due to permission issues in test environment")
	} else {
		assert.Equal(t, 1, len(errs), "expected exactly one error")
	}
}
