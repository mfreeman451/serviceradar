// Package scan provides network scanning functionality
package scan

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/mfreeman451/serviceradar/pkg/models"
)

const (
	defaultCleanupInterval  = 30 * time.Second
	defaultMaxLifetime      = 10 * time.Minute
	defaultIdleTimeout      = 1 * time.Minute
	defaultReadDeadline     = 100 * time.Millisecond
	defaultDNSLookupTimeout = 200 * time.Millisecond
)

// TCPScanner implementation using the improved connection pool.
type TCPScanner struct {
	timeout     time.Duration
	concurrency int
	pool        *connectionPool
	dialer      *net.Dialer
	done        chan struct{}
	closeOnce   sync.Once
}

// connEntry represents a connection in the pool.
type connEntry struct {
	conn      net.Conn
	createdAt time.Time
	lastUsed  time.Time
}

// connectionPool manages a pool of reusable TCP connections with proper lifecycle management.
type connectionPool struct {
	mu          sync.RWMutex
	connections map[string][]*connEntry
	maxIdle     int
	maxLifetime time.Duration
	idleTimeout time.Duration
	cleaner     *cleanupManager
	closeOnce   sync.Once
}

// newConnectionPool creates a new connection pool with proper lifecycle management.
func newConnectionPool(maxIdle int, maxLifetime, idleTimeout time.Duration) *connectionPool {
	p := &connectionPool{
		connections: make(map[string][]*connEntry),
		maxIdle:     maxIdle,
		maxLifetime: maxLifetime,
		idleTimeout: idleTimeout,
	}

	p.cleaner = newCleanupManager(defaultCleanupInterval, p.cleanup)
	p.cleaner.start()

	return p
}

// cleanup removes stale connections from the pool.
func (p *connectionPool) cleanup() {
	if !p.mu.TryLock() {
		// If we can't get the lock immediately, skip this cleanup cycle
		return
	}
	defer p.mu.Unlock()

	now := time.Now()

	for address, entries := range p.connections {
		valid := make([]*connEntry, 0, len(entries))

		for _, entry := range entries {
			// Check if the connection is too old or has been idle too long
			if now.Sub(entry.createdAt) > p.maxLifetime ||
				now.Sub(entry.lastUsed) > p.idleTimeout {
				if err := entry.conn.Close(); err != nil {
					log.Printf("Error closing stale connection: %v", err)
				}

				continue
			}

			valid = append(valid, entry)
		}

		if len(valid) == 0 {
			delete(p.connections, address)
		} else {
			p.connections[address] = valid
		}
	}
}

// get retrieves a connection from the pool or creates a new one.
func (p *connectionPool) get(ctx context.Context, dialer *net.Dialer, address string) (net.Conn, error) {
	// First try to get an existing connection
	if conn := p.getExisting(address); conn != nil {
		return conn, nil
	}

	// Create a new connection if none available
	conn, err := dialer.DialContext(ctx, "tcp", address)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection: %w", err)
	}

	return conn, nil
}

// getExisting retrieves an existing connection from the pool.
func (p *connectionPool) getExisting(address string) net.Conn {
	p.mu.Lock()
	defer p.mu.Unlock()

	entries := p.connections[address]
	now := time.Now()

	// Look for a valid connection
	for i := len(entries) - 1; i >= 0; i-- {
		entry := entries[i]

		// Skip connections that are too old or idle
		if now.Sub(entry.createdAt) > p.maxLifetime ||
			now.Sub(entry.lastUsed) > p.idleTimeout {
			continue
		}

		// Remove this connection from the pool and return it
		p.connections[address] = append(entries[:i], entries[i+1:]...)

		return entry.conn
	}

	return nil
}

// put returns a connection to the pool or closes it if the pool is full.
func (p *connectionPool) put(address string, conn net.Conn) {
	p.mu.Lock()
	defer p.mu.Unlock()

	// If pool is full for this address, close the connection
	if len(p.connections[address]) >= p.maxIdle {
		if err := conn.Close(); err != nil {
			log.Printf("Error closing connection when pool is full: %v", err)
		}

		return
	}

	// Add the connection to the pool
	entry := &connEntry{
		conn:      conn,
		createdAt: time.Now(),
		lastUsed:  time.Now(),
	}

	p.connections[address] = append(p.connections[address], entry)
}

// close closes all connections in the pool and stops the cleanup goroutine.
func (p *connectionPool) close() {
	p.closeOnce.Do(func() {
		if p.cleaner != nil {
			p.cleaner.stop()
		}

		p.mu.Lock()
		defer p.mu.Unlock()

		for _, entries := range p.connections {
			for _, entry := range entries {
				err := entry.conn.Close()
				if err != nil {
					return
				}
			}
		}

		p.connections = nil
	})
}

func NewTCPScanner(timeout time.Duration, concurrency, maxIdle int, maxLifetime, idleTimeout time.Duration) *TCPScanner {
	dialer := &net.Dialer{
		Timeout:   timeout,
		KeepAlive: 30 * time.Second,
		Resolver: &net.Resolver{
			PreferGo: true,
			Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
				d := net.Dialer{Timeout: defaultDNSLookupTimeout}
				return d.DialContext(ctx, network, address)
			},
		},
	}

	return &TCPScanner{
		timeout:     timeout,
		concurrency: concurrency,
		pool:        newConnectionPool(maxIdle, maxLifetime, idleTimeout),
		dialer:      dialer,
		done:        make(chan struct{}), // Initialize the done channel
	}
}

func (s *TCPScanner) scanTarget(ctx context.Context, target models.Target, results chan<- models.Result) {
	// Initialize result
	result := models.Result{
		Target:    target,
		FirstSeen: time.Now(),
		LastSeen:  time.Now(),
	}

	// Create a shorter timeout context for dialing
	connCtx, cancel := context.WithTimeout(ctx, s.timeout/2)
	defer cancel()

	address := fmt.Sprintf("%s:%d", target.Host, target.Port)
	conn, err := s.pool.get(connCtx, s.dialer, address)
	if err != nil {
		s.handleDialError(ctx, err, &result, results)
		return
	}

	// Flag to track successful connection for cleanup purposes.
	success := false
	defer s.cleanupConnection(ctx, address, conn, &success)

	// Check connection and record response time.
	startTime := time.Now()
	result.Available = s.checkConnection(conn)
	result.RespTime = time.Since(startTime)
	success = result.Available

	s.dispatchResult(ctx, results, &result, target)
}

func (s *TCPScanner) handleDialError(ctx context.Context, err error, result *models.Result, results chan<- models.Result) {
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return
	}
	result.Error = err
	result.Available = false
	select {
	case results <- *result:
	case <-ctx.Done():
	}
}

func (s *TCPScanner) cleanupConnection(ctx context.Context, address string, conn net.Conn, success *bool) {
	if *success && ctx.Err() == nil {
		s.pool.put(address, conn)
	} else {
		s.closeConn(conn, "scan completion")
	}
}

func (s *TCPScanner) dispatchResult(ctx context.Context, results chan<- models.Result, result *models.Result, target models.Target) {
	select {
	case results <- *result:
		if result.Available {
			log.Printf("Host %s has port %d open (%.2fms)",
				target.Host, target.Port, float64(result.RespTime.Microseconds())/1000)
		}
	case <-ctx.Done():
	}
}

// closeConn safely closes a connection and logs any errors.
func (*TCPScanner) closeConn(conn net.Conn, reason string) {
	if conn != nil {
		if err := conn.Close(); err != nil && !strings.Contains(err.Error(), "use of closed network connection") {
			log.Printf("Error closing connection (%s): %v", reason, err)
		}
	}
}

func (*TCPScanner) checkConnection(conn net.Conn) bool {
	tcpConn, ok := conn.(*net.TCPConn)
	if !ok {
		return true // Accept non-TCP connections as valid
	}

	// Set read deadline
	if err := tcpConn.SetReadDeadline(time.Now().Add(defaultReadDeadline)); err != nil {
		if !strings.Contains(err.Error(), "use of closed network connection") {
			log.Printf("Error setting read deadline: %v", err)
		}

		return false
	}

	// Reset deadline after we're done
	defer func() {
		if err := tcpConn.SetReadDeadline(time.Time{}); err != nil {
			if !strings.Contains(err.Error(), "use of closed network connection") {
				log.Printf("Error resetting read deadline: %v", err)
			}
		}
	}()

	// Try to read a single byte
	buf := make([]byte, 1)
	_, err := tcpConn.Read(buf)

	// Connection is considered valid if:
	// 1. Read succeeds (service sends banner)
	// 2. Read times out (most services wait for client)
	// 3. Connection is closed by remote (accepts and closes)
	if err != nil {
		var netErr net.Error
		if errors.As(err, &netErr) && netErr.Timeout() {
			return true
		}

		if errors.Is(err, io.EOF) {
			return true
		}

		return false
	}

	return true
}

// sendResultOrCleanup sends the result on the results channel. If the context is done,
// it ensures the connection is closed.
func (s *TCPScanner) sendResultOrCleanup(
	ctx context.Context, results chan<- models.Result, result *models.Result, conn net.Conn) {
	select {
	case results <- *result:
	case <-ctx.Done():
		s.closeConn(conn, "context cancellation during result send")
	}
}

// Scan performs TCP scans for multiple targets concurrently.
func (s *TCPScanner) Scan(ctx context.Context, targets []models.Target) (<-chan models.Result, error) {
	if len(targets) == 0 {
		results := make(chan models.Result)
		close(results)

		return results, nil
	}

	// Create buffered results channel
	results := make(chan models.Result, len(targets))

	// Create scan context
	scanCtx, cancel := context.WithCancel(ctx)

	// Create WaitGroup for all goroutines
	var wg sync.WaitGroup

	// Start processing goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		s.processTargets(scanCtx, targets, results)
	}()

	// Start cleanup goroutine
	go func() {
		// Set up done notification channels
		processDone := make(chan struct{})
		go func() {
			wg.Wait()
			close(processDone)
		}()

		// Wait for either completion or cancellation
		select {
		case <-processDone:
			log.Println("TCP scan completed successfully")
		case <-scanCtx.Done():
			log.Println("TCP scan stopping due to context cancellation")
		case <-s.done:
			log.Println("TCP scan stopping due to scanner shutdown")
		}

		// Cancel context and wait for processing to complete
		cancel()

		// Wait again with timeout to ensure cleanup
		cleanupTimer := time.NewTimer(s.timeout)
		select {
		case <-processDone:
			// Processing completed normally
		case <-cleanupTimer.C:
			log.Println("Warning: TCP scan cleanup timed out")
		}

		// Safe to close results channel now
		close(results)
	}()

	return results, nil
}

func (s *TCPScanner) processTargets(ctx context.Context, targets []models.Target, results chan<- models.Result) {
	log.Println("Starting TCP target processing")
	defer log.Println("TCP target processing completed")

	// Create semaphore for concurrency control
	sem := make(chan struct{}, s.concurrency)

	// Create WaitGroup for tracking target processing
	var targetWg sync.WaitGroup

	for _, target := range targets {
		select {
		case <-ctx.Done():
			return
		default:
		}

		targetWg.Add(1)
		go func(target models.Target) {
			defer targetWg.Done()

			// Acquire semaphore or return on context done
			select {
			case sem <- struct{}{}:
				defer func() { <-sem }()
			case <-ctx.Done():
				return
			}

			// Create connection with timeout
			connCtx, connCancel := context.WithTimeout(ctx, s.timeout)
			defer connCancel()

			result := models.Result{
				Target:    target,
				FirstSeen: time.Now(),
				LastSeen:  time.Now(),
			}

			address := fmt.Sprintf("%s:%d", target.Host, target.Port)
			conn, err := s.pool.get(connCtx, s.dialer, address)
			if err != nil {
				result.Error = err
				result.Available = false
				select {
				case results <- result:
				case <-ctx.Done():
				}
				return
			}

			success := false
			defer func() {
				if success && ctx.Err() == nil {
					s.pool.put(address, conn)
				} else {
					s.closeConn(conn, "scan completion")
				}
			}()

			startTime := time.Now()
			result.Available = s.checkConnection(conn)
			result.RespTime = time.Since(startTime)

			success = result.Available

			select {
			case results <- result:
				if result.Available {
					log.Printf("Host %s has port %d open (%.2fms)",
						target.Host, target.Port, float64(result.RespTime.Microseconds())/1000)
				}
			case <-ctx.Done():
			}
		}(target)
	}

	// Wait for all targets to complete
	targetWg.Wait()
}

func (s *TCPScanner) Stop(context.Context) error {
	log.Println("TCP Scanner Stop called")

	s.closeOnce.Do(func() {
		if s.done != nil {
			close(s.done)
		}

		if s.pool != nil {
			s.pool.close()
		}
	})

	return nil
}
