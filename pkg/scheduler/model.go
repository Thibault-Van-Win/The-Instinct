package scheduler

import (
	"context"
	"time"

	"github.com/Thibault-Van-Win/The-Instinct/pkg/triggerconfig"
)

type EventPublisher interface {
	PublishEvent(ctx context.Context, eventData map[string]any) error
}

type Scheduler interface {
	// AddTrigger adds a new trigger and returns its ID
	AddTrigger(config triggerconfig.TriggerConfig) (string, error)

	// UpdateTrigger updates an existing trigger
	UpdateTrigger(config triggerconfig.TriggerConfig) error

	// RemoveTrigger removes a trigger by ID
	RemoveTrigger(id string) error

	// GetTrigger retrieves a trigger by ID
	GetTrigger(id string) (triggerconfig.TriggerConfig, error)

	// ListTriggers returns all triggers
	ListTriggers(ctx context.Context) ([]triggerconfig.TriggerConfig, error)

	// Start starts the scheduler
	Start() error

	// Stop stops the scheduler
	Stop() context.Context
}

// BaseScheduler implements common functionality for all schedulers
type BaseScheduler struct {
	triggers       map[string]triggerconfig.TriggerConfig
	eventPublisher EventPublisher
}

func NewBaseScheduler(publisher EventPublisher) BaseScheduler {
	return BaseScheduler{
		triggers:       make(map[string]triggerconfig.TriggerConfig),
		eventPublisher: publisher,
	}
}

func (s *BaseScheduler) GetTrigger(id string) (triggerconfig.TriggerConfig, error) {
	trigger, exists := s.triggers[id]
	if !exists {
		return triggerconfig.TriggerConfig{}, ErrTriggerNotFound
	}
	return trigger, nil
}

func (s *BaseScheduler) ListTriggers(ctx context.Context) ([]triggerconfig.TriggerConfig, error) {
	triggers := make([]triggerconfig.TriggerConfig, 0, len(s.triggers))
	for _, trigger := range s.triggers {
		triggers = append(triggers, trigger)
	}
	return triggers, nil
}

func (s *BaseScheduler) createEventWithMetadata(trigger triggerconfig.TriggerConfig) map[string]any {
	eventData := make(map[string]any)
	for k, v := range trigger.EventData {
		eventData[k] = v
	}

	// Add trigger metadata
	eventData["trigger_id"] = trigger.ID
	eventData["trigger_name"] = trigger.Name
	eventData["trigger_type"] = "scheduled"
	eventData["schedule_type"] = trigger.ScheduleType
	eventData["trigger_time"] = time.Now().Format(time.RFC3339)

	return eventData
}

func (s *BaseScheduler) updateLastRun(id string) {
	if trigger, exists := s.triggers[id]; exists {
		trigger.LastRun = time.Now()
		s.triggers[id] = trigger
	}
}
