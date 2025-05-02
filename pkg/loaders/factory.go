package loaders

import (
	"fmt"

	"github.com/Thibault-Van-Win/The-Instinct/internal/config"
	"github.com/Thibault-Van-Win/The-Instinct/pkg/action"
	"github.com/Thibault-Van-Win/The-Instinct/pkg/rule"
)

// LoaderType represents the type of loader
type LoaderType string

const (
	YAML    LoaderType = "yaml"
	MongoDB LoaderType = "mongodb"
)

// LoaderFactory creates rule loaders based on type
type LoaderFactory struct {
	RuleRegistry   *rule.RuleRegistry
	ActionRegistry *action.ActionRegistry
}

// NewLoaderFactory creates a new loader factory
func NewLoaderFactory(ruleRegistry *rule.RuleRegistry, actionRegistry *action.ActionRegistry) *LoaderFactory {
	return &LoaderFactory{
		RuleRegistry:   ruleRegistry,
		ActionRegistry: actionRegistry,
	}
}

// CreateLoader creates a rule loader based on type and configuration
func (f *LoaderFactory) CreateLoader(loaderType LoaderType, conf any) (RuleLoader, error) {
	switch loaderType {
	case YAML:
		yamlConf, ok := conf.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("expected a config map")
		}
		directory, ok := yamlConf["directory"].(string)
		if !ok {
			return nil, fmt.Errorf("yaml loader requires a directory")
		}
		return NewYAMLFileLoader(directory, f.RuleRegistry, f.ActionRegistry), nil
	case MongoDB:
		dbConfig, ok := conf.(*config.DatabaseConfig)
		if !ok {
			return nil, fmt.Errorf("expected a pointer to a database config")
		}
		return NewMongoDBLoader(dbConfig, f.RuleRegistry, f.ActionRegistry), nil
	default:
		return nil, fmt.Errorf("unknown loader type: %s", loaderType)
	}
}
