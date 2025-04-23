package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type ResourceType struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type PostgresResourceTypeStore struct {
	db *sql.DB
}

func NewResourceTypeStore(db *sql.DB) *PostgresResourceTypeStore {
	return &PostgresResourceTypeStore{db: db}
}

type ResourceTypeStore interface {
	CreateResourceType(*ResourceType) (*ResourceType, error)
	GetResourceTypeByID(id int64) (*ResourceType, error)
	UpdateResourceType(*ResourceType) (*ResourceType, error)
}

func (pg *PostgresResourceTypeStore) CreateResourceType(resource_type *ResourceType) (*ResourceType, error) {
	query := `
	 INSERT INTO resource_types(name)
	 VALUES ($1)
	 RETURNING id, name
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := pg.db.QueryRowContext(ctx, query, resource_type.Name).Scan(&resource_type.ID, &resource_type.Name)

	if err != nil {
		return nil, err
	}
	return resource_type, nil
}

func (pg *PostgresResourceTypeStore) GetResourceTypeByID(id int64) (*ResourceType, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	var resourceType ResourceType

	query := `
		SELECT id, name 
		FROM resource_types
		WHERE id = ($1)
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := pg.db.QueryRowContext(ctx, query, id).Scan(&resourceType.ID, &resourceType.Name)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &resourceType, nil

}

func (pg *PostgresResourceTypeStore) UpdateResourceType(resourceType *ResourceType) (*ResourceType, error) {
	if resourceType.ID < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		UPDATE resource_types
		SET name = $2
		WHERE id = $1
		RETURNING id, name
	`

	args := []any{
		resourceType.ID, resourceType.Name,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := pg.db.QueryRowContext(ctx, query, args...).Scan(&resourceType.ID, &resourceType.Name)
	if err != nil {
		return nil, err
	}
	return resourceType, nil

}
