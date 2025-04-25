package loaders

import (
	"github.com/Thibault-Van-Win/The-Instinct/pkg/reflex"
)

// RuleLoader defines the interface for loading reflexes from various sources
type RuleLoader interface {
	LoadReflexes() ([]reflex.Reflex, error)
}
