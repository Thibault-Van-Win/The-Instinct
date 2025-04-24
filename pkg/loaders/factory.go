package loaders

import (
	"fmt"

	"github.com/Thibault-Van-Win/The-Instinct/pkg/action"
)

// LoaderType represents the type of loader
type LoaderType string

const (
	YAMLLoaderType     LoaderType = "yaml"
	PostgresLoaderType LoaderType = "postgres"
)

// LoaderFactory creates rule loaders based on type
type LoaderFactory struct {
	ActionRegistry *action.ActionRegistry
}

// NewLoaderFactory creates a new loader factory
func NewLoaderFactory(registry *action.ActionRegistry) *LoaderFactory {
	return &LoaderFactory{
		ActionRegistry: registry,
	}
}

// CreateLoader creates a rule loader based on type and configuration
func (f *LoaderFactory) CreateLoader(loaderType LoaderType, config map[string]any) (RuleLoader, error) {
	switch loaderType {
	case YAMLLoaderType:
		directory, ok := config["directory"].(string)
		if !ok {
			return nil, fmt.Errorf("yaml loader requires a directory")
		}
		return NewYAMLFileLoader(directory, f.ActionRegistry), nil

	default:
		return nil, fmt.Errorf("unknown loader type: %s", loaderType)
	}
}
