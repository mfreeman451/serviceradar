// pkg/cloud/api/server.go

package api

import (
	"crypto/rand"
	"embed"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/carverauto/serviceradar/pkg/checker/snmp"
	"github.com/carverauto/serviceradar/pkg/db"
	srHttp "github.com/carverauto/serviceradar/pkg/http"
	"github.com/carverauto/serviceradar/pkg/metrics"
	"github.com/carverauto/serviceradar/pkg/models"
	"github.com/gorilla/mux"
)

type ServiceStatus struct {
	Name      string          `json:"name"`
	Available bool            `json:"available"`
	Message   string          `json:"message"`
	Type      string          `json:"type"`    // e.g., "process", "port", "blockchain", etc.
	Details   json.RawMessage `json:"details"` // Flexible field for service-specific data
}

type NodeStatus struct {
	NodeID     string               `json:"node_id"`
	IsHealthy  bool                 `json:"is_healthy"`
	LastUpdate time.Time            `json:"last_update"`
	Services   []ServiceStatus      `json:"services"`
	UpTime     string               `json:"uptime"`
	FirstSeen  time.Time            `json:"first_seen"`
	Metrics    []models.MetricPoint `json:"metrics,omitempty"`
}

type SystemStatus struct {
	TotalNodes   int       `json:"total_nodes"`
	HealthyNodes int       `json:"healthy_nodes"`
	LastUpdate   time.Time `json:"last_update"`
}

type NodeHistory struct {
	NodeID    string
	Timestamp time.Time
	IsHealthy bool
	Services  []ServiceStatus
}

type NodeHistoryPoint struct {
	Timestamp time.Time `json:"timestamp"`
	IsHealthy bool      `json:"is_healthy"`
}

type APIServer struct {
	mu                 sync.RWMutex
	nodes              map[string]*NodeStatus
	router             *mux.Router
	nodeHistoryHandler func(nodeID string) ([]NodeHistoryPoint, error)
	metricsManager     metrics.MetricCollector
	snmpManager        snmp.SNMPManager
	db                 db.Service
	knownPollers       []string
}

func NewAPIServer(options ...func(server *APIServer)) *APIServer {
	s := &APIServer{
		nodes:  make(map[string]*NodeStatus),
		router: mux.NewRouter(),
	}

	for _, o := range options {
		o(s)
	}

	s.setupRoutes()

	return s
}

func WithMetricsManager(m metrics.MetricCollector) func(server *APIServer) {
	return func(server *APIServer) {
		server.metricsManager = m
	}
}

func WithSNMPManager(m snmp.SNMPManager) func(server *APIServer) {
	return func(server *APIServer) {
		server.snmpManager = m
	}
}

func WithDB(db db.Service) func(server *APIServer) {
	return func(server *APIServer) {
		server.db = db
	}
}

//go:embed web/dist/*
var webContent embed.FS

// setupStaticFileServing configures static file serving for the embedded web files.
func (s *APIServer) setupStaticFileServing() {
	fsys, err := fs.Sub(webContent, "web/dist")
	if err != nil {
		log.Printf("Error setting up static file serving: %v", err)
		return
	}
	fileServer := http.FileServer(http.FS(fsys))
	s.router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Static file handler for %s", r.URL.Path)
		if strings.HasPrefix(r.URL.Path, "/api/") || strings.HasPrefix(r.URL.Path, "/web-api/") {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		fileServer.ServeHTTP(w, r)
	})
}

// setProxyHeaders adds headers to the proxied request, including API key.
func (*APIServer) setProxyHeaders(proxyReq, originalReq *http.Request) {
	if apiKey := os.Getenv("API_KEY"); apiKey != "" {
		log.Printf("Attaching API key: %s", apiKey)

		proxyReq.Header.Set("X-API-Key", apiKey)
	} else {
		log.Println("API_KEY not set, skipping header attachment")
	}

	for name, values := range originalReq.Header {
		if !strings.EqualFold(name, "x-api-key") {
			for _, value := range values {
				proxyReq.Header.Add(name, value)
			}
		}
	}
}

// executeRequest performs the HTTP request.
func (*APIServer) executeRequest(req *http.Request) (*http.Response, error) {
	client := &http.Client{}

	return client.Do(req)
}

// copyResponse writes the response headers and body to the writer.
func (*APIServer) copyResponse(w http.ResponseWriter, resp *http.Response) error {
	for name, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(name, value)
		}
	}

	w.WriteHeader(resp.StatusCode)

	_, err := io.Copy(w, resp.Body)
	if err != nil {
		log.Printf("Failed to copy response body: %v", err)
	}

	return err
}

// writeError writes an HTTP error response and logs the issue.
func (*APIServer) writeError(w http.ResponseWriter, msg string, status int) {
	http.Error(w, msg, status)
	log.Printf("%s: %d", msg, status)
}

// generateCSRFToken creates a random token.
func generateCSRFToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func (s *APIServer) setupRoutes() {
	log.Println("Setting up routes...")

	// Apply CSRF middleware to all routes
	s.router.Use(s.csrfMiddleware)

	// Proxy routes
	s.setupWebProxyRoutes("localhost:8090")
	log.Println("Web proxy routes registered for /web-api/")

	// API routes with additional middleware
	apiRouter := s.router.PathPrefix("/api").Subrouter()
	middlewareChain := func(next http.Handler) http.Handler {
		return srHttp.CommonMiddleware(srHttp.APIKeyMiddleware(next))
	}
	apiRouter.Use(middlewareChain)
	apiRouter.HandleFunc("/nodes", s.getNodes).Methods("GET")
	apiRouter.HandleFunc("/status", s.getSystemStatus).Methods("GET")
	apiRouter.HandleFunc("/nodes/{id}", s.getNode).Methods("GET")
	apiRouter.HandleFunc("/nodes/{id}/history", s.getNodeHistory).Methods("GET")
	apiRouter.HandleFunc("/nodes/{id}/metrics", s.getNodeMetrics).Methods("GET")
	apiRouter.HandleFunc("/nodes/{id}/services", s.getNodeServices).Methods("GET")
	apiRouter.HandleFunc("/nodes/{id}/services/{service}", s.getServiceDetails).Methods("GET")
	apiRouter.HandleFunc("/nodes/{id}/snmp", s.getSNMPData).Methods("GET")
	log.Println("API routes registered with middleware under /api/")

	// Static file serving
	s.configureStaticServing()
	log.Println("Static file serving registered")
}

func (s *APIServer) setupWebProxyRoutes(listenAddr string) {
	webProxyRouter := mux.NewRouter()
	webProxyRouter.Use(srHttp.CommonMiddleware)
	proxyHandler := func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Proxy handler hit for %s", r.URL.Path)
		s.proxyAPIRequest(w, r, listenAddr)
	}
	webProxyRouter.PathPrefix("/").HandlerFunc(proxyHandler)
	s.router.PathPrefix("/web-api/").Handler(webProxyRouter)
}

// csrfMiddleware ensures CSRF token is set and validated
func (s *APIServer) csrfMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CSRF token on any GET request not under /web-api/ or /api/
		if r.Method == http.MethodGet && !strings.HasPrefix(r.URL.Path, "/web-api/") && !strings.HasPrefix(r.URL.Path, "/api/") {
			token, err := generateCSRFToken()
			if err != nil {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
			http.SetCookie(w, &http.Cookie{
				Name:     "csrf_token",
				Value:    token,
				Path:     "/",
				HttpOnly: false,
				SameSite: http.SameSiteLaxMode,
			})
			log.Printf("Set CSRF token: %s for %s", token, r.URL.Path)
		}

		// Validate CSRF token for /web-api/ requests
		if strings.HasPrefix(r.URL.Path, "/web-api/") {
			cookie, err := r.Cookie("csrf_token")
			if err != nil || cookie == nil || cookie.Value == "" {
				http.Error(w, "CSRF token missing", http.StatusForbidden)
				log.Printf("CSRF token missing for %s", r.URL.Path)
				return
			}
			headerToken := r.Header.Get("X-CSRF-Token")
			if headerToken == "" || headerToken != cookie.Value {
				http.Error(w, "Invalid CSRF token", http.StatusForbidden)
				log.Printf("Invalid CSRF token for %s", r.URL.Path)
				return
			}
			log.Printf("CSRF token validated for %s", r.URL.Path)
		}

		next.ServeHTTP(w, r)
	})
}

// proxyAPIRequest forwards an incoming request to an internal API server.
func (s *APIServer) proxyAPIRequest(w http.ResponseWriter, r *http.Request, serverAddr string) {
	apiPath := strings.TrimPrefix(r.URL.Path, "/web-api/")
	internalURL := fmt.Sprintf("http://%s/%s", serverAddr, apiPath)
	if r.URL.RawQuery != "" {
		internalURL += "?" + r.URL.RawQuery
	}
	log.Printf("Proxying %s to %s", r.URL.Path, internalURL)
	proxyReq, err := http.NewRequestWithContext(r.Context(), r.Method, internalURL, r.Body)
	if err != nil {
		s.writeError(w, "Failed to create proxy request", http.StatusInternalServerError)
		return
	}
	s.setProxyHeaders(proxyReq, r)
	resp, err := s.executeRequest(proxyReq)
	if err != nil {
		s.writeError(w, "Failed to execute request", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	s.copyResponse(w, resp)
}

// getSNMPData retrieves SNMP data for a specific node.
func (s *APIServer) getSNMPData(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	nodeID := vars["id"]

	// Get start and end times from query parameters
	startStr := r.URL.Query().Get("start")
	endStr := r.URL.Query().Get("end")

	if startStr == "" || endStr == "" {
		http.Error(w, "start and end parameters are required", http.StatusBadRequest)

		return
	}

	startTime, err := time.Parse(time.RFC3339, startStr)
	if err != nil {
		http.Error(w, "Invalid start time format", http.StatusBadRequest)

		return
	}

	endTime, err := time.Parse(time.RFC3339, endStr)
	if err != nil {
		http.Error(w, "Invalid end time format", http.StatusBadRequest)

		return
	}

	// Use the injected snmpManager to fetch SNMP metrics
	snmpMetrics, err := s.snmpManager.GetSNMPMetrics(nodeID, startTime, endTime)
	if err != nil {
		log.Printf("Error fetching SNMP data for node %s: %v", nodeID, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(snmpMetrics); err != nil {
		log.Printf("Error encoding SNMP data response for node %s: %v", nodeID, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func (s *APIServer) getNodeMetrics(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	nodeID := vars["id"]

	if s.metricsManager == nil {
		log.Printf("Metrics not configured for node %s", nodeID)
		http.Error(w, "Metrics not configured", http.StatusInternalServerError)

		return
	}

	m := s.metricsManager.GetMetrics(nodeID)
	if m == nil {
		log.Printf("No metrics found for node %s", nodeID)
		http.Error(w, "No metrics found", http.StatusNotFound)

		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(m); err != nil {
		log.Printf("Error encoding metrics response for node %s: %v", nodeID, err)

		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func (s *APIServer) SetNodeHistoryHandler(handler func(nodeID string) ([]NodeHistoryPoint, error)) {
	s.nodeHistoryHandler = handler
}

func (s *APIServer) UpdateNodeStatus(nodeID string, status *NodeStatus) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.nodes[nodeID] = status
}

func (s *APIServer) getNodeHistory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	nodeID := vars["id"]

	log.Printf("Getting node history for: %s", nodeID)

	if s.nodeHistoryHandler == nil {
		http.Error(w, "History handler not configured", http.StatusInternalServerError)
		return
	}

	points, err := s.nodeHistoryHandler(nodeID)
	if err != nil {
		log.Printf("Error fetching node history: %v", err)
		http.Error(w, "Failed to fetch history", http.StatusInternalServerError)

		return
	}

	log.Printf("Fetched %d history points for node: %s", len(points), nodeID)

	if err := s.encodeJSONResponse(w, points); err != nil {
		log.Printf("Error encoding history response: %v", err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}

func (s *APIServer) getSystemStatus(w http.ResponseWriter, _ *http.Request) {
	s.mu.RLock()
	status := SystemStatus{
		TotalNodes:   len(s.nodes),
		HealthyNodes: 0,
		LastUpdate:   time.Now(),
	}

	for _, node := range s.nodes {
		if node.IsHealthy {
			status.HealthyNodes++
		}
	}
	s.mu.RUnlock()

	log.Printf("System status: total=%d healthy=%d last_update=%s",
		status.TotalNodes, status.HealthyNodes, status.LastUpdate.Format(time.RFC3339))

	if err := s.encodeJSONResponse(w, status); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func (s *APIServer) getNodes(w http.ResponseWriter, _ *http.Request) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Preallocate the slice with the correct length
	nodes := make([]*NodeStatus, 0, len(s.nodes))

	// Append all map values to the slice
	for id, node := range s.nodes {
		// Only include known pollers
		for _, known := range s.knownPollers {
			if id == known {
				nodes = append(nodes, node)
				break
			}
		}
	}

	// Encode and send the response
	if err := s.encodeJSONResponse(w, nodes); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func (s *APIServer) SetKnownPollers(knownPollers []string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.knownPollers = knownPollers
}

func (s *APIServer) getNodeByID(nodeID string) (*NodeStatus, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	node, exists := s.nodes[nodeID]

	return node, exists
}

func (*APIServer) encodeJSONResponse(w http.ResponseWriter, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Error encoding JSON response: %v", err)

		return err
	}

	return nil
}

func (s *APIServer) getNode(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	nodeID := vars["id"]

	// Check if it's a known poller
	isKnown := false

	for _, known := range s.knownPollers {
		if nodeID == known {
			isKnown = true
			break
		}
	}

	if !isKnown {
		http.Error(w, "Node not found", http.StatusNotFound)
		return
	}

	node, exists := s.getNodeByID(nodeID)
	if !exists {
		log.Printf("Node %s not found", nodeID)
		http.Error(w, "Node not found", http.StatusNotFound)

		return
	}

	if s.metricsManager != nil {
		m := s.metricsManager.GetMetrics(nodeID)
		if m != nil {
			node.Metrics = m
			log.Printf("Attached %d metrics points to node %s response", len(m), nodeID)
		}
	}

	if err := s.encodeJSONResponse(w, node); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
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

	if err := s.encodeJSONResponse(w, node.Services); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
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
			if err := s.encodeJSONResponse(w, service); err != nil {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}

			return
		}
	}

	http.Error(w, "Service not found", http.StatusNotFound)
}

const (
	defaultReadTimeout  = 10 * time.Second
	defaultWriteTimeout = 10 * time.Second
	defaultIdleTimeout  = 60 * time.Second
)

func (s *APIServer) Start(addr string) error {
	srv := &http.Server{
		Addr:         addr,
		Handler:      s.router,
		ReadTimeout:  defaultReadTimeout,  // Timeout for reading the entire request, including the body.
		WriteTimeout: defaultWriteTimeout, // Timeout for writing the response.
		IdleTimeout:  defaultIdleTimeout,  // Timeout for idle connections waiting in the Keep-Alive state.
		// Optional: You can also set ReadHeaderTimeout to limit the time for reading request headers
		// ReadHeaderTimeout: 5 * time.Second,
	}

	return srv.ListenAndServe()
}
