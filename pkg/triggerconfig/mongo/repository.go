package mongo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/Thibault-Van-Win/The-Instinct/internal/config"
	"github.com/Thibault-Van-Win/The-Instinct/pkg/triggerconfig"
	"github.com/labstack/gommon/log"
)

const (
	FieldID           = "_id"
	FieldName         = "name"
	FieldScheduleType = "schedule_type"
	FieldDescription  = "description"
	FieldSchedule     = "schedule"
	FieldEventData    = "event_data"
	FieldEnabled      = "enabled"
	FieldUpdatedAt    = "updatedAt"
)

// Database representation of trigger configs
type TriggerConfigDocument struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	Name         string             `bson:"name"`
	ScheduleType string             `bson:"schedule_type"`
	Description  string             `bson:"description"`
	Schedule     string             `bson:"schedule"`
	EventData    map[string]any     `bson:"event_data"`
	Enabled      bool               `bson:"enabled"`
	CreatedAt    time.Time          `bson:"createdAt"`
	UpdatedAt    time.Time          `bson:"updatedAt"`
}

type Repository struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func NewRepository(dbConfig *config.DatabaseConfig) (*Repository, error) {
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
	collection := db.Collection("triggerconfigs")

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
		// Warn just to be sure
		log.Warnf("failed to create index: %v", err)
	}

	return &Repository{
		client:     client,
		collection: collection,
	}, nil
}

// Close closes the MongoDB connection
func (r *Repository) Close(ctx context.Context) error {
	if r.client != nil {
		return r.client.Disconnect(ctx)
	}

	return nil
}

func (r *Repository) Create(ctx context.Context, config triggerconfig.TriggerConfig) (string, error) {
	// Check if a Trigger Config with the same name already exists
	existingConfig, _ := r.GetByName(ctx, config.Name)
	if existingConfig != nil {
		return "", errors.New("trigger config with this name already exists")
	}

	now := time.Now()
	doc := TriggerConfigDocument{
		ID:           primitive.NewObjectID(),
		Name:         config.Name,
		ScheduleType: config.ScheduleType,
		Description:  config.Description,
		Schedule:     config.Schedule,
		EventData:    config.EventData,
		Enabled:      config.Enabled,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	result, err := r.collection.InsertOne(ctx, doc)
	if err != nil {
		return "", fmt.Errorf("failed to insert trigger config: %v", err)
	}

	id, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", errors.New("failed to get inserted ID")
	}

	return id.Hex(), nil
}

func (r *Repository) GetByName(ctx context.Context, name string) (*triggerconfig.TriggerConfig, error) {
	var doc TriggerConfigDocument
	err := r.collection.FindOne(ctx, bson.M{FieldName: name}).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find in mongodb: %v", err)
	}

	return r.documentToTriggerConfig(doc)
}

func (r *Repository) GetByID(ctx context.Context, id string) (*triggerconfig.TriggerConfig, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid ID format")
	}

	var doc TriggerConfigDocument
	err = r.collection.FindOne(ctx, bson.M{FieldID: objectID}).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("trigger config document not found")
		}
		return nil, fmt.Errorf("failed to find in mongodb: %v", err)
	}

	return r.documentToTriggerConfig(doc)
}

// List retrieves all Trigger Configs
func (r *Repository) List(ctx context.Context) ([]*triggerconfig.TriggerConfig, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var docs []TriggerConfigDocument
	if err = cursor.All(ctx, &docs); err != nil {
		return nil, err
	}

	triggerConfigs := make([]*triggerconfig.TriggerConfig, 0, len(docs))
	for _, doc := range docs {
		triggerConfig, err := r.documentToTriggerConfig(doc)
		if err != nil {
			return nil, fmt.Errorf("failed to convert document to trigger config: %v", err)
		}
		triggerConfigs = append(triggerConfigs, triggerConfig)
	}

	return triggerConfigs, nil
}

func (r *Repository) Update(ctx context.Context, id string, config triggerconfig.TriggerConfig) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid ID format")
	}

	update := bson.M{
		"$set": bson.M{
			FieldName:         config.Name,
			FieldScheduleType: config.ScheduleType,
			FieldDescription:  config.Description,
			FieldSchedule:     config.Schedule,
			FieldEventData:    config.EventData,
			FieldEnabled:      config.Enabled,
			FieldUpdatedAt:    time.Now(),
		},
	}

	result, err := r.collection.UpdateOne(ctx, bson.M{FieldID: objectID}, update)
	if err != nil {
		return fmt.Errorf("failed to update in mongoDB: %v", err)
	}

	if result.MatchedCount == 0 {
		return errors.New("trigger config not found")
	}

	return nil
}

// Delete removes a trigger config by its ID
func (r *Repository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid ID format")
	}

	result, err := r.collection.DeleteOne(ctx, bson.M{FieldID: objectID})
	if err != nil {
		return fmt.Errorf("failed to delete in mongoDB: %v", err)
	}

	if result.DeletedCount == 0 {
		return errors.New("trigger config not found")
	}

	return nil
}

// Convert a database reflex representation to the domain representation
func (r *Repository) documentToTriggerConfig(doc TriggerConfigDocument) (*triggerconfig.TriggerConfig, error) {
	return &triggerconfig.TriggerConfig{
		ScheduleType: doc.ScheduleType,
		ID:           doc.ID.Hex(),
		Name:         doc.Name,
		Description:  doc.Description,
		Schedule:     doc.Schedule,
		EventData:    doc.EventData,
		Enabled:      doc.Enabled,
	}, nil
}
