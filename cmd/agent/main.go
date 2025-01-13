// cmd/agent/main.go
package main

import (
	"encoding/json"
	"log"
	"net"
	"os"

	"github.com/mfreeman451/homemon/pkg/agent"
	"github.com/mfreeman451/homemon/pkg/checker"
	"github.com/mfreeman451/homemon/proto"
	"google.golang.org/grpc"
)

// ExternalCheckerConfig represents the configuration for an external checker
type ExternalCheckerConfig struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

func loadExternalCheckers(configDir string) (map[string]checker.Checker, error) {
	checkers := make(map[string]checker.Checker)

	// Specifically look for external.json
	configPath := configDir + "/external.json"

	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			// It's okay if the file doesn't exist - just return empty map
			return checkers, nil
		}
		return nil, err
	}

	var config ExternalCheckerConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	// Create the external checker
	if config.Name != "" && config.Address != "" {
		log.Printf("Creating external checker %s at address %s", config.Name, config.Address)
		checker, err := agent.NewExternalChecker(config.Name, config.Address)
		if err != nil {
			return nil, err
		}
		checkers[config.Name] = checker
	}

	return checkers, nil
}

func main() {
	log.Printf("Starting homemon agent...")
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Initialize built-in checkers
	checkers := map[string]checker.Checker{
		"nginx": &agent.ProcessChecker{ProcessName: "nginx"},
		"ssh":   &agent.PortChecker{Host: "localhost", Port: 22},
	}

	// Load external checkers
	externalCheckers, err := loadExternalCheckers("/etc/homemon/checkers")
	if err != nil {
		log.Printf("Warning: Failed to load external checkers: %v", err)
	} else {
		for name, checker := range externalCheckers {
			log.Printf("Adding external checker: %s", name)
			checkers[name] = checker
		}
	}

	server := grpc.NewServer()
	proto.RegisterAgentServiceServer(server, agent.NewServer(checkers))

	log.Printf("Agent server listening on :50051")
	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
