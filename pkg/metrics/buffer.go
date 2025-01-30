package metrics

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/mfreeman451/serviceradar/pkg/models"
)

// metricPoint represents a single metric data point.
type metricPoint struct {
	timestamp    int64
	responseTime int64
	serviceName  string
}

// LockFreeRingBuffer is a lock-free ring buffer implementation.
type LockFreeRingBuffer struct {
	points []atomic.Pointer[metricPoint]
	pos    atomic.Int64
	size   int64
	pool   sync.Pool
}

// NewBuffer creates a new MetricStore (e.g., RingBuffer or LockFreeRingBuffer).
func NewBuffer(size int) MetricStore {
	return NewLockFreeBuffer(size) // Use LockFreeRingBuffer by default
}

// NewLockFreeBuffer creates a new LockFreeRingBuffer with the specified size.
func NewLockFreeBuffer(size int) MetricStore {
	rb := &LockFreeRingBuffer{
		points: make([]atomic.Pointer[metricPoint], size),
		size:   int64(size),
		pool: sync.Pool{
			New: func() interface{} {
				return &models.MetricPoint{}
			},
		},
	}

	// Initialize atomic pointers
	for i := range rb.points {
		rb.points[i].Store(new(metricPoint))
	}

	return rb
}

// Add adds a new metric point to the buffer.
func (b *LockFreeRingBuffer) Add(timestamp time.Time, responseTime int64, serviceName string) {
	// Create new point
	newPoint := &metricPoint{
		timestamp:    timestamp.UnixNano(),
		responseTime: responseTime,
		serviceName:  serviceName,
	}

	// Atomically increment the position and get the index
	pos := b.pos.Add(1) - 1
	idx := pos % b.size

	// Atomically store the new point
	b.points[idx].Store(newPoint)
}

// GetPoints retrieves all metric points from the buffer.
func (b *LockFreeRingBuffer) GetPoints() []models.MetricPoint {
	// Load the current position atomically
	pos := b.pos.Load()
	points := make([]models.MetricPoint, b.size)

	for i := int64(0); i < b.size; i++ {
		// Calculate the index for the current point
		idx := (pos - i - 1 + b.size) % b.size

		// Atomically load the point
		p := b.points[idx].Load()
		if p == nil {
			continue
		}

		// Get a MetricPoint from the pool
		mp := b.pool.Get().(*models.MetricPoint)
		mp.Timestamp = time.Unix(0, p.timestamp)
		mp.ResponseTime = p.responseTime
		mp.ServiceName = p.serviceName

		points[i] = *mp

		// Return the MetricPoint to the pool
		b.pool.Put(mp)
	}

	return points
}
