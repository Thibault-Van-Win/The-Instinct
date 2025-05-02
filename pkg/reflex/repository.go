package reflex

import (
	"context"
)

type RepositoryType string

var (
	MongoDBRepository RepositoryType = "MongoDBRepo"
)

// Defines the interface for CRUD operations
type Repository interface {
	Create(ctx context.Context, config ReflexConfig) (string, error)
	GetByName(ctx context.Context, name string) (*Reflex, error)
	GetByID(ctx context.Context, id string) (*Reflex, error)
	List(ctx context.Context) ([]*Reflex, error)
	Update(ctx context.Context, id string, config ReflexConfig) error
	Delete(ctx context.Context, id string) error

	// Close closes the repository and all connections it manages
	Close(ctx context.Context) error
}
