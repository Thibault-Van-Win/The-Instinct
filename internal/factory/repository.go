package factory

// Packets used for factories that sit between the db and domain model
// Placing these factories in the respective packages would create an import cycle

import (
	"errors"
	"fmt"

	"github.com/Thibault-Van-Win/The-Instinct/internal/config"
	"github.com/Thibault-Van-Win/The-Instinct/pkg/action"
	"github.com/Thibault-Van-Win/The-Instinct/pkg/reflex"
	mongoReflexRepo "github.com/Thibault-Van-Win/The-Instinct/pkg/reflex/mongo"
	"github.com/Thibault-Van-Win/The-Instinct/pkg/rule"
	"github.com/Thibault-Van-Win/The-Instinct/pkg/triggerconfig"
	mongoTCRepo "github.com/Thibault-Van-Win/The-Instinct/pkg/triggerconfig/mongo"
)

func NewReflexRepository(dbConfig *config.DatabaseConfig, ruleReg *rule.RuleRegistry, actionReg *action.ActionRegistry) (reflex.Repository, error) {
	if dbConfig == nil {
		return nil, errors.New("dbConfig cannot be nil")
	}

	switch dbConfig.Type {
	case config.MongoDB:
		return mongoReflexRepo.NewRepository(dbConfig, ruleReg, actionReg)
	default:
		return nil, unsupportedDatabaseError("reflexes", string(dbConfig.Type))
	}
}

func NewTriggerConfigRepository(dbConfig *config.DatabaseConfig) (triggerconfig.Repository, error) {
	if dbConfig == nil {
		return nil, errors.New("dbConfig cannot be nil")
	}

	switch dbConfig.Type {
	case config.MongoDB:
		return mongoTCRepo.NewRepository(dbConfig)
	default:
		return nil, unsupportedDatabaseError("trigger configs", string(dbConfig.Type))
	}
}

func unsupportedDatabaseError(model, dbType string) error {
	return fmt.Errorf("unsupported database type for %s: %s", model, dbType)
}
