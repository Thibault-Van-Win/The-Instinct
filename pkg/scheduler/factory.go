package scheduler

import "fmt"

type SchedulerType string

const (
	TypeCron SchedulerType = "cron"
)

// SchedulerConfig defines configuration for creating a scheduler
type SchedulerConfig struct {
	Type     SchedulerType `yaml:"type"`
	EventURL string        `yaml:"event_url"`
}

func NewScheduler(config SchedulerConfig) (Scheduler, error) {
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
