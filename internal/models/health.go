package models

import "time"

// HealthStatus represents a health status snapshot for a resource
type HealthStatus struct {
	ID         int       `db:"id" json:"id"`
	ResourceID int       `db:"resource_id" json:"resource_id"`
	Status     string    `db:"status" json:"status"`
	Reason     *string   `db:"reason" json:"reason,omitempty"`
	CheckedAt  time.Time `db:"checked_at" json:"checked_at"`
}

// CreateHealthStatusRequest represents the request body for creating a health status
type CreateHealthStatusRequest struct {
	Status    string `json:"status" binding:"required"`
	Reason    string `json:"reason,omitempty"`
	CheckedAt string `json:"checked_at,omitempty"`
}

// Valid health status values
const (
	HealthStatusHealthy   = "healthy"
	HealthStatusUnhealthy = "unhealthy"
	HealthStatusDegraded  = "degraded"
	HealthStatusUnknown   = "unknown"
)

// IsValidHealthStatus checks if a status value is valid
func IsValidHealthStatus(status string) bool {
	switch status {
	case HealthStatusHealthy, HealthStatusUnhealthy, HealthStatusDegraded, HealthStatusUnknown:
		return true
	}
	return false
}
