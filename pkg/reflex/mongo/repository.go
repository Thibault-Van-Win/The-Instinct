package mongo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/mitchellh/mapstructure"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/Thibault-Van-Win/The-Instinct/internal/config"
	"github.com/Thibault-Van-Win/The-Instinct/pkg/action"
	"github.com/Thibault-Van-Win/The-Instinct/pkg/reflex"
	"github.com/Thibault-Van-Win/The-Instinct/pkg/rule"
)

const (
	FieldID           = "_id"
	FieldName         = "name"
	FieldRuleConfig   = "ruleConfig"
	FieldActionConfig = "actionConfig"
	FieldUpdatedAt    = "updatedAt"
)

// Database representation of Reflexes
type reflexDocument struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	Name         string             `bson:"name"`
	RuleConfig   rule.RuleConfig    `bson:"ruleConfig"`
	ActionConfig bson.Raw           `bson:"actionConfig"`
	CreatedAt    time.Time          `bson:"createdAt"`
	UpdatedAt    time.Time          `bson:"updatedAt"`
}

type Repository struct {
	client         *mongo.Client
	collection     *mongo.Collection
	ruleRegistry   *rule.RuleRegistry
	actionRegistry *action.ActionRegistry
}

func NewRepository(dbConfig *config.DatabaseConfig, ruleReg *rule.RuleRegistry, actionReg *action.ActionRegistry) (*Repository, error) {
	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dbConnString, err := dbConfig.ConnString()
	if err != nil {
		return nil, fmt.Errorf("failed to get connection string from config: %v", err)
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dbConnString))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to mongo db: %v", err)
	}

	// Create database and collections
	db := client.Database("instinct")
	collection := db.Collection("reflexes")

	// Create indexes for faster lookups
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: FieldName, Value: 1}},
		Options: options.Index().SetUnique(true),
	}

	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err = collection.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		// Index might already exists
	}

	return &Repository{
		client:         client,
		collection:     collection,
		ruleRegistry:   ruleReg,
		actionRegistry: actionReg,
	}, nil
}

// Close closes the MongoDB connection
func (r *Repository) Close(ctx context.Context) error {
	if r.client != nil {
		return r.client.Disconnect(ctx)
	}

	return nil
}

func (r *Repository) Create(ctx context.Context, config reflex.ReflexConfig) (string, error) {
	// Check if a reflex with the same name already exists
	existingReflex, _ := r.GetByName(ctx, config.Name)
	if existingReflex != nil {
		return "", errors.New("reflex with this name already exists")
	}

	rawActionConfig, err := bson.Marshal(config.ActionConfig)
	if err != nil {
		return "", fmt.Errorf("failed to marshal action config: %w", err)
	}

	now := time.Now()
	doc := reflexDocument{
		Name:         config.Name,
		RuleConfig:   config.RuleConfig,
		ActionConfig: rawActionConfig,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	result, err := r.collection.InsertOne(ctx, doc)
	if err != nil {
		return "", err
	}

	id, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", errors.New("failed to get inserted ID")
	}

	return id.Hex(), nil
}

func (r *Repository) GetByName(ctx context.Context, name string) (*reflex.Reflex, error) {
	var doc reflexDocument
	err := r.collection.FindOne(ctx, bson.M{FieldName: name}).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	return r.documentToReflex(doc)
}

func (r *Repository) GetByID(ctx context.Context, id string) (*reflex.Reflex, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid ID format")
	}

	var doc reflexDocument
	err = r.collection.FindOne(ctx, bson.M{FieldID: objectID}).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("reflex not found")
		}
		return nil, err
	}

	return r.documentToReflex(doc)
}

// List retrieves all reflexes
func (r *Repository) List(ctx context.Context) ([]*reflex.Reflex, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var docs []reflexDocument
	if err = cursor.All(ctx, &docs); err != nil {
		return nil, err
	}

	reflexes := make([]*reflex.Reflex, 0, len(docs))
	for _, doc := range docs {
		reflex, err := r.documentToReflex(doc)
		if err != nil {
			return nil, err
		}
		reflexes = append(reflexes, reflex)
	}

	return reflexes, nil
}

func (r *Repository) Update(ctx context.Context, id string, config reflex.ReflexConfig) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid ID format")
	}

	update := bson.M{
		"$set": bson.M{
			FieldName:         config.Name,
			FieldRuleConfig:   config.RuleConfig,
			FieldActionConfig: config.ActionConfig,
			FieldUpdatedAt:    time.Now(),
		},
	}

	result, err := r.collection.UpdateOne(ctx, bson.M{FieldID: objectID}, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("reflex not found")
	}

	return nil
}

// Delete removes a reflex by its ID
func (r *Repository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid ID format")
	}

	result, err := r.collection.DeleteOne(ctx, bson.M{FieldID: objectID})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("reflex not found")
	}

	return nil
}

// Convert a database reflex representation to the domain representation
func (r *Repository) documentToReflex(doc reflexDocument) (*reflex.Reflex, error) {
	// Convert rule configuration to actual rule
	ruleInstance, err := r.ruleRegistry.Create(doc.RuleConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create rule from config: %v", err)
	}

	// First unmarshal to a temporary structure with standard Go types
	var tempMap map[string]any
	if err := bson.Unmarshal(doc.ActionConfig, &tempMap); err != nil {
		return nil, fmt.Errorf("failed to decode action config: %w", err)
	}

	// Convert any MongoDB primitive types to standard Go types
	sanitizedMap := sanitizeBSONPrimitives(tempMap)

	// Now create an ActionConfig from the sanitized map
	var actionConfig action.ActionConfig
	if err := mapstructure.Decode(sanitizedMap, &actionConfig); err != nil {
		return nil, fmt.Errorf("failed to decode sanitized action config: %w", err)
	}

	actionInstance, err := r.actionRegistry.Create(actionConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create action from config: %v", err)
	}

	return reflex.NewReflex(
		doc.Name,
		ruleInstance,
		actionInstance,
	), nil
}

func sanitizeBSONPrimitives(input any) any {
	switch v := input.(type) {
	case primitive.A:
		// Convert MongoDB array to Go slice
		result := make([]any, len(v))
		for i, item := range v {
			result[i] = sanitizeBSONPrimitives(item)
		}
		return result

	case primitive.D:
		// Convert MongoDB document to Go map
		result := make(map[string]any)
		for _, elem := range v {
			result[elem.Key] = sanitizeBSONPrimitives(elem.Value)
		}
		return result

	case primitive.M:
		// Convert MongoDB map to Go map
		result := make(map[string]any)
		for key, val := range v {
			result[key] = sanitizeBSONPrimitives(val)
		}
		return result

	case map[string]any:
		// Recursively sanitize map values
		result := make(map[string]any)
		for key, val := range v {
			result[key] = sanitizeBSONPrimitives(val)
		}
		return result

	case []any:
		// Recursively sanitize slice elements
		result := make([]any, len(v))
		for i, item := range v {
			result[i] = sanitizeBSONPrimitives(item)
		}
		return result

	default:
		// Return other types as-is
		return v
	}
}
