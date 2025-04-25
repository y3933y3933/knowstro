package store

import (
	"database/sql"
	"testing"

	_ "github.com/jackc/pgx/v4/stdlib"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("pgx", "host=localhost user=postgres password=postgres dbname=postgres port=5433")
	if err != nil {
		t.Fatalf("opening test db: %v", err)
	}

	err = Migrate(db, "../../migrations/")
	if err != nil {
		t.Fatalf("migrating test db error: %v", err)
	}

	_, err = db.Exec(`TRUNCATE resource_types CASCADE`)
	if err != nil {
		t.Fatalf("truncating tables %v", err)
	}

	return db
}

func TestCreateResourceType(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	store := NewPostgresResourceTypeStore(db)

	tests := []struct {
		name         string
		resourceType ResourceType
		wantErr      bool
	}{
		{
			name: "valid resourceType",
			resourceType: ResourceType{
				Name:        "course",
				Description: "線上/線下課程",
			},
			wantErr: false,
		},
		{
			name: "resource type with invalid name",
			resourceType: ResourceType{
				Name: "123456789012345678901234567890123456789012345678901234567890",
			},
			wantErr: true,
		},
		{
			name: "resource type with empty description",
			resourceType: ResourceType{
				Name: "book",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			createResourceType, err := store.CreateResourceType(&tt.resourceType)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.resourceType.Name, createResourceType.Name)
			assert.Equal(t, tt.resourceType.Description, createResourceType.Description)

			retrieved, err := store.GetResourceTypeByID(int64(createResourceType.ID))
			require.NoError(t, err)

			assert.Equal(t, createResourceType.ID, retrieved.ID)
		})
	}
}
