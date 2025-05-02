package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/Thibault-Van-Win/The-Instinct/internal/config"
	"github.com/Thibault-Van-Win/The-Instinct/pkg/action"
	"github.com/Thibault-Van-Win/The-Instinct/pkg/api"
	"github.com/Thibault-Van-Win/The-Instinct/pkg/instinct"
	"github.com/Thibault-Van-Win/The-Instinct/pkg/loaders"
	"github.com/Thibault-Van-Win/The-Instinct/pkg/reflex"
	mongoRepo "github.com/Thibault-Van-Win/The-Instinct/pkg/reflex/mongo"
	"github.com/Thibault-Van-Win/The-Instinct/pkg/rule"
)

var (
	system *instinct.Instinct
)

func main() {
	conf, err := config.Instance()
	if err != nil {
		log.Fatalf("failed to retrieve config instance: %v", err)
	}

	// Create the rule registry
	ruleRegistry := rule.NewRuleRegistry(
		rule.WithStandardRules(),
	)

	// Create action registry
	actionRegistry := action.NewActionRegistry(
		action.WithStandardActions(),
	)

	// Create a new instinct system
	system = instinct.New(ruleRegistry, actionRegistry)

	// Load the reflexes
	dbConnString, err := conf.DbConfig.ConnString()
	if err != nil {
		log.Fatalf("failed to retrieve the db connection string: %v", err)
	}

	if err := system.LoadReflexes(loaders.MongoDB, map[string]any{
		"uri":        dbConnString,
		"database":   "instinct",
		"collection": "reflexes",
	}); err != nil {
		log.Fatalf("Failed to load reflexes: %v", err)
	}

	log.Printf("Loaded %d reflexes\n", len(system.Reflexes))

	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dbConnString))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	// Create database and collections
	db := client.Database("instinct")

	// Initialize repository and service (dependency injection)
	repository := mongoRepo.NewRepository(db, ruleRegistry, actionRegistry)
	service := reflex.NewReflexService(repository)

	e := echo.New()

	e.Use(
		middleware.Logger(),
		middleware.Recover(),
	)

	reflexController := api.NewReflexController(service)
	reflexController.Register(e)

	e.POST("/event", handleEvent)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", conf.WebServerConfig.Port)))
}

func handleEvent(c echo.Context) error {

	var event map[string]any
	err := c.Bind(&event)
	if err != nil {
		return c.String(http.StatusBadRequest, "Bad request")
	}

	fmt.Printf("Received the following event: %v\n", event)

	if err := system.ProcessEvent(event); err != nil {
		return c.String(http.StatusInternalServerError, "Something went wrong")
	}

	return c.NoContent(http.StatusNoContent)
}
