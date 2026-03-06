package repository

import (
	"fmt"
	"time"

	"github.com/HarkHorning/infra-health-dashboard/internal/models"
	"github.com/jmoiron/sqlx"
)

// MetricRepository handles database operations for metrics
type MetricRepository struct {
	db *sqlx.DB
}

// NewMetricRepository creates a new MetricRepository
func NewMetricRepository(db *sqlx.DB) *MetricRepository {
	return &MetricRepository{db: db}
}

// Insert stores a new metric reading for a resource
func (r *MetricRepository) Insert(resourceID int, req models.CreateMetricRequest) (*models.Metric, error) {
	recordedAt := time.Now()
	if req.RecordedAt != "" {
		parsed, err := time.Parse(time.RFC3339, req.RecordedAt)
		if err != nil {
			return nil, fmt.Errorf("invalid recorded_at format (use RFC3339): %w", err)
		}
		recordedAt = parsed
	}

	query := `
		INSERT INTO metrics (resource_id, metric_name, metric_value, unit, recorded_at)
		VALUES (?, ?, ?, ?, ?)
	`

	var unit interface{}
	if req.Unit != "" {
		unit = req.Unit
	}

	result, err := r.db.Exec(query, resourceID, req.MetricName, req.MetricValue, unit, recordedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to insert metric: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert id: %w", err)
	}

	return r.GetByID(id)
}

// InsertBatch inserts multiple metrics in a single transaction
func (r *MetricRepository) InsertBatch(resourceID int, metrics []models.CreateMetricRequest) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := `
		INSERT INTO metrics (resource_id, metric_name, metric_value, unit, recorded_at)
		VALUES (?, ?, ?, ?, ?)
	`

	stmt, err := tx.Preparex(query)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	for _, req := range metrics {
		recordedAt := time.Now()
		if req.RecordedAt != "" {
			parsed, err := time.Parse(time.RFC3339, req.RecordedAt)
			if err != nil {
				return fmt.Errorf("invalid recorded_at format (use RFC3339): %w", err)
			}
			recordedAt = parsed
		}

		var unit interface{}
		if req.Unit != "" {
			unit = req.Unit
		}

		_, err := stmt.Exec(resourceID, req.MetricName, req.MetricValue, unit, recordedAt)
		if err != nil {
			return fmt.Errorf("failed to insert metric: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// GetByID fetches a single metric by its ID
func (r *MetricRepository) GetByID(id int64) (*models.Metric, error) {
	var metric models.Metric
	query := `SELECT id, resource_id, metric_name, metric_value, unit, recorded_at, created_at FROM metrics WHERE id = ?`

	err := r.db.Get(&metric, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get metric: %w", err)
	}

	return &metric, nil
}

// GetByResourceID fetches metrics for a resource within a time range
func (r *MetricRepository) GetByResourceID(resourceID int, query models.MetricQuery) ([]models.Metric, error) {
	var metrics []models.Metric
	var args []interface{}

	sql := `
		SELECT id, resource_id, metric_name, metric_value, unit, recorded_at, created_at
		FROM metrics
		WHERE resource_id = ?
	`
	args = append(args, resourceID)

	if query.MetricName != "" {
		sql += ` AND metric_name = ?`
		args = append(args, query.MetricName)
	}

	if !query.From.IsZero() {
		sql += ` AND recorded_at >= ?`
		args = append(args, query.From)
	}

	if !query.To.IsZero() {
		sql += ` AND recorded_at <= ?`
		args = append(args, query.To)
	}

	sql += ` ORDER BY recorded_at DESC`

	err := r.db.Select(&metrics, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get metrics: %w", err)
	}

	return metrics, nil
}

// GetLatestByResourceID fetches the most recent metric for each metric name
func (r *MetricRepository) GetLatestByResourceID(resourceID int) ([]models.Metric, error) {
	var metrics []models.Metric
	query := `
		SELECT m.id, m.resource_id, m.metric_name, m.metric_value, m.unit, m.recorded_at, m.created_at
		FROM metrics m
		INNER JOIN (
			SELECT metric_name, MAX(recorded_at) as max_recorded
			FROM metrics
			WHERE resource_id = ?
			GROUP BY metric_name
		) latest ON m.metric_name = latest.metric_name AND m.recorded_at = latest.max_recorded
		WHERE m.resource_id = ?
	`

	err := r.db.Select(&metrics, query, resourceID, resourceID)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest metrics: %w", err)
	}

	return metrics, nil
}

// DeleteOlderThan removes metrics older than the specified time (for cleanup)
func (r *MetricRepository) DeleteOlderThan(before time.Time) (int64, error) {
	query := `DELETE FROM metrics WHERE recorded_at < ?`

	result, err := r.db.Exec(query, before)
	if err != nil {
		return 0, fmt.Errorf("failed to delete old metrics: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get rows affected: %w", err)
	}

	return rowsAffected, nil
}
