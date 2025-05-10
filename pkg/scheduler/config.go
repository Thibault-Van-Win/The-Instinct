package scheduler

// SchedulerConfig defines configuration for creating a scheduler
type SchedulerConfig struct {
	Type     string `yaml:"type"`
	EventURL string `yaml:"event_url"`
}
