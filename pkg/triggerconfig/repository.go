package triggerconfig

import (
	"context"
)

type RepositoryType string

const (
	MongoDBRepository RepositoryType = "MongoDBRepo"
)

// Defines the interface for CRUD operations
type Repository interface {
	Create(ctx context.Context, config TriggerConfig) (string, error)
	GetByName(ctx context.Context, name string) (*TriggerConfig, error)
	GetByID(ctx context.Context, id string) (*TriggerConfig, error)
	List(ctx context.Context) ([]*TriggerConfig, error)
	Update(ctx context.Context, id string, config TriggerConfig) error
	Delete(ctx context.Context, id string) error

	// Close closes the repository and all connections it manages
	Close(ctx context.Context) error
}
