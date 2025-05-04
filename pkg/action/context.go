package action

type Status string

const (
	StatusRunning   Status = "Running"
	StatusFailed    Status = "Failed"
	StatusCompleted Status = "Completed"
)

type SecurityContext struct {
	EventData       map[string]any    `json:"event_data"`
	Variables       map[string]any    `json:"variables"`
	ExecutionStatus map[string]Status `json:"execution_status"`
}
