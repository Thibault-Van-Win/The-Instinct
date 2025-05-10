package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Thibault-Van-Win/The-Instinct/pkg/scheduler"
	"github.com/Thibault-Van-Win/The-Instinct/pkg/triggerconfig"
)

func main() {
	// TODO: fetch all triggers from the DB

	cronConfig := scheduler.SchedulerConfig{
		Type:     scheduler.TypeCron,
		EventURL: "http://localhost:8080/event",
	}

	reg, err := scheduler.NewSchedulerRegistry(
		scheduler.WithSchedulerFromConfig(cronConfig),
		scheduler.WithTriggerConfigs([]triggerconfig.TriggerConfig{
			{
				ScheduleType: "cron",
				Name:         "Minute Trigger",
				Description:  "Runs every minute",
				Schedule:     "* * * * *",
				EventData: map[string]any{
					"type":   "maintenance",
					"action": "check",
					"source": "scheduler",
				},
				Enabled: true,
			},
		}),
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
