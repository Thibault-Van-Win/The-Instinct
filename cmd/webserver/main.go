package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/Thibault-Van-Win/The-Instinct/internal/config"
	"github.com/Thibault-Van-Win/The-Instinct/internal/factory"
	"github.com/Thibault-Van-Win/The-Instinct/pkg/action"
	"github.com/Thibault-Van-Win/The-Instinct/pkg/controllers"
	"github.com/Thibault-Van-Win/The-Instinct/pkg/instinct"
	"github.com/Thibault-Van-Win/The-Instinct/pkg/reflex"
	"github.com/Thibault-Van-Win/The-Instinct/pkg/rule"
	"github.com/Thibault-Van-Win/The-Instinct/pkg/triggerconfig"
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
		action.WithPlugins(),
	)
	defer actionRegistry.Close()

	// Create a new instinct system
	system = instinct.New(ruleRegistry, actionRegistry)

	// Initialize repository and service (dependency injection)
	repository, err := factory.NewReflexRepository(&conf.DbConfig, ruleRegistry, actionRegistry)
	if err != nil {
		log.Fatalf("Failed to create reflex repository: %v", err)
	}
	service := reflex.NewReflexService(repository)
	defer service.Close(context.Background())

	// Load the reflexes
	if err := system.LoadReflexes(service); err != nil {
		log.Fatalf("Failed to load reflexes: %v", err)
	}
	log.Printf("Loaded %d reflexes\n", len(system.Reflexes))

	e := echo.New()

	e.Use(
		middleware.Logger(),
		middleware.Recover(),
	)

	// Register the needed controllers
	reflexController := controllers.NewReflexController(service)
	reflexController.Register(e)

	triggerconfigRepository, err := factory.NewTriggerConfigRepository(&conf.DbConfig)
	if err != nil {
		log.Fatalf("Failed to init trigger config repository: %v", err)
	}
	triggerconfigService := triggerconfig.NewTriggerConfigService(triggerconfigRepository)
	defer triggerconfigService.Close(context.Background())
	triggerconfigController := controllers.NewTriggerConfigController(triggerconfigService)
	triggerconfigController.Register(e)

	e.POST("/event", handleEvent)

	// Start the server
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		if err := e.Start(fmt.Sprintf(":%d", conf.WebServerConfig.Port)); err != nil && err != http.ErrServerClosed {
			log.Println("Shutting down server")
		}

	}()
	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Println("Gracefully terminating server")
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}

func handleEvent(c echo.Context) error {

	var event map[string]any
	err := c.Bind(&event)
	if err != nil {
		return c.String(http.StatusBadRequest, "Bad request")
	}

	fmt.Printf("Received the following event: %v\n", event)

	if err := system.ProcessEvent(event); err != nil {
		log.Println(err)
		return c.String(http.StatusInternalServerError, "Something went wrong")
	}

	return c.NoContent(http.StatusNoContent)
}
