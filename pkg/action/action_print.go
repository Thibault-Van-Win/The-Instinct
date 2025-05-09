package action

import (
	"fmt"

	"github.com/Thibault-Van-Win/The-Instinct/pkg/security_context"
)

const (
	ActionTypePrint = "print"
)

type PrintAction struct {
	BaseAction
	Message string `json:"message"`
}

func NewPrintAction(params map[string]any) (*PrintAction, error) {
	name, ok := params["name"].(string)
	if !ok {
		return nil, fmt.Errorf("print action requires a name")
	}

	message, ok := params["message"].(string)
	if !ok {
		return nil, fmt.Errorf("print action requires a message parameter")
	}
	return &PrintAction{
		BaseAction: BaseAction{
			Type: ActionTypePrint,
			Name: name,
		},
		Message: message,
	}, nil
}

// PrintAction creates a simple action that prints a message
func (pa *PrintAction) Execute(ctx *security_context.SecurityContext) error {
	fmt.Println(pa.Message)
	ctx.ExecutionStatus[pa.Name] = security_context.StatusCompleted
	return nil
}

func (pa *PrintAction) Validate() error {
	if err := pa.BaseAction.Validate(); err != nil {
		return fmt.Errorf("basic validation failed: %v", err)
	}

	if pa.Message == "" {
		return fmt.Errorf("message cannot be empty")
	}

	return nil
}
