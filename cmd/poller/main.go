// cmd/poller/main.go
package main

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mfreeman451/homemon/pkg/poller"
)

func main() {
	config := poller.Config{
		Agents: map[string]string{
			"local-agent": "localhost:50051",
		},
		CloudAddress: "cloud-service:50052",
		PollInterval: 30 * time.Second,
		PollerID:     "home-poller-1",
	}

	p, err := poller.New(config)
	if err != nil {
		log.Fatalf("Failed to create poller: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigCh
		cancel()
	}()

	if err := p.Start(ctx); err != nil && !errors.Is(err, context.Canceled) {
		log.Fatalf("Poller failed: %v", err)
	}
}
