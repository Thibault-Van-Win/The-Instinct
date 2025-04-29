package reflex

import (
	"context"
	"errors"
)

type ReflexService struct {
	repo Repository
}

func NewReflexService(repo Repository) *ReflexService {
	return &ReflexService{
		repo: repo,
	}
}

func (s *ReflexService) CreateReflex(ctx context.Context, config ReflexConfig) (string, error) {
	// Validate the configuration
	if config.Name == "" {
		return "", errors.New("reflex name cannot be empty")
	}

	// Validate rule configuration
	// TODO

	// Validate action configurations
	// TODO

	return s.repo.Create(ctx, config)
}

// Retrieve a reflex by name
func (s *ReflexService) GetReflexByName(ctx context.Context, name string) (*Reflex, error) {
	if name == "" {
		return nil, errors.New("reflex name cannot be empty")
	}

	return s.repo.GetByName(ctx, name)
}

// Retrieve a reflex by ID
func (s *ReflexService) GetReflexByID(ctx context.Context, id string) (*Reflex, error) {
	if id == "" {
		return nil, errors.New("reflex ID cannot be empty")
	}

	return s.repo.GetByID(ctx, id)
}

// Retrieve all reflexes from the db
func (s *ReflexService) ListReflexes(ctx context.Context) ([]*Reflex, error) {
	return s.repo.List(ctx)
}

// UpdateReflex updates an existing reflex
func (s *ReflexService) UpdateReflex(ctx context.Context, id string, config ReflexConfig) error {
	if id == "" {
		return errors.New("reflex ID cannot be empty")
	}
	
	// Validate the configuration
	if config.Name == "" {
		return errors.New("reflex name cannot be empty")
	}
	
	// Validate rule configuration
	// TODO
	
	// Validate action configurations
	// TODO	
	
	return s.repo.Update(ctx, id, config)
}

// DeleteReflex deletes a reflex by its ID
func (s *ReflexService) DeleteReflex(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("reflex ID cannot be empty")
	}
	
	return s.repo.Delete(ctx, id)
}
