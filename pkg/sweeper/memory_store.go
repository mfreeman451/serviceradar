package sweeper

import (
	"context"
	"sync"
	"time"
)

// InMemoryStore implements Store interface for temporary storage.
type InMemoryStore struct {
	mu      sync.RWMutex
	results []Result
}

// NewInMemoryStore creates a new in-memory store for sweep results.
func NewInMemoryStore() Store {
	return &InMemoryStore{
		results: make([]Result, 0),
	}
}

// SaveHostResult updates the last-seen time (and possibly availability)
// for the given host. For in-memory store, we'll store the latest host
// result for each host.
func (s *InMemoryStore) SaveHostResult(_ context.Context, result *HostResult) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i := range s.results {
		// use index to avoid copying the entire Result struct
		existing := &s.results[i]

		// guard clause to skip non-matching hosts
		if existing.Target.Host != result.Host {
			continue
		}

		// Found matching host; update LastSeen and availability
		existing.LastSeen = result.LastSeen
		if result.Available {
			existing.Available = true
		}

		return nil
	}

	return nil
}

// GetHostResults returns a slice of HostResult based on the provided filter.
func (s *InMemoryStore) GetHostResults(_ context.Context, filter *ResultFilter) ([]HostResult, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// group results by host
	hostMap := make(map[string]*HostResult)

	for i := range s.results {
		r := &s.results[i]
		if !s.matchesFilter(r, filter) {
			continue
		}

		host, exists := hostMap[r.Target.Host]
		if !exists {
			host = &HostResult{
				Host:        r.Target.Host,
				FirstSeen:   r.FirstSeen,
				LastSeen:    r.LastSeen,
				Available:   false,
				PortResults: make([]*PortResult, 0),
			}
			hostMap[r.Target.Host] = host
		}

		if r.Available {
			host.Available = true
			if r.Target.Mode == ModeTCP {
				portResult := &PortResult{
					Port:      r.Target.Port,
					Available: true,
					RespTime:  r.RespTime,
				}
				host.PortResults = append(host.PortResults, portResult)
			}
		}

		// update timestamps
		if r.FirstSeen.Before(host.FirstSeen) {
			host.FirstSeen = r.FirstSeen
		}
		if r.LastSeen.After(host.LastSeen) {
			host.LastSeen = r.LastSeen
		}
	}

	// Convert map to slice
	hosts := make([]HostResult, 0, len(hostMap))
	for _, host := range hostMap {
		hosts = append(hosts, *host)
	}

	return hosts, nil
}

// GetSweepSummary gathers high-level sweep information.
func (s *InMemoryStore) GetSweepSummary(_ context.Context) (*SweepSummary, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	hostMap := make(map[string]*HostResult)
	portCounts := make(map[int]int)
	totalHosts := 0
	var lastSweep time.Time

	for i := range s.results {
		r := &s.results[i]

		// track latest sweep time
		if r.LastSeen.After(lastSweep) {
			lastSweep = r.LastSeen
		}

		// count ports
		if r.Available && r.Target.Mode == ModeTCP {
			portCounts[r.Target.Port]++
		}

		// group by host
		host, exists := hostMap[r.Target.Host]
		if !exists {
			totalHosts++
			host = &HostResult{
				Host:        r.Target.Host,
				FirstSeen:   r.FirstSeen,
				LastSeen:    r.LastSeen,
				Available:   false,
				PortResults: make([]*PortResult, 0),
			}
			hostMap[r.Target.Host] = host
		}

		if r.Available {
			host.Available = true
		}
	}

	// Count available hosts
	availableHosts := 0
	hosts := make([]HostResult, 0, len(hostMap))
	for _, host := range hostMap {
		if host.Available {
			availableHosts++
		}
		hosts = append(hosts, *host)
	}

	// Create port counts
	ports := make([]PortCount, 0, len(portCounts))
	for port, count := range portCounts {
		ports = append(ports, PortCount{
			Port:      port,
			Available: count,
		})
	}

	summary := &SweepSummary{
		Network:        "",
		TotalHosts:     totalHosts,
		AvailableHosts: availableHosts,
		LastSweep:      lastSweep.Unix(),
		Hosts:          hosts,
		Ports:          ports,
	}

	return summary, nil
}

// SaveResult stores (or updates) a Result in memory.
func (s *InMemoryStore) SaveResult(_ context.Context, result *Result) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i := range s.results {
		// if the same target already exists, overwrite
		if s.results[i].Target == result.Target {
			s.results[i] = *result
			return nil
		}
	}

	// If not found, append it
	s.results = append(s.results, *result)
	return nil
}

// GetResults returns a list of Results that match the filter.
func (s *InMemoryStore) GetResults(_ context.Context, filter *ResultFilter) ([]Result, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	filtered := make([]Result, 0, len(s.results))

	for i := range s.results {
		r := &s.results[i]
		if s.matchesFilter(r, filter) {
			filtered = append(filtered, *r)
		}
	}
	return filtered, nil
}

// PruneResults removes old results that haven't been seen since 'age' ago.
func (s *InMemoryStore) PruneResults(_ context.Context, age time.Duration) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	cutoff := time.Now().Add(-age)
	newResults := make([]Result, 0, len(s.results))

	for i := range s.results {
		r := &s.results[i]
		if r.LastSeen.After(cutoff) {
			newResults = append(newResults, *r)
		}
	}
	s.results = newResults

	return nil
}

// matchesFilter checks if a Result matches the provided filter.
func (*InMemoryStore) matchesFilter(result *Result, filter *ResultFilter) bool {
	checks := []func(*Result, *ResultFilter) bool{
		checkTimeRange,
		checkHost,
		checkPort,
		checkAvailability,
	}

	for _, check := range checks {
		if !check(result, filter) {
			return false
		}
	}
	return true
}

// checkTimeRange verifies if the result falls within the specified time range.
func checkTimeRange(result *Result, filter *ResultFilter) bool {
	if !filter.StartTime.IsZero() && result.LastSeen.Before(filter.StartTime) {
		return false
	}
	if !filter.EndTime.IsZero() && result.LastSeen.After(filter.EndTime) {
		return false
	}
	return true
}

// checkHost verifies if the result matches the specified host.
func checkHost(result *Result, filter *ResultFilter) bool {
	return filter.Host == "" || result.Target.Host == filter.Host
}

// checkPort verifies if the result matches the specified port.
func checkPort(result *Result, filter *ResultFilter) bool {
	return filter.Port == 0 || result.Target.Port == filter.Port
}

// checkAvailability verifies if the result matches the specified availability.
func checkAvailability(result *Result, filter *ResultFilter) bool {
	return filter.Available == nil || result.Available == *filter.Available
}
