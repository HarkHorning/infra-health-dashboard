package models

import "time"

// Metric represents a time-series metric data point
type Metric struct {
	ID          int64     `db:"id" json:"id"`
	ResourceID  int       `db:"resource_id" json:"resource_id"`
	MetricName  string    `db:"metric_name" json:"metric_name"`
	MetricValue float64   `db:"metric_value" json:"metric_value"`
	Unit        *string   `db:"unit" json:"unit,omitempty"`
	RecordedAt  time.Time `db:"recorded_at" json:"recorded_at"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
}

// CreateMetricRequest represents the request body for inserting a metric
type CreateMetricRequest struct {
	MetricName  string  `json:"metric_name" binding:"required"`
	MetricValue float64 `json:"metric_value"`
	Unit        string  `json:"unit,omitempty"`
	RecordedAt  string  `json:"recorded_at,omitempty"` // Optional, defaults to now
}

// MetricQuery represents query parameters for fetching metrics
type MetricQuery struct {
	MetricName string
	From       time.Time
	To         time.Time
}
