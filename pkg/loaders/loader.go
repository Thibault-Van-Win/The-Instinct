package loaders

import (
	"context"

	"github.com/Thibault-Van-Win/The-Instinct/pkg/reflex"
)

// RuleLoader defines the interface for loading reflexes from various sources
type RuleLoader interface {
	ListReflexes(ctx context.Context) ([]*reflex.Reflex, error)
}
