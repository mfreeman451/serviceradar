package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/gorilla/websocket"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
)

type ConfigFile struct {
	NodeAddress string `json:"node_address"`
	Timeout     string `json:"timeout"`
	ListenAddr  string `json:"listen_addr"`
}

type Config struct {
	NodeAddress string
	Timeout     time.Duration
	ListenAddr  string
}

type DuskChecker struct {
	config    Config
	ws        *websocket.Conn
	sessionID string
	lastBlock time.Time
	mu        sync.RWMutex
	done      chan struct{}
}

// HealthServer implements the gRPC health check service
type HealthServer struct {
	grpc_health_v1.UnimplementedHealthServer
	checker *DuskChecker
}

func (s *HealthServer) Check(ctx context.Context, req *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	s.checker.mu.RLock()
	defer s.checker.mu.RUnlock()

	if s.checker.ws == nil {
		log.Printf("Health check failed: WebSocket connection not established")
		return &grpc_health_v1.HealthCheckResponse{
			Status: grpc_health_v1.HealthCheckResponse_NOT_SERVING,
		}, nil
	}

	if s.checker.lastBlock.IsZero() {
		// Connected but haven't received any blocks yet
		log.Printf("Health check warning: Connected but no blocks received yet")
		return &grpc_health_v1.HealthCheckResponse{
			Status: grpc_health_v1.HealthCheckResponse_NOT_SERVING,
		}, nil
	}

	timeSinceLastBlock := time.Since(s.checker.lastBlock)
	if timeSinceLastBlock > s.checker.config.Timeout {
		log.Printf("Health check failed: No blocks received in %v", timeSinceLastBlock)
		return &grpc_health_v1.HealthCheckResponse{
			Status: grpc_health_v1.HealthCheckResponse_NOT_SERVING,
		}, nil
	}

	// Everything is good
	log.Printf("Health check passed: Last block received %v ago", timeSinceLastBlock)
	return &grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	}, nil
}

func main() {
	log.Printf("Starting Dusk checker...")

	configFile := flag.String("config", "/etc/homemon/checkers/dusk.json", "Path to config file")
	flag.Parse()

	// Load config
	log.Printf("Loading config from: %s", *configFile)
	config, err := loadConfig(*configFile)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	log.Printf("Loaded config: %+v", config)

	// Create checker
	log.Printf("Creating Dusk checker...")
	checker := &DuskChecker{
		config: config,
		done:   make(chan struct{}),
	}

	// Start monitoring
	log.Printf("Starting monitoring...")
	if err := checker.startMonitoring(); err != nil {
		log.Fatalf("Failed to start monitoring: %v", err)
	}
	log.Printf("Monitoring started successfully")

	// Start gRPC server
	log.Printf("Starting gRPC server on %s", config.ListenAddr)
	lis, err := net.Listen("tcp", config.ListenAddr)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, &HealthServer{checker: checker})
	log.Printf("Registered health server")

	// Handle shutdown gracefully
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigCh
		log.Printf("Received signal %v, shutting down", sig)
		checker.Close()
		grpcServer.GracefulStop()
	}()

	log.Printf("Dusk checker listening on %s", config.ListenAddr)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func loadConfig(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("failed to read config file: %w", err)
	}

	var fileConfig ConfigFile
	if err := json.Unmarshal(data, &fileConfig); err != nil {
		return Config{}, fmt.Errorf("failed to parse config: %w", err)
	}

	// Parse the timeout string into a duration
	timeout, err := time.ParseDuration(fileConfig.Timeout)
	if err != nil {
		return Config{}, fmt.Errorf("invalid timeout value %q: %w", fileConfig.Timeout, err)
	}

	config := Config{
		NodeAddress: fileConfig.NodeAddress,
		Timeout:     timeout,
		ListenAddr:  fileConfig.ListenAddr,
	}

	// Set defaults if needed
	if config.Timeout == 0 {
		config.Timeout = 5 * time.Minute
	}
	if config.ListenAddr == "" {
		config.ListenAddr = "127.0.0.1:50052"
	}

	return config, nil
}

func (d *DuskChecker) startMonitoring() error {
	log.Printf("Connecting to Dusk node at %s", d.config.NodeAddress)

	// Connect to WebSocket
	u := url.URL{Scheme: "ws", Host: d.config.NodeAddress, Path: "/on"}
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return fmt.Errorf("websocket dial failed: %w", err)
	}
	d.ws = conn

	// Read initial session ID
	_, message, err := d.ws.ReadMessage()
	if err != nil {
		d.ws.Close()
		return fmt.Errorf("failed to read session ID: %w", err)
	}
	d.sessionID = string(message)
	log.Printf("Received session ID: %s", d.sessionID)

	// Subscribe to block events using session ID
	if err := d.subscribeToBlocks(); err != nil {
		d.ws.Close()
		return fmt.Errorf("failed to subscribe to blocks: %w", err)
	}

	go d.listenForEvents()
	return nil
}

func (d *DuskChecker) subscribeToBlocks() error {
	// Format subscription request with session ID
	subscribeReq := struct {
		Method  string `json:"method"`
		Path    string `json:"path"`
		Headers struct {
			RuskSessionID string `json:"Rusk-Session-Id"`
		} `json:"headers"`
	}{
		Method: "GET",
		Path:   "/on/blocks/accepted",
	}
	subscribeReq.Headers.RuskSessionID = d.sessionID

	log.Printf("Subscribing to blocks with request: %+v", subscribeReq)
	if err := d.ws.WriteJSON(subscribeReq); err != nil {
		return fmt.Errorf("failed to subscribe to blocks: %w", err)
	}
	log.Printf("Block subscription request sent")

	return nil
}

func (d *DuskChecker) listenForEvents() {
	log.Printf("Starting event listener")
	defer d.ws.Close()

	for {
		select {
		case <-d.done:
			return
		default:
			messageType, data, err := d.ws.ReadMessage()
			if err != nil {
				log.Printf("Error reading message: %v", err)
				return
			}

			// Log raw message for debugging
			log.Printf("Received message type %d: %s", messageType, string(data))

			// Parse block events - they might come in different formats
			// For now, just look for block acceptance indicators
			if strings.Contains(string(data), "/on/blocks/accepted") {
				log.Printf("Block acceptance event detected")
				d.mu.Lock()
				d.lastBlock = time.Now()
				d.mu.Unlock()
			}
		}
	}
}

func (d *DuskChecker) isHealthy() bool {
	d.mu.RLock()
	defer d.mu.RUnlock()

	if d.lastBlock.IsZero() {
		return false
	}

	return time.Since(d.lastBlock) <= d.config.Timeout
}

func (d *DuskChecker) Close() error {
	close(d.done)
	if d.ws != nil {
		return d.ws.Close()
	}
	return nil
}
