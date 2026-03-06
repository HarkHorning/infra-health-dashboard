package repository

import (
	"database/sql"
	"fmt"

	"github.com/HarkHorning/infra-health-dashboard/internal/models"
	"github.com/jmoiron/sqlx"
)

// ResourceRepository handles database operations for resources
type ResourceRepository struct {
	db *sqlx.DB
}

// NewResourceRepository creates a new ResourceRepository
func NewResourceRepository(db *sqlx.DB) *ResourceRepository {
	return &ResourceRepository{db: db}
}

// Create inserts a new resource into the database
func (r *ResourceRepository) Create(req models.CreateResourceRequest) (*models.Resource, error) {
	query := `
		INSERT INTO resources (resource_id, resource_type, name, region, tags)
		VALUES (?, ?, ?, ?, ?)
	`

	var region interface{}
	if req.Region != "" {
		region = req.Region
	}

	var tags interface{}
	if len(req.Tags) > 0 {
		tags = req.Tags
	}

	result, err := r.db.Exec(query, req.ResourceID, req.ResourceType, req.Name, region, tags)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert id: %w", err)
	}

	return r.GetByID(int(id))
}

// GetByID fetches a single resource by its ID
func (r *ResourceRepository) GetByID(id int) (*models.Resource, error) {
	var resource models.Resource
	query := `SELECT id, resource_id, resource_type, name, region, COALESCE(tags, '{}') as tags, created_at, updated_at FROM resources WHERE id = ?`

	err := r.db.Get(&resource, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("resource not found: %d", id)
		}
		return nil, fmt.Errorf("failed to get resource: %w", err)
	}

	return &resource, nil
}

// GetByResourceID fetches a resource by its external resource ID (Azure ID or K8s UID)
func (r *ResourceRepository) GetByResourceID(resourceID string) (*models.Resource, error) {
	var resource models.Resource
	query := `SELECT id, resource_id, resource_type, name, region, COALESCE(tags, '{}') as tags, created_at, updated_at FROM resources WHERE resource_id = ?`

	err := r.db.Get(&resource, query, resourceID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("resource not found: %s", resourceID)
		}
		return nil, fmt.Errorf("failed to get resource: %w", err)
	}

	return &resource, nil
}

// List fetches all resources from the database
func (r *ResourceRepository) List() ([]models.Resource, error) {
	var resources []models.Resource
	query := `SELECT id, resource_id, resource_type, name, region, COALESCE(tags, '{}') as tags, created_at, updated_at FROM resources ORDER BY created_at DESC`

	err := r.db.Select(&resources, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list resources: %w", err)
	}

	return resources, nil
}

// ListByType fetches all resources of a specific type
func (r *ResourceRepository) ListByType(resourceType string) ([]models.Resource, error) {
	var resources []models.Resource
	query := `SELECT id, resource_id, resource_type, name, region, COALESCE(tags, '{}') as tags, created_at, updated_at FROM resources WHERE resource_type = ? ORDER BY created_at DESC`

	err := r.db.Select(&resources, query, resourceType)
	if err != nil {
		return nil, fmt.Errorf("failed to list resources by type: %w", err)
	}

	return resources, nil
}

// Delete removes a resource from the database
func (r *ResourceRepository) Delete(id int) error {
	query := `DELETE FROM resources WHERE id = ?`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete resource: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("resource not found: %d", id)
	}

	return nil
}

// Update updates an existing resource
func (r *ResourceRepository) Update(id int, req models.CreateResourceRequest) (*models.Resource, error) {
	query := `
		UPDATE resources
		SET resource_id = ?, resource_type = ?, name = ?, region = ?, tags = ?
		WHERE id = ?
	`

	var region interface{}
	if req.Region != "" {
		region = req.Region
	}

	var tags interface{}
	if len(req.Tags) > 0 {
		tags = req.Tags
	}

	result, err := r.db.Exec(query, req.ResourceID, req.ResourceType, req.Name, region, tags, id)
	if err != nil {
		return nil, fmt.Errorf("failed to update resource: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return nil, fmt.Errorf("resource not found: %d", id)
	}

	return r.GetByID(id)
}
