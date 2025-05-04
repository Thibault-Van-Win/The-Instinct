package action

import (
	"errors"
	"fmt"

	"github.com/mitchellh/mapstructure"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	ActionTypeSequential = "sequential"
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

	// This is dirty...
	// Primitive.A is not correctly converted a a slice...
	var childConfigs []any
	switch children := params["children"].(type) {
    case []any:
        // Standard Go slice - use directly
        childConfigs = children
    case primitive.A:
        // MongoDB primitive.A - convert to []any
        childConfigs = make([]any, len(children))
		copy(childConfigs, children)
    default:
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

func (sa *SequentialAction) Execute(ctx *SecurityContext) error {
	ctx.ExecutionStatus[sa.Name] = StatusRunning

	for _, child := range sa.Children {
		if err := child.Execute(ctx); err != nil {
			ctx.ExecutionStatus[sa.Name] = StatusFailed
			// Sequential action fails if one child fails
			return err
		}
	}

	ctx.ExecutionStatus[sa.Name] = StatusCompleted
	return nil
}

func (sa *SequentialAction) Validate() error {
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
