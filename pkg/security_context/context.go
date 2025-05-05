package security_context

type Status string

const (
	StatusRunning   Status = "Running"
	StatusFailed    Status = "Failed"
	StatusCompleted Status = "Completed"
)

// ? Should a lock be included here?
type SecurityContext struct {
	Event           map[string]any    `json:"event"`
	Variables       map[string]any    `json:"variables"`
	ExecutionStatus map[string]Status `json:"execution_status"`
}
