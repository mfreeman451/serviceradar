/*
 * Copyright 2025 Carver Automation Corporation.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package scan

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"sync"
	"syscall"
	"time"

	"github.com/carverauto/serviceradar/pkg/models"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

const (
	defaultICMPRateLimit = 1000 // packets per second
	defaultICMPTimeout   = 5 * time.Second
	batchInterval        = 10 * time.Millisecond
)

type ICMPSweeper struct {
	rateLimit   int
	timeout     time.Duration
	identifier  int
	rawSocketFD int
	conn        *icmp.PacketConn
	mu          sync.Mutex
	results     map[string]models.Result
	cancel      context.CancelFunc
}

var _ Scanner = (*ICMPSweeper)(nil)

const (
	defaultIdentifierMod = 65536
)

// NewICMPSweeper creates a new scanner for ICMP sweeping.
func NewICMPSweeper(timeout time.Duration, rateLimit int) (*ICMPSweeper, error) {
	if timeout == 0 {
		timeout = defaultICMPTimeout
	}

	if rateLimit == 0 {
		rateLimit = defaultICMPRateLimit
	}

	// Create identifier for this scanner instance
	identifier := int(time.Now().UnixNano() % defaultIdentifierMod)

	// Create raw socket for sending
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_ICMP)
	if err != nil {
		return nil, fmt.Errorf("failed to create raw socket: %w", err)
	}

	// Create listener for receiving
	conn, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		err := syscall.Close(fd)
		if err != nil {
			log.Printf("Failed to close ICMP listener: %v", err)
			return nil, err
		}

		return nil, fmt.Errorf("failed to create ICMP listener: %w", err)
	}

	s := &ICMPSweeper{
		rateLimit:   rateLimit,
		timeout:     timeout,
		identifier:  identifier,
		rawSocketFD: fd,
		conn:        conn,
		results:     make(map[string]models.Result),
	}

	return s, nil
}

// Scan performs the ICMP sweep and returns results.
func (s *ICMPSweeper) Scan(ctx context.Context, targets []models.Target) (<-chan models.Result, error) {
	icmpTargets := filterICMPTargets(targets)

	if len(icmpTargets) == 0 {
		ch := make(chan models.Result)
		close(ch)

		return ch, nil
	}

	scanCtx, cancel := context.WithCancel(ctx)
	s.cancel = cancel

	resultCh := make(chan models.Result, len(icmpTargets))

	// Reset results map for this scan
	s.mu.Lock()
	s.results = make(map[string]models.Result)
	s.mu.Unlock()

	// Start listener goroutine
	listenerDone := make(chan struct{})

	go func() {
		defer close(listenerDone)
		s.listenForReplies(scanCtx, icmpTargets)
	}()

	// Start sender goroutine
	senderDone := make(chan struct{})

	go func() {
		defer close(senderDone)

		s.sendPings(scanCtx, icmpTargets)
	}()

	// Process results after scanning is done or timeout
	go func() {
		defer close(resultCh)

		// Wait for sender to finish or context to be canceled
		select {
		case <-senderDone:
			// Wait for all replies or timeout
			timer := time.NewTimer(s.timeout)
			select {
			case <-timer.C:
				// Timeout reached
			case <-scanCtx.Done():
				if !timer.Stop() {
					<-timer.C
				}
			}
		case <-scanCtx.Done():
		}
		cancel()
		<-listenerDone

		// Process and send results
		s.processResults(icmpTargets, resultCh)
	}()

	return resultCh, nil
}

// sendPings sends ICMP echo requests to all targets with rate limiting.
func (s *ICMPSweeper) sendPings(ctx context.Context, targets []models.Target) {
	packetsPerInterval := s.calculatePacketsPerInterval()

	log.Printf("Sending ICMP pings to %d targets (rate: %d/sec, batch: %d)",
		len(targets), s.rateLimit, packetsPerInterval)

	data, err := s.prepareEchoRequest()
	if err != nil {
		log.Printf("Error marshaling ICMP message: %v", err)

		return
	}

	s.sendBatches(ctx, targets, data, packetsPerInterval)
}

const (
	defaultRateLimitDivisor = 1000
)

// calculatePacketsPerInterval determines the batch size based on rate limit.
func (s *ICMPSweeper) calculatePacketsPerInterval() int {
	packets := s.rateLimit / int(defaultRateLimitDivisor/batchInterval.Milliseconds())
	if packets < 1 {
		return 1
	}

	return packets
}

// prepareEchoRequest builds the ICMP echo request template.
func (s *ICMPSweeper) prepareEchoRequest() ([]byte, error) {
	msg := icmp.Message{
		Type: ipv4.ICMPTypeEcho,
		Code: 0,
		Body: &icmp.Echo{
			ID:   s.identifier,
			Seq:  1,
			Data: []byte("ping"),
		},
	}

	return msg.Marshal(nil)
}

// sendBatches manages the sending of ping batches.
func (s *ICMPSweeper) sendBatches(ctx context.Context, targets []models.Target, data []byte, batchSize int) {
	ticker := time.NewTicker(batchInterval)
	defer ticker.Stop()

	targetIndex := 0

	for range ticker.C {
		if ctx.Err() != nil {
			return
		}

		batchEnd := s.calculateBatchEnd(targetIndex, batchSize, len(targets))
		s.processBatch(targets[targetIndex:batchEnd], data)

		targetIndex = batchEnd
		if targetIndex >= len(targets) {
			return
		}
	}
}

// calculateBatchEnd determines the end index for the current batch.
func (*ICMPSweeper) calculateBatchEnd(index, batchSize, totalTargets int) int {
	end := index + batchSize
	if end > totalTargets {
		return totalTargets
	}

	return end
}

// processBatch sends pings to a batch of targets.
func (s *ICMPSweeper) processBatch(targets []models.Target, data []byte) {
	for _, target := range targets {
		s.sendPingToTarget(target, data)
	}
}

// sendPingToTarget sends a single ICMP ping and records initial result.
func (s *ICMPSweeper) sendPingToTarget(target models.Target, data []byte) {
	ipAddr := net.ParseIP(target.Host)
	if ipAddr == nil || ipAddr.To4() == nil {
		log.Printf("Invalid IPv4 address: %s", target.Host)

		return
	}

	addr := [4]byte{}
	copy(addr[:], ipAddr.To4())
	sockaddr := &syscall.SockaddrInet4{Addr: addr}

	if err := syscall.Sendto(s.rawSocketFD, data, 0, sockaddr); err != nil {
		log.Printf("Error sending ICMP to %s: %v", target.Host, err)
	}

	s.recordInitialResult(target)
}

// recordInitialResult stores the initial ping result.
func (s *ICMPSweeper) recordInitialResult(target models.Target) {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	s.results[target.Host] = models.Result{
		Target:     target,
		Available:  false,
		FirstSeen:  now,
		LastSeen:   now,
		PacketLoss: 100,
	}
	fmt.Println(s.results[target.Host])
}

const (
	defaultBytesRead    = 1500
	defaultReadDeadline = 100 * time.Millisecond
)

// listenForReplies listens for and processes ICMP echo replies.
func (s *ICMPSweeper) listenForReplies(ctx context.Context, targets []models.Target) {
	targetMap := make(map[string]struct{})
	for _, t := range targets {
		targetMap[t.Host] = struct{}{}
	}

	buf := make([]byte, defaultBytesRead)

	const readDeadline = defaultReadDeadline

	for {
		if ctx.Err() != nil {
			return
		}

		if err := s.conn.SetReadDeadline(time.Now().Add(readDeadline)); err != nil {
			log.Printf("Error setting read deadline: %v", err)

			continue
		}

		reply, err := s.readReply(buf)
		if err != nil {
			continue // Error handling already logged in readReply
		}

		if err := s.processReply(reply, targetMap); err != nil {
			continue // Error handling already logged in processReply
		}
	}
}

// readReply reads an ICMP reply from the connection.
func (s *ICMPSweeper) readReply(buf []byte) (reply struct {
	n    int
	addr net.Addr
	data []byte
}, err error) {
	n, addr, err := s.conn.ReadFrom(buf)
	if err != nil {
		var netErr net.Error

		if errors.As(err, &netErr) && netErr.Timeout() {
			return reply, nil // Timeout is not an error in this context
		}

		log.Printf("Error reading ICMP reply: %v", err)

		return reply, err
	}

	return struct {
		n    int
		addr net.Addr
		data []byte
	}{n, addr, buf[:n]}, nil
}

// processReply processes a valid ICMP reply.
func (s *ICMPSweeper) processReply(reply struct {
	n    int
	addr net.Addr
	data []byte
}, targetMap map[string]struct{}) error {
	if reply.addr == nil {
		// Timeout or invalid reply, skip processing
		return nil
	}

	ip := reply.addr.String()

	// Verify this is one of our targets
	if _, ok := targetMap[ip]; !ok {
		return nil // Not an error, just not our target
	}

	// Parse the ICMP message
	msg, err := icmp.ParseMessage(1, reply.data)
	if err != nil {
		log.Printf("Error parsing ICMP message from %s: %v", ip, err)
		return err
	}

	// Verify it's an echo reply with our identifier
	echo, ok := msg.Body.(*icmp.Echo)
	if !ok || msg.Type != ipv4.ICMPTypeEchoReply || echo.ID != s.identifier {
		return nil // Not an error, just not our reply
	}

	// Update the result
	s.mu.Lock()
	defer s.mu.Unlock()

	if result, ok := s.results[ip]; ok {
		result.Available = true
		result.RespTime = time.Since(result.FirstSeen)
		result.PacketLoss = 0
		result.LastSeen = time.Now()
		s.results[ip] = result
	}

	return nil
}

// processResults sends final results to the result channel.
func (s *ICMPSweeper) processResults(targets []models.Target, ch chan<- models.Result) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Send all results to the channel
	for _, target := range targets {
		if result, ok := s.results[target.Host]; ok {
			ch <- result
		} else {
			// If we somehow don't have a result for this target, create a default one
			ch <- models.Result{
				Target:     target,
				Available:  false,
				PacketLoss: 100,
				FirstSeen:  time.Now(),
				LastSeen:   time.Now(),
			}
		}
	}
}

// Stop stops the scanner and releases resources.
func (s *ICMPSweeper) Stop(_ context.Context) error {
	if s.cancel != nil {
		s.cancel()
	}

	// Close the connection and socket
	if s.conn != nil {
		err := s.conn.Close()
		if err != nil {
			log.Printf("Error closing ICMP connection: %v", err)

			return err
		}
	}

	if s.rawSocketFD != 0 {
		err := syscall.Close(s.rawSocketFD)
		if err != nil {
			log.Printf("Error closing raw socket: %v", err)

			return err
		}

		s.rawSocketFD = 0
	}

	return nil
}

// filterICMPTargets filters only ICMP targets from the given slice.
func filterICMPTargets(targets []models.Target) []models.Target {
	var filtered []models.Target

	for _, t := range targets {
		if t.Mode == models.ModeICMP {
			filtered = append(filtered, t)
		}
	}

	return filtered
}
