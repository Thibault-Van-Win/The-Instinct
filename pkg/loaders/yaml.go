package loaders

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"github.com/Thibault-Van-Win/The-Instinct/pkg/action"
	"github.com/Thibault-Van-Win/The-Instinct/pkg/reflex"
	"github.com/Thibault-Van-Win/The-Instinct/pkg/rule"
)

// YAMLFileLoader loads reflexes from YAML files in a directory
type YAMLFileLoader struct {
	Directory      string
	RuleRegistry *rule.RuleRegistry
	ActionRegistry *action.ActionRegistry
}

// NewYAMLFileLoader creates a new YAML file loader
func NewYAMLFileLoader(directory string, ruleRegistry *rule.RuleRegistry, actionRegistry *action.ActionRegistry) *YAMLFileLoader {
	return &YAMLFileLoader{
		Directory:      directory,
		RuleRegistry: ruleRegistry,
		ActionRegistry: actionRegistry,
	}
}

// LoadReflexes implements the RuleLoader interface
func (l *YAMLFileLoader) LoadReflexes() ([]reflex.Reflex, error) {
	var reflexes []reflex.Reflex

	// Read all files in the directory
	files, err := os.ReadDir(l.Directory)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	// Process each YAML file
	for _, file := range files {
		if file.IsDir() || (!(filepath.Ext(file.Name()) == ".yaml") && !(filepath.Ext(file.Name()) == ".yml")) {
			continue
		}

		// Read the file
		filePath := filepath.Join(l.Directory, file.Name())
		data, err := os.ReadFile(filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to read file %s: %w", filePath, err)
		}

		// Parse the YAML
		var configs []ReflexConfig
		if err := yaml.Unmarshal(data, &configs); err != nil {
			return nil, fmt.Errorf("failed to parse YAML in file %s: %w", filePath, err)
		}

		// Create reflexes from the configs
		for _, config := range configs {
			reflex, err := l.createReflex(config)
			if err != nil {
				return nil, fmt.Errorf("failed to create reflex from file %s: %w", filePath, err)
			}
			reflexes = append(reflexes, *reflex)
		}
	}

	return reflexes, nil
}

// createReflex creates a reflex from a config
func (l *YAMLFileLoader) createReflex(config ReflexConfig) (*reflex.Reflex, error) {
	// Create the rule
	ruleInstance, err := l.RuleRegistry.Create(config.RuleConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create rule: %v", err)
	}

	// Create the actions
	actions := make([]action.Action, 0, len(config.Actions))
	for _, actionConfig := range config.Actions {
		actionInstance, err := l.ActionRegistry.Create(actionConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to create action for reflex %s: %w", config.Name, err)
		}
		actions = append(actions, actionInstance)
	}

	// Create the reflex
	return reflex.NewReflex(config.Name, ruleInstance, actions), nil
}
