// cmd/agent/main.go
package main

import (
	"encoding/json"
	"log"
	"net"
	"os"
	"path/filepath"

	"github.com/mfreeman451/homemon/pkg/agent"
	"github.com/mfreeman451/homemon/pkg/checker"
	"github.com/mfreeman451/homemon/proto"
	"google.golang.org/grpc"
)

type ExternalCheckerConfig struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

func loadExternalCheckers(configDir string) (map[string]checker.Checker, error) {
	checkers := make(map[string]checker.Checker)

	// Read all .json files in the checkers config directory
	files, err := filepath.Glob(filepath.Join(configDir, "*.json"))
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		var config ExternalCheckerConfig
		data, err := os.ReadFile(file)
		if err != nil {
			log.Printf("Warning: Failed to read checker config %s: %v", file, err)
			continue
		}

		if err := json.Unmarshal(data, &config); err != nil {
			log.Printf("Warning: Failed to parse checker config %s: %v", file, err)
			continue
		}

		c, err := agent.NewExternalChecker(config.Name, config.Address)
		if err != nil {
			log.Printf("Warning: Failed to create checker from %s: %v", file, err)
			continue
		}

		checkers[config.Name] = c
	}

	return checkers, nil
}

func main() {
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
		// Add external checkers to the map
		for name, c := range externalCheckers {
			checkers[name] = c
		}
	}

	server := grpc.NewServer()
	proto.RegisterAgentServiceServer(server, agent.NewServer(checkers))

	log.Printf("Agent server listening on :50051")
	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
