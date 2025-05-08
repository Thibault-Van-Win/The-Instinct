package security_context

import (
	"encoding/json"
	"fmt"
)

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

// Cannot give the ctx directly as the eval expect a map[string]any
// Use the json tags to marshall this into a map
// Another option would be to create a new map, benefits:
// 	- Faster
//	- Type safety
// Negatives:
//	- Need to add a new field for each extension
func (ctx *SecurityContext) ToMap() (map[string]any, error) {
	jsonBytes, err := json.Marshal(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal context: %v", err)
	}

	var activation map[string]any
	if err := json.Unmarshal(jsonBytes, &activation); err != nil {
		return nil, fmt.Errorf("failed to unmarshal into map: %v", err)
	}

	return activation, nil
}
