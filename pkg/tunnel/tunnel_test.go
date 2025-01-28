// pkg/tunnel/tunnel_test.go

package tunnel

import (
	"bytes"
	"context"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockConn implements net.Conn for testing.
type mockConn struct {
	readData  *bytes.Buffer
	writeData *bytes.Buffer
	closed    bool
	mu        sync.Mutex
}

func newMockConn() *mockConn {
	return &mockConn{
		readData:  new(bytes.Buffer),
		writeData: new(bytes.Buffer),
	}
}

func (m *mockConn) Read(p []byte) (n int, err error) {
	return m.readData.Read(p)
}

func (m *mockConn) Write(p []byte) (n int, err error) {
	return m.writeData.Write(p)
}

func (m *mockConn) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.closed = true
	return nil
}

func (*mockConn) LocalAddr() net.Addr                { return nil }
func (*mockConn) RemoteAddr() net.Addr               { return nil }
func (m *mockConn) SetDeadline(time.Time) error      { return nil }
func (m *mockConn) SetReadDeadline(time.Time) error  { return nil }
func (m *mockConn) SetWriteDeadline(time.Time) error { return nil }

func TestNewTunnel(t *testing.T) {
	tests := []struct {
		name    string
		conn    net.Conn
		wantErr bool
	}{
		{
			name:    "valid connection",
			conn:    newMockConn(),
			wantErr: false,
		},
		{
			name:    "nil connection",
			conn:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tunnel, err := NewTunnel(tt.conn)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, tunnel)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, tunnel)

			// Test cleanup
			assert.NoError(t, tunnel.Close())
		})
	}
}

func TestTunnel_OpenStream(t *testing.T) {
	mockConn := newMockConn()
	tunnel, err := NewTunnel(mockConn)
	require.NoError(t, err)
	defer tunnel.Close()

	tests := []struct {
		name    string
		id      uint32
		wantErr bool
	}{
		{
			name:    "valid stream ID",
			id:      1,
			wantErr: false,
		},
		{
			name:    "duplicate stream ID",
			id:      1,
			wantErr: false, // Should return existing stream
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stream, err := tunnel.OpenStream(tt.id)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, stream)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, stream)
		})
	}
}

func TestTunnel_StartPacketCapture(t *testing.T) {
	mockConn := newMockConn()
	tunnel, err := NewTunnel(mockConn)
	require.NoError(t, err)
	defer func(tunnel Tunnel) {
		err := tunnel.Close()
		if err != nil {
			t.Errorf("Failed to close tunnel: %v", err)
		}
	}(tunnel)

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	// Start capture in background
	errChan := make(chan error)
	go func() {
		errChan <- tunnel.StartPacketCapture(ctx, "lo0") // Use loopback for testing
	}()

	// Wait for either error or timeout
	select {
	case err := <-errChan:
		if err != nil && !isPermissionError(err) { // Ignore permission errors in tests
			t.Errorf("StartPacketCapture failed: %v", err)
		}
	case <-ctx.Done():
		// This is expected since we're using a short timeout
	}
}

func TestTunnel_GetPacketStream(t *testing.T) {
	mockConn := newMockConn()
	tunnel, err := NewTunnel(mockConn)
	require.NoError(t, err)
	defer func(tunnel Tunnel) {
		err := tunnel.Close()
		if err != nil {
			t.Errorf("Failed to close tunnel: %v", err)
		}
	}(tunnel)

	ctx := context.Background()
	stream, err := tunnel.GetPacketStream(ctx)
	require.NoError(t, err)
	require.NotNil(t, stream)

	// Test stream can be closed
	assert.NoError(t, stream.Close())
}

func TestTunnel_Close(t *testing.T) {
	mockConn := newMockConn()
	tunnel, err := NewTunnel(mockConn)
	require.NoError(t, err)

	// Open some streams
	stream1, err := tunnel.OpenStream(1)
	require.NoError(t, err)
	stream2, err := tunnel.OpenStream(2)
	require.NoError(t, err)

	// Close tunnel
	assert.NoError(t, tunnel.Close())

	// Verify all streams are closed
	_, err = stream1.Write([]byte("test"))
	assert.Error(t, err, "stream should be closed")
	_, err = stream2.Write([]byte("test"))
	assert.Error(t, err, "stream should be closed")

	// Verify underlying connection is closed
	assert.True(t, mockConn.closed, "underlying connection should be closed")
}

// Helper function to check if error is a permission error
func isPermissionError(err error) bool {
	if err == nil {
		return false
	}
	return err.Error() == "operation not permitted" ||
		err.Error() == "permission denied"
}
