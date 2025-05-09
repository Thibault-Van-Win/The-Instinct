package action

import (
	"errors"
	"fmt"

	"github.com/Thibault-Van-Win/The-Instinct/pkg/security_context"
	"github.com/mitchellh/mapstructure"
)

const (
	ActionTypeSequential string = "sequential"
)

type SequentialAction struct {
	BaseAction
	Children []Action `json:"children"`
}

func NewSequentialAction(params map[string]any, reg *ActionRegistry) (*SequentialAction, error) {
	name, ok := params["name"].(string)
	if !ok {
		return nil, fmt.Errorf("a sequential action requires a name")
	}

	childConfigs, ok := params["children"].([]any)
	if !ok {
		return nil, fmt.Errorf("a sequential action requires children action configs, got %T", params["children"])
	}

	children := make([]Action, 0, len(childConfigs))
	for _, rawChildConfig := range childConfigs {
		// Need to convert to an ActionConfig
		var childConfig ActionConfig
		err := mapstructure.Decode(rawChildConfig, &childConfig)
		if err != nil {
			return nil, fmt.Errorf("child config has an unexpected structure: %v", err)
		}

		child, err := reg.Create(childConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to create child: %v", err)
		}
		children = append(children, child)
	}

	instance := &SequentialAction{
		BaseAction: BaseAction{
			Type: ActionTypeSequential,
			Name: name,
		},
		Children: children,
	}

	if err := instance.Validate(); err != nil {
		return nil, fmt.Errorf("sequential validation failed: %v", err)
	}

	return instance, nil
}

func (sa *SequentialAction) Execute(ctx *security_context.SecurityContext) error {
	ctx.ExecutionStatus[sa.Name] = security_context.StatusRunning

	for _, child := range sa.Children {
		if err := child.Execute(ctx); err != nil {
			ctx.ExecutionStatus[sa.Name] = security_context.StatusFailed
			// Sequential action fails if one child fails
			return err
		}
	}

	ctx.ExecutionStatus[sa.Name] = security_context.StatusCompleted
	return nil
}

func (sa *SequentialAction) Validate() error {
	if err := sa.BaseAction.Validate(); err != nil {
		return fmt.Errorf("basic validation failed: %v", err)
	}

	if len(sa.Children) == 0 {
		return fmt.Errorf("sequential action %s has no children", sa.Name)
	}

	var errs []error
	for _, child := range sa.Children {
		errs = append(errs, child.Validate())
	}

	if err := errors.Join(errs...); err != nil {
		return fmt.Errorf("sequential action %s has invalid children: %v", sa.Name, err)
	}

	return nil
}
