// pkg/cloud/api/server.go

package api

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

// Generic status types
type ServiceStatus struct {
	Name      string          `json:"name"`
	Available bool            `json:"available"`
	Message   string          `json:"message"`
	Type      string          `json:"type"`    // e.g., "process", "port", "blockchain", etc.
	Details   json.RawMessage `json:"details"` // Flexible field for service-specific data
}

type NodeStatus struct {
	NodeID     string          `json:"node_id"`
	IsHealthy  bool            `json:"is_healthy"`
	LastUpdate time.Time       `json:"last_update"`
	Services   []ServiceStatus `json:"services"`
	UpTime     string          `json:"uptime"`
	FirstSeen  time.Time       `json:"first_seen"`
}

type SystemStatus struct {
	TotalNodes   int       `json:"total_nodes"`
	HealthyNodes int       `json:"healthy_nodes"`
	LastUpdate   time.Time `json:"last_update"`
}

type APIServer struct {
	mu            sync.RWMutex
	nodes         map[string]*NodeStatus
	statusHistory map[string][]NodeStatus
	maxHistory    int
	router        *mux.Router
}

func NewAPIServer() *APIServer {
	s := &APIServer{
		nodes:         make(map[string]*NodeStatus),
		statusHistory: make(map[string][]NodeStatus),
		maxHistory:    1000,
		router:        mux.NewRouter(),
	}
	s.setupRoutes()
	return s
}

func (s *APIServer) setupRoutes() {
	// Basic endpoints
	s.router.HandleFunc("/api/nodes", s.getNodes).Methods("GET")
	s.router.HandleFunc("/api/nodes/{id}", s.getNode).Methods("GET")
	s.router.HandleFunc("/api/status", s.getSystemStatus).Methods("GET")

	// History endpoint
	s.router.HandleFunc("/api/nodes/{id}/history", s.getNodeHistory).Methods("GET")

	// Service-specific endpoints
	s.router.HandleFunc("/api/nodes/{id}/services", s.getNodeServices).Methods("GET")
	s.router.HandleFunc("/api/nodes/{id}/services/{service}", s.getServiceDetails).Methods("GET")
}

func (s *APIServer) getNodeHistory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	nodeID := vars["id"]

	s.mu.RLock()
	history, exists := s.statusHistory[nodeID]
	s.mu.RUnlock()

	if !exists {
		http.Error(w, "Node not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(history)
}

func (s *APIServer) getNodes(w http.ResponseWriter, r *http.Request) {
	s.mu.RLock()
	nodes := make([]*NodeStatus, 0, len(s.nodes))
	for _, node := range s.nodes {
		nodes = append(nodes, node)
	}
	s.mu.RUnlock()

	json.NewEncoder(w).Encode(nodes)
}

func (s *APIServer) getNode(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	nodeID := vars["id"]

	s.mu.RLock()
	node, exists := s.nodes[nodeID]
	s.mu.RUnlock()

	if !exists {
		http.Error(w, "Node not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(node)
}

func (s *APIServer) getNodeServices(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	nodeID := vars["id"]

	s.mu.RLock()
	node, exists := s.nodes[nodeID]
	s.mu.RUnlock()

	if !exists {
		http.Error(w, "Node not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(node.Services)
}

func (s *APIServer) getServiceDetails(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	nodeID := vars["id"]
	serviceName := vars["service"]

	s.mu.RLock()
	node, exists := s.nodes[nodeID]
	s.mu.RUnlock()

	if !exists {
		http.Error(w, "Node not found", http.StatusNotFound)
		return
	}

	for _, service := range node.Services {
		if service.Name == serviceName {
			json.NewEncoder(w).Encode(service)
			return
		}
	}

	http.Error(w, "Service not found", http.StatusNotFound)
}

func (s *APIServer) getSystemStatus(w http.ResponseWriter, r *http.Request) {
	s.mu.RLock()
	status := SystemStatus{
		TotalNodes: len(s.nodes),
		LastUpdate: time.Now(),
	}

	for _, node := range s.nodes {
		if node.IsHealthy {
			status.HealthyNodes++
		}
	}
	s.mu.RUnlock()

	json.NewEncoder(w).Encode(status)
}

func (s *APIServer) UpdateNodeStatus(nodeID string, status *NodeStatus) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.nodes[nodeID] = status

	// Update history
	s.statusHistory[nodeID] = append(s.statusHistory[nodeID], *status)
	if len(s.statusHistory[nodeID]) > s.maxHistory {
		s.statusHistory[nodeID] = s.statusHistory[nodeID][1:]
	}
}

func (s *APIServer) Start(addr string) error {
	return http.ListenAndServe(addr, s.router)
}
