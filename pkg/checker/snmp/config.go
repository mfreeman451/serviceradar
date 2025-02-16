package snmp

import (
	"fmt"
	"time"

	"github.com/mfreeman451/serviceradar/pkg/config"
)

const (
	defaultTimeout = 5 * time.Minute
)

var (
	errNodeAddressRequired = fmt.Errorf("node_address is required")
	errListenAddrRequired  = fmt.Errorf("listen_addr is required")
)

// Validate implements config.Validator interface.
func (c *Config) Validate() error {
	if c.NodeAddress == "" {
		return errNodeAddressRequired
	}

	if c.ListenAddr == "" {
		return errListenAddrRequired
	}

	// Cast to time.Duration for comparison
	if time.Duration(c.Timeout) == 0 {
		// Assign a default by constructing a config.Duration
		c.Timeout = config.Duration(defaultTimeout)
	}

	return nil
}
