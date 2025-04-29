package loaders

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/Thibault-Van-Win/The-Instinct/pkg/action"
	"github.com/Thibault-Van-Win/The-Instinct/pkg/reflex"
	mongoRepo "github.com/Thibault-Van-Win/The-Instinct/pkg/reflex/mongo"
	"github.com/Thibault-Van-Win/The-Instinct/pkg/rule"
)

// MongoDBLoader loads reflexes from MongoDB
type MongoDBLoader struct {
	URI            string
	DatabaseName   string
	CollectionName string
	RuleRegistry   *rule.RuleRegistry
	ActionRegistry *action.ActionRegistry
}

// NewMongoDBLoader creates a new MongoDB loader
func NewMongoDBLoader(uri, dbName, collName string, ruleRegistry *rule.RuleRegistry, actionRegistry *action.ActionRegistry) *MongoDBLoader {
	return &MongoDBLoader{
		URI:            uri,
		DatabaseName:   dbName,
		CollectionName: collName,
		RuleRegistry:   ruleRegistry,
		ActionRegistry: actionRegistry,
	}
}

// LoadReflexes implements the RuleLoader interface
func (l *MongoDBLoader) LoadReflexes() ([]reflex.Reflex, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(l.URI))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}
	defer client.Disconnect(ctx)

	repo := mongoRepo.NewRepository(client.Database("instinct"), l.RuleRegistry, l.ActionRegistry)
	service := reflex.NewReflexService(repo)

	reflexes, err := service.ListReflexes(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to list reflexes: %v", err)
	}

	var reflexesVal []reflex.Reflex
	for _, r := range reflexes {
		reflexesVal = append(reflexesVal, *r)
	}

	return reflexesVal, nil
}
