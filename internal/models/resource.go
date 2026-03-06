package models

import (
	"encoding/json"
	"time"
)

// Resource represents a monitored infrastructure resource (VM, pod, node, etc.)
type Resource struct {
	ID           int             `db:"id" json:"id"`
	ResourceID   string          `db:"resource_id" json:"resource_id"`
	ResourceType string          `db:"resource_type" json:"resource_type"`
	Name         string          `db:"name" json:"name"`
	Region       *string         `db:"region" json:"region,omitempty"`
	Tags         json.RawMessage `db:"tags" json:"tags,omitempty"`
	CreatedAt    time.Time       `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time       `db:"updated_at" json:"updated_at"`
}

// CreateResourceRequest represents the request body for creating a resource
type CreateResourceRequest struct {
	ResourceID   string          `json:"resource_id" binding:"required"`
	ResourceType string          `json:"resource_type" binding:"required"`
	Name         string          `json:"name" binding:"required"`
	Region       string          `json:"region,omitempty"`
	Tags         json.RawMessage `json:"tags,omitempty"`
}
