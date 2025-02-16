// Package models pkg/models/metrics.go
package models

import "time"

type MetricType string

const (
	MetricTypeNumeric MetricType = "numeric"
	MetricTypeString  MetricType = "string"
	MetricTypeBoolean MetricType = "boolean"
	MetricTypeGauge   MetricType = "gauge"
	MetricTypeCounter MetricType = "counter"
)

type MetricPoint struct {
	Timestamp   time.Time   `json:"timestamp"`
	Value       interface{} `json:"value"`
	ValueType   MetricType  `json:"value_type"`
	ServiceName string      `json:"service_name"`
}

type MetricsConfig struct {
	Enabled   bool `json:"metrics_enabled"`
	Retention int  `json:"metrics_retention"`
	MaxNodes  int  `json:"max_nodes"`
}

const MetricPointSize = 32 // 8 bytes timestamp + 8 bytes response + 16 bytes name
