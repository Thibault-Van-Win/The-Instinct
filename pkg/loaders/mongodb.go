package loaders

import (
	"context"
	"fmt"

	"github.com/Thibault-Van-Win/The-Instinct/internal/config"
	"github.com/Thibault-Van-Win/The-Instinct/pkg/action"
	"github.com/Thibault-Van-Win/The-Instinct/pkg/reflex"
	mongoRepo "github.com/Thibault-Van-Win/The-Instinct/pkg/reflex/mongo"
	"github.com/Thibault-Van-Win/The-Instinct/pkg/rule"
)

// MongoDBLoader loads reflexes from MongoDB
type MongoDBLoader struct {
	dbConfig *config.DatabaseConfig	
	RuleRegistry   *rule.RuleRegistry
	ActionRegistry *action.ActionRegistry
}

// NewMongoDBLoader creates a new MongoDB loader
func NewMongoDBLoader(dbConfig *config.DatabaseConfig, ruleRegistry *rule.RuleRegistry, actionRegistry *action.ActionRegistry) *MongoDBLoader {
	return &MongoDBLoader{
		dbConfig: dbConfig,
		RuleRegistry:   ruleRegistry,
		ActionRegistry: actionRegistry,
	}
}

// LoadReflexes implements the RuleLoader interface
func (l *MongoDBLoader) ListReflexes(ctx context.Context) ([]*reflex.Reflex, error) {

	repo, err := mongoRepo.NewRepository(l.dbConfig, l.RuleRegistry, l.ActionRegistry)
	if err != nil {
		return nil, fmt.Errorf("failed to create reflex repository: %v", err)
	}
	service := reflex.NewReflexService(repo)
	defer service.Close(ctx)

	return  service.ListReflexes(context.Background())
}
