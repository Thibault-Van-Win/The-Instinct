package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Thibault-Van-Win/The-Instinct/internal/config"
	"github.com/Thibault-Van-Win/The-Instinct/internal/factory"
	"github.com/Thibault-Van-Win/The-Instinct/pkg/scheduler"
	"github.com/Thibault-Van-Win/The-Instinct/pkg/triggerconfig"
)

func main() {
	cronConfig := scheduler.SchedulerConfig{
		Type:     scheduler.TypeCron,
		EventURL: "http://localhost:8080/event",
	}

	generalConfig, err := config.Instance()
	if err != nil {
		log.Fatalf("Failed to retrieve general config: %v", err)
	}

	triggerConfigRepo, err := factory.NewTriggerConfigRepository(&generalConfig.DbConfig)
	if err != nil {
		log.Fatalf("Failed to create a trigger config repository: %v", err)
	}
	triggerService := triggerconfig.NewTriggerConfigService(triggerConfigRepo)
	defer triggerService.Close(context.Background())

	triggerconfigsPtr, err := triggerService.ListTriggerConfigs(context.Background())
	if err != nil {
		log.Fatalf("Failed to retrieve configured trigger configs: %v", err)
	}

	// Pointers need to be dereferenced
	var triggerconfigVals []triggerconfig.TriggerConfig
	for _, ptr := range triggerconfigsPtr {
		triggerconfigVals = append(triggerconfigVals, *ptr)
	}
	log.Printf("Loaded %d trigger configs from the DB", len(triggerconfigVals))

	reg, err := scheduler.NewSchedulerRegistry(
		scheduler.WithSchedulerFromConfig(cronConfig),
		scheduler.WithTriggerConfigs(triggerconfigVals),
	)
	if err != nil {
		log.Fatalf("Failed to initiate scheduler registry: %v", err)
	}

	reg.StartAll()

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	// Wait for signal
	<-stop
	log.Println("Shutting down schedulers...")
	reg.StopAll()
	log.Println("All schedulers stopped")
}
