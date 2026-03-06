package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

// RunMigrations creates all tables if they don't exist
func RunMigrations(db *sqlx.DB) error {
	migrations := []string{
		// Resources table
		`CREATE TABLE IF NOT EXISTS resources (
			id INT AUTO_INCREMENT PRIMARY KEY,
			resource_id VARCHAR(255) NOT NULL UNIQUE,
			resource_type VARCHAR(50) NOT NULL,
			name VARCHAR(255) NOT NULL,
			region VARCHAR(100),
			tags JSON,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
		)`,

		// Metrics table
		`CREATE TABLE IF NOT EXISTS metrics (
			id BIGINT AUTO_INCREMENT PRIMARY KEY,
			resource_id INT NOT NULL,
			metric_name VARCHAR(100) NOT NULL,
			metric_value DOUBLE NOT NULL,
			unit VARCHAR(50),
			recorded_at TIMESTAMP NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (resource_id) REFERENCES resources(id) ON DELETE CASCADE,
			INDEX idx_resource_metric_time (resource_id, metric_name, recorded_at)
		)`,

		// Health status table
		`CREATE TABLE IF NOT EXISTS health_status (
			id INT AUTO_INCREMENT PRIMARY KEY,
			resource_id INT NOT NULL,
			status VARCHAR(20) NOT NULL,
			reason TEXT,
			checked_at TIMESTAMP NOT NULL,
			FOREIGN KEY (resource_id) REFERENCES resources(id) ON DELETE CASCADE,
			INDEX idx_resource_checked (resource_id, checked_at)
		)`,

		// Alerts table
		`CREATE TABLE IF NOT EXISTS alerts (
			id INT AUTO_INCREMENT PRIMARY KEY,
			resource_id INT NOT NULL,
			alert_type VARCHAR(100) NOT NULL,
			severity VARCHAR(20) NOT NULL,
			message TEXT,
			triggered_at TIMESTAMP NOT NULL,
			resolved_at TIMESTAMP,
			FOREIGN KEY (resource_id) REFERENCES resources(id) ON DELETE CASCADE,
			INDEX idx_resource_alert (resource_id, triggered_at)
		)`,
	}

	for i, migration := range migrations {
		_, err := db.Exec(migration)
		if err != nil {
			return fmt.Errorf("migration %d failed: %w", i+1, err)
		}
	}

	return nil
}
