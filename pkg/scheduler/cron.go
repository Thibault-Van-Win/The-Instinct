package scheduler

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/robfig/cron/v3"

	"github.com/Thibault-Van-Win/The-Instinct/pkg/triggerconfig"
)

const (
	TypeCron string = "cron"
)

// CronScheduler implements Scheduler using cron expressions
type CronScheduler struct {
	BaseScheduler
	cron      *cron.Cron
	entryIDs  map[string]cron.EntryID
	mu        sync.RWMutex
	isRunning bool
}

// NewCronScheduler creates a new cron-based scheduler
func NewCronScheduler(publisher EventPublisher) *CronScheduler {
	//cronOptions := cron.WithSeconds()
	return &CronScheduler{
		BaseScheduler: NewBaseScheduler(publisher),
		cron:          cron.New(),
		entryIDs:      make(map[string]cron.EntryID),
		isRunning:     false,
	}
}

// AddTrigger adds a new time-based trigger
func (s *CronScheduler) AddTrigger(config triggerconfig.TriggerConfig) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Set schedule type if not specified
	if config.ScheduleType == "" {
		config.ScheduleType = TypeCron
	} else if config.ScheduleType != TypeCron {
		return "", fmt.Errorf("invalid schedule type for CronScheduler: %s", config.ScheduleType)
	}

	// Generate ID if not provided
	if config.ID == "" {
		config.ID = uuid.New().String()
	}

	// Validate cron expression
	if _, err := cron.ParseStandard(config.Schedule); err != nil {
		return "", fmt.Errorf("%w: %v", ErrInvalidSchedule, err)
	}

	// Add to triggers map
	s.triggers[config.ID] = config

	// Schedule if enabled and running
	if config.Enabled && s.isRunning {
		if err := s.scheduleTrigger(config); err != nil {
			delete(s.triggers, config.ID)
			return "", err
		}
	}

	return config.ID, nil
}

// UpdateTrigger updates an existing trigger
func (s *CronScheduler) UpdateTrigger(config triggerconfig.TriggerConfig) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.triggers[config.ID]; !exists {
		return ErrTriggerNotFound
	}

	// Validate cron expression
	if _, err := cron.ParseStandard(config.Schedule); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidSchedule, err)
	}

	// If running, remove the existing schedule
	if s.isRunning {
		if entryID, exists := s.entryIDs[config.ID]; exists {
			s.cron.Remove(entryID)
			delete(s.entryIDs, config.ID)
		}
	}

	// Update in triggers map
	s.triggers[config.ID] = config

	// Reschedule if enabled and running
	if config.Enabled && s.isRunning {
		return s.scheduleTrigger(config)
	}

	return nil
}

// RemoveTrigger removes a trigger by ID
func (s *CronScheduler) RemoveTrigger(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.triggers[id]; !exists {
		return ErrTriggerNotFound
	}

	// If running, remove the schedule
	if s.isRunning {
		if entryID, exists := s.entryIDs[id]; exists {
			s.cron.Remove(entryID)
			delete(s.entryIDs, id)
		}
	}

	// Remove from triggers map
	delete(s.triggers, id)

	return nil
}

// scheduleTrigger adds the trigger to the cron scheduler
func (s *CronScheduler) scheduleTrigger(config triggerconfig.TriggerConfig) error {
	entryID, err := s.cron.AddFunc(config.Schedule, func() {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		triggerID := config.ID // Capture for closure

		// Create event with metadata
		eventData := s.createEventWithMetadata(config)

		// Publish event
		if err := s.eventPublisher.PublishEvent(ctx, eventData); err != nil {
			log.Printf("Error publishing event for trigger %s: %v", triggerID, err)
			return
		}

		// Update last run time
		s.mu.Lock()
		s.updateLastRun(triggerID)
		s.mu.Unlock()
	})

	if err != nil {
		return fmt.Errorf("failed to schedule trigger: %w", err)
	}

	s.entryIDs[config.ID] = entryID
	return nil
}

// Start starts the scheduler
func (s *CronScheduler) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.isRunning {
		return ErrAlreadyStarted
	}

	// Schedule all enabled triggers
	for _, trigger := range s.triggers {
		if trigger.Enabled {
			if err := s.scheduleTrigger(trigger); err != nil {
				// Cleanup any scheduled triggers
				for scheduledID, entryID := range s.entryIDs {
					s.cron.Remove(entryID)
					delete(s.entryIDs, scheduledID)
				}
				return err
			}
		}
	}

	s.cron.Start()
	s.isRunning = true
	log.Println("Cron scheduler started")

	return nil
}

// Stop stops the scheduler
func (s *CronScheduler) Stop() context.Context {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.isRunning {
		ctx, _ := context.WithCancel(context.Background())
		return ctx
	}

	ctx := s.cron.Stop()
	s.isRunning = false
	s.entryIDs = make(map[string]cron.EntryID)
	log.Println("Cron scheduler stopped")

	return ctx
}

// GetTrigger retrieves a trigger by ID
func (s *CronScheduler) GetTrigger(id string) (triggerconfig.TriggerConfig, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.BaseScheduler.GetTrigger(id)
}

// ListTriggers returns all triggers
func (s *CronScheduler) ListTriggers(ctx context.Context) ([]triggerconfig.TriggerConfig, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.BaseScheduler.ListTriggers(ctx)
}
