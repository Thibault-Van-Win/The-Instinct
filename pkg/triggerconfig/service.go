package triggerconfig

import (
	"context"
	"errors"
	"fmt"
)

type TriggerConfigService struct {
	repo Repository
}

func NewTriggerConfigService(repo Repository) *TriggerConfigService {
	return &TriggerConfigService{
		repo: repo,
	}
}

// Close the repository and all the connections it manages
func (s *TriggerConfigService) Close(ctx context.Context) error {
	return s.repo.Close(ctx)
}

func (s *TriggerConfigService) CreateTriggerConfig(ctx context.Context, config TriggerConfig) (string, error) {
	if err := config.Validate(); err != nil {
		return "", fmt.Errorf("trigger config validation failed: %v", err)
	}

	return s.repo.Create(ctx, config)
}

// Retrieve a trigger config by name
func (s *TriggerConfigService) GetTriggerConfigByName(ctx context.Context, name string) (*TriggerConfig, error) {
	if name == "" {
		return nil, errors.New("trigger config name cannot be empty")
	}

	return s.repo.GetByName(ctx, name)
}

// Retrieve a trigger config by ID
func (s *TriggerConfigService) GetTriggerConfigByID(ctx context.Context, id string) (*TriggerConfig, error) {
	if id == "" {
		return nil, errors.New("trigger config ID cannot be empty")
	}

	return s.repo.GetByID(ctx, id)
}

// Retrieve all trigger configs from the db
func (s *TriggerConfigService) ListTriggerConfigs(ctx context.Context) ([]*TriggerConfig, error) {
	return s.repo.List(ctx)
}

// UpdateTriggerConfig updates an existing trigger config
func (s *TriggerConfigService) UpdateTriggerConfig(ctx context.Context, id string, config TriggerConfig) error {
	if id == "" {
		return errors.New("trigger config ID cannot be empty")
	}

	if err := config.Validate(); err != nil {
		return fmt.Errorf("failed to validate trigger config: %v", err)
	}

	return s.repo.Update(ctx, id, config)
}

// DeleteTriggerConfig deletes a trigger config by its ID
func (s *TriggerConfigService) DeleteTriggerConfig(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("reflex ID cannot be empty")
	}

	return s.repo.Delete(ctx, id)
}
