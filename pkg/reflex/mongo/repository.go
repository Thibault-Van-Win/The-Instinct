package mongo

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/Thibault-Van-Win/The-Instinct/pkg/action"
	"github.com/Thibault-Van-Win/The-Instinct/pkg/reflex"
	"github.com/Thibault-Van-Win/The-Instinct/pkg/rule"
)

// Database representation of Reflexes
type reflexDocument struct {
	ID            primitive.ObjectID    `bson:"_id,omitempty"`
	Name          string                `bson:"name"`
	RuleConfig    rule.RuleConfig       `bson:"ruleConfig"`
	ActionConfigs []action.ActionConfig `bson:"actionConfigs"`
	CreatedAt     time.Time             `bson:"createdAt"`
	UpdatedAt     time.Time             `bson:"updatedAt"`
}

type Repository struct {
	collection     *mongo.Collection
	ruleRegistry   *rule.RuleRegistry
	actionRegistry *action.ActionRegistry
}

func NewRepository(db *mongo.Database, ruleReg *rule.RuleRegistry, actionReg *action.ActionRegistry) *Repository {
	collection := db.Collection("reflexes")

	// Create indexes for faster lookups
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "name", Value: 1}},
		Options: options.Index().SetUnique(true),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := collection.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		// Index might already exists
	}

	return &Repository{
		collection:     collection,
		ruleRegistry:   ruleReg,
		actionRegistry: actionReg,
	}
}

func (r *Repository) Create(ctx context.Context, config reflex.ReflexConfig) (string, error) {
	// Check if a reflex with the same name already exists
	existingReflex, _ := r.GetByName(ctx, config.Name)
	if existingReflex != nil {
		return "", errors.New("reflex with this name already exists")
	}

	now := time.Now()
	doc := reflexDocument{
		Name:          config.Name,
		RuleConfig:    config.RuleConfig,
		ActionConfigs: config.ActionConfigs,
		CreatedAt:     now,
		UpdatedAt:     now,
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
	err := r.collection.FindOne(ctx, bson.M{"name": name}).Decode(&doc)
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
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&doc)
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
			"name":          config.Name,
			"ruleConfig":    config.RuleConfig,
			"actionConfigs": config.ActionConfigs,
			"updatedAt":     time.Now(),
		},
	}

	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
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

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
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
		return nil, errors.New("failed to create rule from config: " + err.Error())
	}

	// Convert action configurations to actual actions
	actions := make([]action.Action, 0, len(doc.ActionConfigs))
	for _, actionConfig := range doc.ActionConfigs {
		actionInstance, err := r.actionRegistry.Create(actionConfig)
		if err != nil {
			return nil, errors.New("failed to create action from config: " + err.Error())
		}
		actions = append(actions, actionInstance)
	}

	return reflex.NewReflex(
		doc.Name,
		ruleInstance,
		actions,
	), nil
}
