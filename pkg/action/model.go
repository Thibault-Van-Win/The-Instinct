package action

import (
	"fmt"

	"github.com/Thibault-Van-Win/The-Instinct/pkg/security_context"
)

type Action interface {
	Execute(ctx *security_context.SecurityContext) error
	GetType() string
	GetName() string
	Validate() error
}

// BaseAction provides common functionality for all actions
type BaseAction struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

func (ba *BaseAction) GetType() string {
	return ba.Type
}

func (ba *BaseAction) GetName() string {
	return ba.Name
}

func (ba *BaseAction) Validate() error {
	if ba.Type == "" {
		return fmt.Errorf("missing type: all actions must have a type")
	}

	if ba.Name == "" {
		return fmt.Errorf("missing name: all actions must have a unique name")
	}

	return nil
}
