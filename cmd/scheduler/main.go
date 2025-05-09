package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Thibault-Van-Win/The-Instinct/pkg/scheduler"
)

func main() {
	cronConfig := scheduler.SchedulerConfig{
		Type:     scheduler.TypeCron,
		EventURL: "http://localhost:8080/event",
	}

	cronScheduler, err := scheduler.NewScheduler(cronConfig)
	if err != nil {
		log.Fatalf("Failed to create scheduler: %v", err)
	}

	// TODO: fetch all triggers from the DB
	// TODO: loop over all trigger configs and add them in the correct scheduler
	//? Might need a trigger registry to allow a mapping of which trigger config belongs to which scheduler

	cronTrigger := scheduler.TriggerConfig{
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
	}

	cronID, err := cronScheduler.AddTrigger(cronTrigger)
	if err != nil {
		log.Fatalf("Failed to add cron trigger: %v", err)
	}
	log.Printf("Added cron trigger with ID: %s", cronID)

	cronScheduler.Start()

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	// Wait for signal
	<-stop
	log.Println("Shutting down scheduler...")

	ctx := cronScheduler.Stop()

	// Wait for ongoing jobs to finish
	<-ctx.Done()
	log.Println("Scheduler stopped gracefully")
}
