// pkg/cloud/snmp_handler.go

package cloud

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/mfreeman451/serviceradar/pkg/cloud/api"
	"github.com/mfreeman451/serviceradar/pkg/db"
)

type SNMPServiceStatus struct {
	NodeID    string                    `json:"node_id"`
	LastPoll  time.Time                 `json:"last_poll"`
	OIDStatus map[string]OIDStatusEntry `json:"oid_status"`
	Available bool                      `json:"available"`
}

type OIDStatusEntry struct {
	LastValue  interface{} `json:"last_value"`
	LastUpdate time.Time   `json:"last_update"`
	ErrorCount int         `json:"error_count"`
	LastError  string      `json:"last_error,omitempty"`
	Type       string      `json:"type"`
}

func (s *Server) handleSNMPService(nodeID string, service *api.ServiceStatus) error {
	var snmpStatus SNMPServiceStatus
	if err := json.Unmarshal([]byte(service.Message), &snmpStatus); err != nil {
		return fmt.Errorf("failed to parse SNMP status: %w", err)
	}

	// Update API state with processed SNMP data
	apiStatus := &api.NodeStatus{
		NodeID:     nodeID,
		LastUpdate: time.Now(),
		IsHealthy:  snmpStatus.Available,
		Services: []api.ServiceStatus{
			{
				Name:      service.Name,
				Type:      "snmp",
				Available: snmpStatus.Available,
				Details:   json.RawMessage(service.Message),
			},
		},
	}

	s.updateAPIState(nodeID, apiStatus)

	// Store historical data if needed
	if err := s.storeNodeStatus(nodeID, snmpStatus.Available, time.Now()); err != nil {
		return fmt.Errorf("failed to store SNMP status: %w", err)
	}

	return nil
}

func (s *Server) processSNMPData(nodeID string, status *api.ServiceStatus) error {
	var snmpData SNMPServiceStatus
	if err := json.Unmarshal([]byte(status.Message), &snmpData); err != nil {
		return fmt.Errorf("failed to unmarshal SNMP data: %w", err)
	}

	// Calculate health based on OID errors
	totalErrors := 0
	for _, oidStatus := range snmpData.OIDStatus {
		totalErrors += oidStatus.ErrorCount
	}

	// Consider the service unhealthy if there are too many errors
	isHealthy := totalErrors < len(snmpData.OIDStatus)*3 // Allow up to 3 errors per OID

	// Store the processed status
	if err := s.db.UpdateServiceStatus(&db.ServiceStatus{
		NodeID:      nodeID,
		ServiceName: status.Name,
		ServiceType: "snmp",
		Available:   isHealthy,
		Details:     status.Message,
		Timestamp:   time.Now(),
	}); err != nil {
		return fmt.Errorf("failed to store SNMP service status: %w", err)
	}

	return nil
}
