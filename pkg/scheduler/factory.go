package scheduler

import (
	"fmt"
	"log"
)

// Factory method
func newScheduler(config SchedulerConfig) (Scheduler, error) {
	if config.EventURL == "" {
		return nil, fmt.Errorf("event URL is required when using HTTP publisher")
	}

	// Currently, only one that is supported
	publisher := NewHTTPEventPublisher(config.EventURL)

	// Create the appropriate scheduler
	switch config.Type {
	case TypeCron:
		return NewCronScheduler(publisher), nil
	default:
		return nil, fmt.Errorf("unknown scheduler type: %s", config.Type)
	}
}

type SchedulerRegistry struct {
	Schedulers map[string]Scheduler
}

func (sr *SchedulerRegistry) StartAll() {
	for _, scheduler := range sr.Schedulers {
		scheduler.Start()
	}
}

func (sr *SchedulerRegistry) StopAll() {
	for typ, scheduler := range sr.Schedulers {
		ctx := scheduler.Stop()
		<-ctx.Done()
		log.Printf("%s Scheduler stopped gracefully", typ)
	}
}

type SchedulerRegistryOption func(*SchedulerRegistry) error

func NewSchedulerRegistry(opts ...SchedulerRegistryOption) (*SchedulerRegistry, error) {
	instance := SchedulerRegistry{
		Schedulers: make(map[string]Scheduler),
	}

	for _, opt := range opts {
		if err := opt(&instance); err != nil {
			return nil, fmt.Errorf("failed to apply option: %v", err)
		}
	}

	return &instance, nil
}

func WithSchedulerFromConfig(config SchedulerConfig) SchedulerRegistryOption {
	return func(sr *SchedulerRegistry) error {
		instance, err := newScheduler(config)
		if err != nil {
			return fmt.Errorf("failed to create scheduler from config (type: %s): %v", config.Type, err)
		}

		sr.Schedulers[config.Type] = instance
		return nil
	}
}

func WithTriggerConfigs(triggerConfigs []TriggerConfig) SchedulerRegistryOption {
	return func(sr *SchedulerRegistry) error {
		for _, triggerConfig := range triggerConfigs {
			instance, ok := sr.Schedulers[triggerConfig.ScheduleType]
			if !ok {
				return fmt.Errorf("failed to create registry: no %s scheduler", triggerConfig.ScheduleType)
			}

			cronID, err := instance.AddTrigger(triggerConfig)
			if err != nil {
				return fmt.Errorf("failed to add trigger config (name: %s) to scheduler: %v", triggerConfig.Name, err)
			}
			log.Printf("Added cron trigger with ID: %s", cronID)
		}

		return nil
	}
}
