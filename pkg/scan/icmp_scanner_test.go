package scan

import (
	"context"
	"testing"
	"time"

	"github.com/mfreeman451/serviceradar/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestICMPChecksum(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		expected uint16
	}{
		{
			name:     "Empty data",
			data:     []byte{},
			expected: 0xFFFF,
		},
		{
			name:     "Simple ICMP header",
			data:     []byte{8, 0, 0, 0, 0, 1, 0, 1},
			expected: 0xF7FD,
		},
		{
			name:     "Odd length data",
			data:     []byte{8, 0, 0, 0, 0, 1, 0, 1, 0},
			expected: 0xF7FD, // Corrected expected value
		},
	}

	scanner := &ICMPScanner{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := scanner.calculateChecksum(tt.data)
			if result != tt.expected {
				t.Errorf("calculateChecksum() = %#x, want %#x", result, tt.expected)
			}
		})
	}
}

func TestNewICMPScanner_Error(t *testing.T) {
	// Simulate an error by passing invalid parameters
	_, err := NewICMPScanner(0, 0, 0) // All parameters are invalid
	require.Error(t, err, "Expected error for invalid parameters")
}

func TestICMPScanner_SocketError(t *testing.T) {
	scanner, err := NewICMPScanner(1*time.Second, 1, 3)
	require.NoError(t, err, "Expected error for invalid socket")

	scanner.rawSocket = -1 // Invalid socket

	targets := []models.Target{
		{Host: "127.0.0.1", Mode: models.ModeICMP},
	}

	_, err = scanner.Scan(context.Background(), targets)
	require.Error(t, err, "Expected error for invalid socket")
}

func TestICMPScanner_Scan_InvalidTargets(t *testing.T) {
	scanner, err := NewICMPScanner(1*time.Second, 1, 3)
	require.NoError(t, err)

	targets := []models.Target{
		{Host: "invalid.host", Mode: models.ModeICMP},
	}

	results, err := scanner.Scan(context.Background(), targets)
	require.NoError(t, err)

	// Count results channel to ensure proper behavior
	resultCount := 0
	for range results {
		resultCount++
	}

	// We expect one result for the invalid target, with Available=false
	assert.Equal(t, 1, resultCount, "Expected one result for invalid target")

	// Clean up
	err = scanner.Stop()
	require.NoError(t, err)
}
