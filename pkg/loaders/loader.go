package loaders

import (
	"github.com/Thibault-Van-Win/The-Instinct/pkg/action"
	"github.com/Thibault-Van-Win/The-Instinct/pkg/reflex"
)

// RuleLoader defines the interface for loading reflexes from various sources
type RuleLoader interface {
	LoadReflexes() ([]reflex.Reflex, error)
}

// ReflexConfig represents the structure of a reflex configuration
type ReflexConfig struct {
	Name       string                `yaml:"name" json:"name"`
	Expression string                `yaml:"expression" json:"expression"`
	Actions    []action.ActionConfig `yaml:"actions" json:"actions"`
}
