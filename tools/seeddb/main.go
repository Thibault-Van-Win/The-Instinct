package main

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/yaml.v3"

	"github.com/Thibault-Van-Win/The-Instinct/pkg/reflex"
)

var (
	configDir = "./config/security_reflexes.yaml"
)

func main() {
	log.Println("Seeding MongoDB database")
	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Use the connection string with authentication
	uri := "mongodb://user:secret@localhost:27017"
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(ctx)

	// Ping to check connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}
	log.Println("Connected to MongoDB successfully")

	// Load the yaml rules
	data, err := os.ReadFile(configDir)
	if err != nil {
		log.Fatalf("Failed to read the config file %s: %v", configDir, err)
	}

	// Parse the YAML
	var configs []reflex.ReflexConfig
	if err := yaml.Unmarshal(data, &configs); err != nil {
		log.Fatalf("Failed to parse YAML in file %s: %v", configDir, err)
	}

	// Add to MongoDB
	// Push to the instinct db, reflexes collection
	db := client.Database("instinct")
	collection := db.Collection("reflexes")

	// Insert or update each reflex config
	for _, config := range configs {
		filter := bson.M{"name": config.Name}

		// Upsert the docs
		opts := options.Replace().SetUpsert(true)

		// Use ReplaceOne to either insert a new document or replace an existing one
		result, err := collection.ReplaceOne(ctx, filter, config, opts)
		if err != nil {
			log.Printf("Error upserting reflex %s: %v", config.Name, err)
			continue
		}

		if result.UpsertedCount > 0 {
			log.Printf("Added new reflex: %s", config.Name)
		} else if result.ModifiedCount > 0 {
			log.Printf("Updated existing reflex: %s", config.Name)
		} else {
			log.Printf("No changes needed for reflex: %s", config.Name)
		}
	}

	log.Println("Database seeding completed successfully")

}
