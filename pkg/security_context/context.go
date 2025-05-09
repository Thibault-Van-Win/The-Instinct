package security_context

import (
	"encoding/json"
	"errors"
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

type SecurityContextOption func(*SecurityContext) error

// Cannot give the ctx directly as the eval expect a map[string]any
// Use the json tags to marshall this into a map
// Another option would be to create a new map, benefits:
//   - Faster
//   - Type safety
//
// Negatives:
//   - Need to add a new field for each extension
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

func New(opts ...SecurityContextOption) (*SecurityContext, error) {
	instance := SecurityContext{
		Event:           make(map[string]any),
		Variables:       make(map[string]any),
		ExecutionStatus: make(map[string]Status),
	}

	for _, opt := range opts {
		opt(&instance)
	}

	return &instance, nil
}

func WithEvent(event map[string]any) SecurityContextOption {
	return func(sc *SecurityContext) error {
		if event == nil {
			return errors.New("event map cannot be nil")
		}

		sc.Event = event
		return nil
	}
}

func WithVariable(key string, value any) SecurityContextOption {
	return func(sc *SecurityContext) error {
		if key == "" {
			return errors.New("key for variable cannot be empty")
		}

		sc.Variables[key] = value

		return nil
	}
}

func WithVariables(vars map[string]any) SecurityContextOption {
	return func(sc *SecurityContext) error {
		for k, v := range vars {
			if k == "" {
				return errors.New("variable key cannot be empty")
			}
			sc.Variables[k] = v
		}
		return nil
	}
}
