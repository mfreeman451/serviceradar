// Package snmp pkg/checker/snmp/client.go

package snmp

import (
	"fmt"
	"sync"
	"time"

	"github.com/gosnmp/gosnmp"
)

// SNMPClientImpl implements the SNMPClient interface using gosnmp.
type SNMPClientImpl struct {
	client     *gosnmp.GoSNMP
	target     *Target
	mu         sync.RWMutex
	connected  bool
	lastError  error
	reconnects int
}

// SNMPError wraps SNMP-specific errors with additional context.
type SNMPError struct {
	Op      string
	Target  string
	Wrapped error
}

func (e *SNMPError) Error() string {
	return fmt.Sprintf("SNMP %s failed for target %s: %v", e.Op, e.Target, e.Wrapped)
}

func newSNMPClient(target *Target) (SNMPClient, error) {
	if err := validateTarget(target); err != nil {
		return nil, fmt.Errorf("invalid target: %w", err)
	}

	client := &gosnmp.GoSNMP{
		Target:             target.Host,
		Port:               target.Port,
		Community:          target.Community,
		Timeout:            time.Duration(target.Timeout),
		Retries:            target.Retries,
		ExponentialTimeout: true,
		MaxOids:            gosnmp.MaxOids,
	}

	// Set SNMP version based on configuration
	switch target.Version {
	case Version1:
		client.Version = gosnmp.Version1
	case Version2c:
		client.Version = gosnmp.Version2c
	case Version3:
		return nil, fmt.Errorf("SNMPv3 not yet implemented")
	default:
		return nil, fmt.Errorf("unsupported SNMP version: %s", target.Version)
	}

	return &SNMPClientImpl{
		client: client,
		target: target,
	}, nil
}

// Connect implements SNMPClient interface.
func (s *SNMPClientImpl) Connect() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.connected {
		return nil
	}

	if err := s.client.Connect(); err != nil {
		s.lastError = &SNMPError{
			Op:      "connect",
			Target:  s.target.Host,
			Wrapped: err,
		}

		return s.lastError
	}

	s.connected = true

	return nil
}

// Get implements SNMPClient interface.
func (s *SNMPClientImpl) Get(oids []string) (map[string]interface{}, error) {
	s.mu.Lock()
	if !s.connected {
		if err := s.client.Connect(); err != nil {
			s.mu.Unlock()

			return nil, &SNMPError{
				Op:      "connect",
				Target:  s.target.Host,
				Wrapped: err,
			}
		}

		s.connected = true
	}

	s.mu.Unlock()

	// Split OIDs into chunks of MaxOids size
	var allResults = make(map[string]interface{})

	for i := 0; i < len(oids); i += gosnmp.MaxOids {
		end := i + gosnmp.MaxOids
		if end > len(oids) {
			end = len(oids)
		}

		chunk := oids[i:end]

		result, err := s.client.Get(chunk)
		if err != nil {
			s.handleError(err)

			return nil, &SNMPError{
				Op:      "get",
				Target:  s.target.Host,
				Wrapped: err,
			}
		}

		// Process results from this chunk
		for _, variable := range result.Variables {
			value, err := s.convertVariable(variable)
			if err != nil {
				return nil, &SNMPError{
					Op:      "convert",
					Target:  s.target.Host,
					Wrapped: err,
				}
			}

			allResults[variable.Name] = value
		}
	}

	return allResults, nil
}

// Close implements SNMPClient interface.
func (s *SNMPClientImpl) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.connected {
		return nil
	}

	err := s.client.Conn.Close()
	if err != nil {
		return err
	}

	s.connected = false

	return nil
}

// GetLastError returns the last error encountered.
func (s *SNMPClientImpl) GetLastError() error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.lastError
}

// handleError processes SNMP errors and manages reconnection.
func (s *SNMPClientImpl) handleError(err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.lastError = err
	s.connected = false
	s.reconnects++
}

const defaultTimeTickDuration = time.Second / 100

// convertVariable converts an SNMP variable to the appropriate Go type.
func (*SNMPClientImpl) convertVariable(variable gosnmp.SnmpPDU) (interface{}, error) {
	switch variable.Type {
	// Handling the known types in gosnmp.Asn1BER
	case gosnmp.Boolean:
		return variable.Value.(bool), nil
	case gosnmp.BitString:
		return variable.Value, nil // Needs custom decoding
	case gosnmp.Null:
		return nil, nil
	case gosnmp.ObjectDescription:
		return string(variable.Value.(byte)), nil
	case gosnmp.Opaque:
		return variable.Value, nil // Needs custom decoding
	case gosnmp.NsapAddress:
		return variable.Value, nil
	case gosnmp.Uinteger32:
		return uint64(variable.Value.(uint)), nil
	case gosnmp.OpaqueFloat:
		return variable.Value.(float32), nil
	case gosnmp.OpaqueDouble:
		return variable.Value.(float64), nil
	case gosnmp.NoSuchObject:
		return nil, fmt.Errorf("SNMP NoSuchObject")
	case gosnmp.NoSuchInstance:
		return nil, fmt.Errorf("SNMP NoSuchInstance")
	case gosnmp.EndOfMibView:
		return nil, fmt.Errorf("SNMP EndOfMibView")
	// Removed duplicate case for gosnmp.UnknownType
	case gosnmp.UnknownType:
		return nil, fmt.Errorf("SNMP UnknownType")
	// Handling other types
	case gosnmp.Integer:
		return variable.Value.(int), nil
	case gosnmp.OctetString:
		return string(variable.Value.(byte)), nil
	case gosnmp.ObjectIdentifier:
		return variable.Value.(string), nil
	case gosnmp.IPAddress:
		return variable.Value.(string), nil
	case gosnmp.Counter32, gosnmp.Gauge32:
		return uint64(variable.Value.(uint)), nil
	case gosnmp.Counter64:
		return variable.Value.(uint64), nil
	case gosnmp.TimeTicks:
		return time.Duration(variable.Value.(uint32)) * defaultTimeTickDuration, nil
	// Default case for unsupported types
	default:
		return nil, fmt.Errorf("unsupported SNMP type: %v, Value: %v", variable.Type, variable.Value)
	}
}

const defaultSNMPPort = 161

// validateTarget performs basic validation of target configuration.
func validateTarget(target *Target) error {
	if target == nil {
		return fmt.Errorf("target configuration is nil")
	}

	if target.Host == "" {
		return fmt.Errorf("target host is required")
	}

	if target.Port == 0 {
		target.Port = defaultPort
	}

	if target.Timeout == 0 {
		target.Timeout = Duration(defaultTimeout)
	}

	if target.Retries == 0 {
		target.Retries = defaultRetries
	}

	return nil
}
