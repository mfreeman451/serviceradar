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

// Package models provides data models for the sweeper service.
package models

import "time"

// SweepData represents network sweep results.
type SweepData struct {
	Network        string       `json:"network"`
	TotalHosts     int32        `json:"total_hosts"`
	AvailableHosts int32        `json:"available_hosts"`
	LastSweep      int64        `json:"last_sweep"`
	Ports          []PortStatus `json:"ports"`
}

// PortStatus represents port availability information.
type PortStatus struct {
	Port      int32 `json:"port"`
	Available int32 `json:"available"`
}

// Config defines sweeper configuration.
type Config struct {
	Networks     []string      `json:"networks"`
	Ports        []int         `json:"ports"`
	SweepModes   []SweepMode   `json:"sweep_modes"`
	Interval     time.Duration `json:"interval"`
	Concurrency  int           `json:"concurrency"`
	Timeout      time.Duration `json:"timeout"`
	ICMPCount    int           `json:"icmp_count"`
	MaxIdle      int           `json:"max_idle"`
	MaxLifetime  time.Duration `json:"max_lifetime"`
	IdleTimeout  time.Duration `json:"idle_timeout"`
	ICMPSettings struct {
		RateLimit int // Packets per second
		Timeout   time.Duration
		MaxBatch  int
	}
	TCPSettings struct {
		Concurrency int
		Timeout     time.Duration
		MaxBatch    int
	}
	EnableHighPerformanceICMP bool `json:"high_perf_icmp,omitempty"`
	ICMPRateLimit             int  `json:"icmp_rate_limit,omitempty"`
}

type SweepMode string

const (
	ModeTCP  SweepMode = "tcp"
	ModeICMP SweepMode = "icmp"
)

// Target represents a network target to be scanned.
type Target struct {
	Host     string
	Port     int
	Mode     SweepMode
	Metadata map[string]interface{} // Additional metadata about the scan

}

// Result represents the outcome of a sweep against a target.
type Result struct {
	Target     Target
	Available  bool
	FirstSeen  time.Time
	LastSeen   time.Time
	RespTime   time.Duration
	PacketLoss float64
	Error      error
	Metadata   map[string]interface{}
}

// ResultFilter defines criteria for retrieving results.
type ResultFilter struct {
	Host      string
	Port      int
	StartTime time.Time
	EndTime   time.Time
	Available *bool
}

// ContainsMode checks if a mode is in a list of modes.
func ContainsMode(modes []SweepMode, mode SweepMode) bool {
	for _, m := range modes {
		if m == mode {
			return true
		}
	}

	return false
}

// HostResult represents all results for a single host.
type HostResult struct {
	Host         string        `json:"host"`
	Available    bool          `json:"available"`
	FirstSeen    time.Time     `json:"first_seen"`
	LastSeen     time.Time     `json:"last_seen"`
	PortResults  []*PortResult `json:"port_results,omitempty"`
	ICMPStatus   *ICMPStatus   `json:"icmp_status,omitempty"`
	ResponseTime time.Duration `json:"response_time"`
}

// ICMPStatus represents ICMP ping results.
type ICMPStatus struct {
	Available  bool          `json:"available"`
	RoundTrip  time.Duration `json:"round_trip"`
	PacketLoss float64       `json:"packet_loss"`
}

// PortResult represents a single port scan result.
type PortResult struct {
	Port      int           `json:"port"`
	Available bool          `json:"available"`
	RespTime  time.Duration `json:"response_time"`
	Service   string        `json:"service,omitempty"` // Optional service identification
}
type PortCount struct {
	Port      int `json:"port"`
	Available int `json:"available"`
}

// SweepSummary provides aggregated sweep results.
type SweepSummary struct {
	Network        string       `json:"network"`
	TotalHosts     int          `json:"total_hosts"`
	AvailableHosts int          `json:"available_hosts"`
	LastSweep      int64        `json:"last_sweep"` // Unix timestamp
	Ports          []PortCount  `json:"ports"`
	Hosts          []HostResult `json:"hosts"`
}
