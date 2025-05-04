package reflex

import (
	"errors"
	"fmt"

	"github.com/Thibault-Van-Win/The-Instinct/pkg/action"
	"github.com/Thibault-Van-Win/The-Instinct/pkg/rule"
)

// Domain model for a reflex configuration
type ReflexConfig struct {
	Name         string              `yaml:"name" json:"name"`
	RuleConfig   rule.RuleConfig     `yaml:"rule" json:"rule"`
	ActionConfig action.ActionConfig `yaml:"action" json:"action"`
}

func (rc *ReflexConfig) Validate() error {
	if rc.Name == "" {
		return errors.New("reflex name cannot be empty")
	}

	// Validate rule configuration
	if err := rc.RuleConfig.Validate(); err != nil {
		return fmt.Errorf("failed to validate rule config: %v", err)
	}

	// Validate action configurations
	if err := rc.ActionConfig.Validate(); err != nil {
		return fmt.Errorf("failed to validate action config: %v", err)
	}

	return nil
}
