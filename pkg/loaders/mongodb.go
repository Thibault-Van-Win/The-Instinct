package loaders

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/Thibault-Van-Win/The-Instinct/pkg/action"
	"github.com/Thibault-Van-Win/The-Instinct/pkg/reflex"
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
	var reflexes []reflex.Reflex

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(l.URI))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}
	defer client.Disconnect(ctx)

	// Ping to verify connection
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	collection := client.Database(l.DatabaseName).Collection(l.CollectionName)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to query MongoDB collection: %w", err)
	}
	defer cursor.Close(ctx)

	// Decode results into ReflexConfig objects
	var configs []reflex.ReflexConfig
	if err := cursor.All(ctx, &configs); err != nil {
		return nil, fmt.Errorf("failed to decode MongoDB documents: %w", err)
	}

	// Create reflexes from the configs
	for _, config := range configs {
		reflex, err := reflex.ReflexFromConfig(config, l.RuleRegistry, l.ActionRegistry)
		if err != nil {
			return nil, fmt.Errorf("failed to create reflex from MongoDB document: %w", err)
		}
		reflexes = append(reflexes, *reflex)
	}

	return reflexes, nil
}
