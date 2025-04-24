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
)

var (
	system *instinct.Instinct
)

func init() {
	registry := action.NewActionRegistry()
	registry.RegisterStandardActions()

	system = instinct.New(registry)

	// Load the reflexes
	if err := system.LoadReflexes(loaders.YAMLLoader, map[string]any{
		"directory": "./config",
	}); err != nil {
		log.Fatalf("Failed to load reflexes from YAML: %v", err)
	}
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
	err := c.Bind(&event); if err != nil {
		return c.String(http.StatusBadRequest, "Bad request")
	}

	fmt.Printf("Received the following event: %v\n", event)

	if err := system.ProcessEvent(event); err != nil {
		return c.String(http.StatusInternalServerError, "Something went wrong")
	}

	return c.NoContent(http.StatusNoContent)
}