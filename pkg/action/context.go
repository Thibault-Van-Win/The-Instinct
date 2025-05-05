package action

type Status string

const (
	StatusRunning   Status = "Running"
	StatusFailed    Status = "Failed"
	StatusCompleted Status = "Completed"
)

type SecurityContext struct {
	Event           map[string]any    `json:"event"`
	Variables       map[string]any    `json:"variables"`
	ExecutionStatus map[string]Status `json:"execution_status"`
}
