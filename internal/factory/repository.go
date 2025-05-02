package factory

// Package used for a the Reflex Repository
// This could not be placed in pkg/reflex/repository as this would create an import cycle.
// Might move all factories over here eventually

import (
	"fmt"

	"github.com/Thibault-Van-Win/The-Instinct/internal/config"
	"github.com/Thibault-Van-Win/The-Instinct/pkg/action"
	"github.com/Thibault-Van-Win/The-Instinct/pkg/reflex"
	mongoRepo "github.com/Thibault-Van-Win/The-Instinct/pkg/reflex/mongo"
	"github.com/Thibault-Van-Win/The-Instinct/pkg/rule"
)

func NewReflexRepository(dbConfig *config.DatabaseConfig, ruleReg *rule.RuleRegistry, actionReg *action.ActionRegistry) (reflex.Repository, error) {
	switch dbConfig.Type {
	case config.MongoDB:
		return mongoRepo.NewRepository(dbConfig, ruleReg, actionReg)
	default:
		return nil, fmt.Errorf("unsupported database type for reflexes: %s", dbConfig.Type)
	}
}
