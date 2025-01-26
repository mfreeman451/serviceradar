package metrics

import (
	"bytes"
	"encoding/binary"
	"log"
	"math"
	"sync"
	"time"

	"github.com/mfreeman451/serviceradar/pkg/models"
)

const (
	metricPointSize = 32
)

type NodeMetrics struct {
	buffer MetricBuffer
	mu     sync.RWMutex // Fine-grained lock for this node
}

type MetricBuffer struct {
	buffer []byte
	pos    int
	size   int
	mu     sync.RWMutex
}

func NewMetricBuffer(points int) *MetricBuffer {
	return &MetricBuffer{
		buffer: make([]byte, points*metricPointSize),
		size:   points,
	}
}

func (b *MetricBuffer) Add(timestamp time.Time, responseTime int64, serviceName string) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	offset := b.pos * metricPointSize

	// Safely convert timestamp.UnixNano() to uint64
	unixNano := timestamp.UnixNano()
	if unixNano < 0 {
		log.Printf("Warning: negative timestamp for service '%s': %d (clamped to 0)", serviceName, unixNano)
		unixNano = 0
	}

	binary.LittleEndian.PutUint64(b.buffer[offset:], uint64(unixNano)) //nolint:gosec // Safe: negative values are clamped to 0

	// Safely convert responseTime to uint64
	if responseTime < 0 {
		log.Printf("Warning: negative response time for service '%s': %d (clamped to 0)", serviceName, responseTime)
		responseTime = 0
	}

	binary.LittleEndian.PutUint64(b.buffer[offset+8:], uint64(responseTime)) //nolint:gosec // Safe: negative values are clamped to 0

	// Write service name (ensure it fits within 16 bytes)
	copy(b.buffer[offset+16:offset+32], serviceName)

	b.pos = (b.pos + 1) % b.size

	return nil
}

func (b *MetricBuffer) GetPoints() []models.MetricPoint {
	b.mu.RLock()
	defer b.mu.RUnlock()

	points := make([]models.MetricPoint, b.size)

	for i := 0; i < b.size; i++ {
		pos := (b.pos - i - 1 + b.size) % b.size
		offset := pos * metricPointSize

		ts := binary.LittleEndian.Uint64(b.buffer[offset:])
		rt := binary.LittleEndian.Uint64(b.buffer[offset+8:])
		sn := string(bytes.TrimRight(b.buffer[offset+16:offset+32], "\x00"))

		// Safely handle uint64 to int64 conversion
		timestamp := safeUint64ToInt64(ts, sn, "timestamp")
		responseTime := safeUint64ToInt64(rt, sn, "response time")

		points[i] = models.MetricPoint{
			Timestamp:    time.Unix(0, timestamp),
			ResponseTime: responseTime,
			ServiceName:  sn,
		}
	}

	return points
}

// safeUint64ToInt64 safely converts a uint64 to int64.
func safeUint64ToInt64(value uint64, serviceName, fieldName string) int64 {
	if value > math.MaxInt64 {
		log.Printf("Warning: %s overflow for service '%s': %d (clamped to MaxInt64)", fieldName, serviceName, value)
		return math.MaxInt64
	}

	return int64(value)
}
