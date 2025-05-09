package scheduler

import "time"

// SchedulerConfig defines configuration for creating a scheduler
type SchedulerConfig struct {
	Type     string `yaml:"type"`
	EventURL string `yaml:"event_url"`
}

// TriggerConfig defines the configuration for a certain trigger
type TriggerConfig struct {
	ScheduleType string         `json:"schedule_type"` // "cron", "interval", etc.
	ID           string         `json:"id"`
	Name         string         `json:"name"`
	Description  string         `json:"description"`
	Schedule     string         `json:"schedule"` // Format depends on scheduler type
	EventData    map[string]any `json:"event_data"`
	Enabled      bool           `json:"enabled"`
	LastRun      time.Time      `json:"last_run,omitempty"`
}
