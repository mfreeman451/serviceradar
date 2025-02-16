// Package snmp pkg/checker/snmp/collector.go

package snmp

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

const (
	defaultByteBuffer               = 1024
	defaultErrorChan                = 10
	defaultDataChanBufferMultiplier = 2
)

// SNMPCollector implements the Collector interface.
type SNMPCollector struct {
	target     *Target
	client     SNMPClient
	dataChan   chan DataPoint
	errorChan  chan error
	done       chan struct{}
	closeOnce  sync.Once
	mu         sync.RWMutex
	status     TargetStatus
	bufferPool *sync.Pool
}

// NewCollector creates a new SNMP collector for a target.
func NewCollector(target *Target) (Collector, error) {
	if err := validateTarget(target); err != nil {
		return nil, fmt.Errorf("%w %w", ErrInvalidTargetConfig, err)
	}

	// Initialize the SNMP client
	client, err := newSNMPClient(target)
	if err != nil {
		return nil, fmt.Errorf("%w %w", ErrSNMPConnect, err)
	}

	collector := &SNMPCollector{
		target:    target,
		client:    client,
		dataChan:  make(chan DataPoint, len(target.OIDs)*defaultDataChanBufferMultiplier), // Buffer for 2 polls per OID
		errorChan: make(chan error, defaultErrorChan),
		done:      make(chan struct{}),
		status: TargetStatus{
			OIDStatus: make(map[string]OIDStatus),
		},
		bufferPool: &sync.Pool{
			New: func() interface{} {
				return make([]byte, 0, defaultByteBuffer)
			},
		},
	}

	return collector, nil
}

// Start implements the Collector interface.
func (c *SNMPCollector) Start(ctx context.Context) error {
	// Connect to the SNMP device
	if err := c.client.Connect(); err != nil {
		return fmt.Errorf("%w - %w", ErrSNMPConnect, err)
	}

	// Start collection goroutine
	go c.collect(ctx)

	// Start error handling goroutine
	go c.handleErrors(ctx)

	return nil
}

// Stop implements the Collector interface.
func (c *SNMPCollector) Stop() error {
	c.closeOnce.Do(func() {
		close(c.done)

		if err := c.client.Close(); err != nil {
			log.Printf("Error closing SNMP client for target %s: %v", c.target.Name, err)
		}
	})

	return nil
}

// GetResults implements the Collector interface.
func (c *SNMPCollector) GetResults() <-chan DataPoint {
	return c.dataChan
}

// collect runs the main collection loop.
func (c *SNMPCollector) collect(ctx context.Context) {
	ticker := time.NewTicker(time.Duration(c.target.Interval))
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-c.done:
			return
		case <-ticker.C:
			if err := c.pollTarget(ctx); err != nil {
				select {
				case c.errorChan <- err:
				default:
					log.Printf("Error channel full, dropping error: %v", err)
				}
			}
		}
	}
}

// pollTarget performs a single poll of all OIDs for the target.
func (c *SNMPCollector) pollTarget(ctx context.Context) error {
	log.Printf("Polling target %s (%s) for %d OIDs", c.target.Name, c.target.Host, len(c.target.OIDs))

	oids := make([]string, len(c.target.OIDs))
	for i, oid := range c.target.OIDs {
		oids[i] = oid.OID
	}

	// Get SNMP data
	results, err := c.client.Get(oids)
	if err != nil {
		c.updateStatus(false, err.Error())
		return fmt.Errorf("%w - %w", ErrSNMPGet, err)
	}

	log.Printf("Successfully polled target %s, processing %d results", c.target.Name, len(results))
	c.updateStatus(true, "")

	// Process each result
	for oid, value := range results {
		if err := c.processResult(ctx, oid, value); err != nil {
			log.Printf("Error processing result for OID %s: %v", oid, err)
		}
	}

	return nil
}

// processResult handles a single OID result.
func (c *SNMPCollector) processResult(ctx context.Context, oid string, value interface{}) error {
	// Find OID config
	var oidConfig *OIDConfig

	for _, cfg := range c.target.OIDs {
		if cfg.OID == oid {
			oidConfig = &cfg

			break
		}
	}

	if oidConfig == nil {
		return fmt.Errorf("%w %s", ErrNoOIDConfig, oid)
	}

	// Convert value based on type
	converted, err := c.convertValue(value, oidConfig)
	if err != nil {
		return fmt.Errorf("%w - %w", ErrSNMPConvert, err)
	}

	// Create data point
	point := DataPoint{
		OIDName:   oidConfig.Name,
		Value:     converted,
		Timestamp: time.Now(),
	}

	log.Printf("Collected data point for %s: %v", point.OIDName, point.Value)

	// Update OID status
	c.updateOIDStatus(oidConfig.Name, point)

	// Send data point
	select {
	case c.dataChan <- point:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	case <-c.done:
		return ErrCollectorStopped
	}
}

// convertValue converts an SNMP value based on the OID configuration.
func (c *SNMPCollector) convertValue(value interface{}, config *OIDConfig) (interface{}, error) {
	switch config.DataType {
	case TypeCounter:
		return c.convertCounter(value, config.Scale)
	case TypeGauge:
		return c.convertGauge(value, config.Scale)
	case TypeBoolean:
		return c.convertBoolean(value)
	case TypeBytes:
		return c.convertBytes(value, config.Scale)
	case TypeString:
		return c.convertString(value)
	default:
		return nil, fmt.Errorf("%w %v", ErrUnsupportedDataType, config.DataType)
	}
}

// handleErrors processes errors from the collection process.
func (c *SNMPCollector) handleErrors(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-c.done:
			return
		case err := <-c.errorChan:
			log.Printf("Error collecting from target %s: %v", c.target.Name, err)
		}
	}
}

// updateStatus updates the collector's status.
func (c *SNMPCollector) updateStatus(available bool, errorMsg string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.status.Available = available
	c.status.LastPoll = time.Now()
	c.status.Error = errorMsg
}

// updateOIDStatus updates the status for a specific OID.
func (c *SNMPCollector) updateOIDStatus(oidName string, point DataPoint) {
	c.mu.Lock()
	defer c.mu.Unlock()

	status := c.status.OIDStatus[oidName]
	status.LastValue = point.Value
	status.LastUpdate = point.Timestamp

	c.status.OIDStatus[oidName] = status
}

// GetStatus returns the current status of the collector.
func (c *SNMPCollector) GetStatus() TargetStatus {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.status
}

// convertCounter converts a counter value to a uint64.
func (*SNMPCollector) convertCounter(value interface{}, scale float64) (uint64, error) {
	v, ok := value.(uint64)
	if !ok {
		return 0, fmt.Errorf("%w %T", ErrInvalidCounterType, value)
	}

	return uint64(float64(v) * scale), nil
}

// convertGauge converts a gauge value to a float64.
func (*SNMPCollector) convertGauge(value interface{}, scale float64) (float64, error) {
	switch v := value.(type) {
	case uint64:
		return float64(v) * scale, nil
	case int64:
		return float64(v) * scale, nil
	case float64:
		return v * scale, nil
	default:
		return 0, fmt.Errorf("%w %T", ErrInvalidGaugeType, value)
	}
}

// convertBoolean converts a boolean value to a bool.
func (*SNMPCollector) convertBoolean(value interface{}) (bool, error) {
	switch v := value.(type) {
	case int:
		return v != 0, nil
	case bool:
		return v, nil
	default:
		return false, fmt.Errorf("%w %T", ErrInvalidBooleanType, value)
	}
}

// convertBytes converts a byte value to a uint64.
func (*SNMPCollector) convertBytes(value interface{}, scale float64) (uint64, error) {
	v, ok := value.(uint64)
	if !ok {
		return 0, fmt.Errorf("%w %T", ErrInvalidBytesType, value)
	}

	return uint64(float64(v) * scale), nil
}

// convertString converts a string value to a string.
func (*SNMPCollector) convertString(value interface{}) (string, error) {
	switch v := value.(type) {
	case []byte:
		return string(v), nil
	case string:
		return v, nil
	default:
		return "", fmt.Errorf("%w %T", ErrInvalidStringType, value)
	}
}
