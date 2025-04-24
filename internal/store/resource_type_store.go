package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type ResourceType struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
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
	DeleteResourceType(id int64) error
	GetAllResourceType() ([]*ResourceType, error)
}

func (pg *PostgresResourceTypeStore) CreateResourceType(resource_type *ResourceType) (*ResourceType, error) {
	query := `
	 INSERT INTO resource_types(name, description)
	 VALUES ($1, $2)
	 RETURNING id, name, description
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{
		resource_type.Name,
		resource_type.Description,
	}

	err := pg.db.QueryRowContext(ctx, query, args...).Scan(&resource_type.ID, &resource_type.Name, &resource_type.Description)

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
		SELECT id, name , description
		FROM resource_types
		WHERE id = ($1)
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := pg.db.QueryRowContext(ctx, query, id).Scan(&resourceType.ID, &resourceType.Name, &resourceType.Description)
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
		SET name = $1, description = $2
		WHERE id = $3
		RETURNING id, name, description
	`

	args := []any{
		resourceType.Name, resourceType.Description, resourceType.ID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := pg.db.QueryRowContext(ctx, query, args...).Scan(&resourceType.ID, &resourceType.Name, &resourceType.Description)
	if err != nil {
		return nil, err
	}
	return resourceType, nil

}

func (pg *PostgresResourceTypeStore) DeleteResourceType(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
		DELETE FROM resource_types
		WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := pg.db.ExecContext(ctx, query, id)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}

func (pg *PostgresResourceTypeStore) GetAllResourceType() ([]*ResourceType, error) {
	query := `
		SELECT id, name, description
		FROM resource_types
		ORDER BY id
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := pg.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	resourceTypes := []*ResourceType{}

	for rows.Next() {
		var resourceType ResourceType
		err := rows.Scan(
			&resourceType.ID,
			&resourceType.Name,
			&resourceType.Description,
		)
		if err != nil {
			return nil, err
		}

		resourceTypes = append(resourceTypes, &resourceType)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return resourceTypes, nil
}
