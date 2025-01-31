// Package agent pkg/agent/icmp_checker.go
package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/mfreeman451/serviceradar/pkg/models"
	"github.com/mfreeman451/serviceradar/pkg/scan"
)

type ICMPChecker struct {
	Host  string
	Count int
}

type ICMPResponse struct {
	Host         string  `json:"host"`
	ResponseTime int64   `json:"response_time"` // in nanoseconds
	PacketLoss   float64 `json:"packet_loss"`
	Available    bool    `json:"available"`
}

const (
	combinedScannerTimeout = 5 * time.Second
	combinedScannerConc    = 10
	combinedScannerICMP    = 1
	combinedScannerMaxIdle = 10
	combinedScannerMaxLife = 10 * time.Minute
	combinedScannerIdle    = 5 * time.Minute
)

func (p *ICMPChecker) Check(ctx context.Context) (available bool, response string) {
	scanner := scan.NewCombinedScanner(
		combinedScannerTimeout,
		combinedScannerConc,
		combinedScannerICMP,
		combinedScannerMaxIdle,
		combinedScannerMaxLife,
		combinedScannerIdle,
	)

	target := models.Target{
		Host: p.Host,
		Mode: models.ModeICMP,
	}

	resultChan, err := scanner.Scan(ctx, []models.Target{target})
	if err != nil {
		return false, fmt.Sprintf(`{"error": "%v"}`, err)
	}

	// Take first result as we only sent one target
	var result models.Result
	for r := range resultChan {
		result = r
		break
	}

	// Format response
	responseStruct := struct {
		Host         string  `json:"host"`
		ResponseTime int64   `json:"response_time"` // in nanoseconds
		PacketLoss   float64 `json:"packet_loss"`
		Available    bool    `json:"available"`
	}{
		Host:         p.Host,
		ResponseTime: result.RespTime.Nanoseconds(),
		PacketLoss:   result.PacketLoss,
		Available:    result.Available,
	}

	jsonResponse, err := json.Marshal(responseStruct)
	if err != nil {
		return false, fmt.Sprintf(`{"error": "Failed to marshal response: %v"}`, err)
	}

	return result.Available, string(jsonResponse)
}
