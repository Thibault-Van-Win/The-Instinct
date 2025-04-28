package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/Thibault-Van-Win/The-Instinct/pkg/action"
	"github.com/Thibault-Van-Win/The-Instinct/pkg/instinct"
	"github.com/Thibault-Van-Win/The-Instinct/pkg/loaders"
	"github.com/Thibault-Van-Win/The-Instinct/pkg/rule"
)

var (
	system *instinct.Instinct
)

func init() {
	// Create action registry
	actionRegistry := action.NewActionRegistry()
	actionRegistry.RegisterStandardActions()

	// Create the rule registry
	ruleRegistry := rule.NewRuleRegistry()
	ruleRegistry.RegisterStandardRules()

	// Create a new instinct system
	system = instinct.New(ruleRegistry, actionRegistry)
	// Load the reflexes
	if err := system.LoadReflexes(loaders.MongoDB, map[string]any{
		"uri":        "mongodb://user:secret@localhost:27017",
		"database":   "instinct",
		"collection": "reflexes",
	}); err != nil {
		log.Fatalf("Failed to load reflexes: %v", err)
	}

	log.Printf("Loaded %d reflexes\n", len(system.Reflexes))
}

func main() {
	e := echo.New()

	e.Use(
		middleware.Logger(),
		middleware.Recover(),
	)

	e.POST("/event", handleEvent)

	e.Logger.Fatal(e.Start(":8080"))
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
