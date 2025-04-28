package loaders

import (
	"fmt"

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
func (f *LoaderFactory) CreateLoader(loaderType LoaderType, config map[string]any) (RuleLoader, error) {
	switch loaderType {
	case YAML:
		directory, ok := config["directory"].(string)
		if !ok {
			return nil, fmt.Errorf("yaml loader requires a directory")
		}
		return NewYAMLFileLoader(directory, f.RuleRegistry, f.ActionRegistry), nil
	case MongoDB:
		uri, ok := config["uri"].(string)
		if !ok {
			return nil, fmt.Errorf("mongoDB loader requires a uri")
		}

		database, ok := config["database"].(string)
		if !ok {
			return nil, fmt.Errorf("mongoDB loader requires a database name")
		}

		collection, ok := config["collection"].(string)
		if !ok {
			return nil, fmt.Errorf("mongoDB loader requires a collection name")
		}

		return NewMongoDBLoader(uri, database, collection, f.RuleRegistry, f.ActionRegistry), nil
	default:
		return nil, fmt.Errorf("unknown loader type: %s", loaderType)
	}
}
