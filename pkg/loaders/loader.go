package loaders

import (
	"github.com/Thibault-Van-Win/The-Instinct/pkg/action"
	"github.com/Thibault-Van-Win/The-Instinct/pkg/reflex"
	"github.com/Thibault-Van-Win/The-Instinct/pkg/rule"
)

// RuleLoader defines the interface for loading reflexes from various sources
type RuleLoader interface {
	LoadReflexes() ([]reflex.Reflex, error)
}

// ReflexConfig represents the structure of a reflex configuration
type ReflexConfig struct {
	Name       string                `yaml:"name" json:"name"`
	RuleConfig rule.RuleConfig       `yaml:"rule" json:"rule"`
	Actions    []action.ActionConfig `yaml:"actions" json:"actions"`
}
