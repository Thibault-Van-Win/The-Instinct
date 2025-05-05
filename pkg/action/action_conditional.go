package action

import (
	"fmt"

	"github.com/Thibault-Van-Win/The-Instinct/pkg/rule"
	"github.com/Thibault-Van-Win/The-Instinct/pkg/security_context"
)

const (
	ActionTypeConditional = "conditional"
)

type ConditionalAction struct {
	BaseAction
	Matcher Matcher `json:"matcher"`
	ThenAction Action `json:"then_action"`
	ElseAction Action `json:"else_action"`
}

type Matcher interface {
	Match(ctx *security_context.SecurityContext) (bool, error)
}

func NewConditionalAction(params map[string]any, reg *ActionRegistry) (*ConditionalAction, error) {
	name, ok := params["name"].(string)
	if !ok {
		return nil, fmt.Errorf("conditional action requires a name")
	}

	// Create rule matcher
	matcherRawConf, ok := params["rule_config"]
	if !ok {
		return nil, fmt.Errorf("conditional action requires a rule config")
	}

	matcherConf, err := rule.ConvertToConfig(matcherRawConf)
	if err != nil {
		return nil, fmt.Errorf("failed to decode rule_config: %v", err)
	}

	// Still a bit dirty
	ruleReg := rule.NewRuleRegistry(
		rule.WithStandardRules(),
	)

	rule, err := ruleReg.Create(*matcherConf)
	if err != nil {
		return nil, fmt.Errorf("failed to create matcher for conditional: %v", err)
	}

	// Create then action
	thenRawActionConfig, ok := params["then_action"]
	if !ok {
		return nil, fmt.Errorf("conditional action requires a 'then_action'")
	}

	thenActionConfig, err := convertToConfig(thenRawActionConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to decode then_action config: %v", err)
	}

	thenAction, err := reg.Create(*thenActionConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create then_action: %v", err)
	}

	// Create else action
	var elseAction Action
	elseRawActionConfig, ok := params["else_action"]
	if !ok {
		// Else action is not mandatory
		elseAction = nil
	} else {
		elseActionConfig, err := convertToConfig(elseRawActionConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to decode else_action: %v", err)
		}

		elseAction, err = reg.Create(*elseActionConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to create else_action: %v", err)
		}
	}

	instance := &ConditionalAction{
		BaseAction: BaseAction{
			Type: ActionTypeConditional,
			Name: name,
		},
		Matcher: rule,
		ThenAction: thenAction,
		ElseAction: elseAction,
	}

	if err := instance.Validate(); err != nil {
		return nil, fmt.Errorf("conditional validation failed: %v", err)
	}

	return instance, nil
}

func (ca *ConditionalAction) Execute(ctx *security_context.SecurityContext) error {
	ctx.ExecutionStatus[ca.Name] = security_context.StatusRunning

	conditionMet, err := ca.Matcher.Match(ctx)
	if err != nil {
		ctx.ExecutionStatus[ca.Name] = security_context.StatusFailed
		return fmt.Errorf("conditional failed: failed to evaluate condition: %v", err)
	}

	if conditionMet {
		err = ca.ThenAction.Execute(ctx)
	} else {
		if ca.ElseAction != nil {
			err = ca.ElseAction.Execute(ctx)
		}
	}

	if err != nil {
		ctx.ExecutionStatus[ca.Name] = security_context.StatusFailed
		return fmt.Errorf("failed to execute conditional action: %v", err)
	}

	ctx.ExecutionStatus[ca.Name] = security_context.StatusCompleted
	return nil
}

func (ca *ConditionalAction) Validate() error {
	if err := ca.BaseAction.Validate(); err != nil {
		return fmt.Errorf("basic validation failed: %v", err)
	}

	if ca.Matcher == nil {
		return fmt.Errorf("conditional action %s has no rule matcher", ca.Name)
	}

	if ca.ThenAction == nil {
		return fmt.Errorf("conditional action %s has no 'then' action", ca.Name)
	}

	if err := ca.ThenAction.Validate(); err != nil {
		return fmt.Errorf("conditional action %s has an invalid 'then' action: %v", ca.Name, err)
	}

	if ca.ElseAction == nil {
		return nil
	}

	if err := ca.ElseAction.Validate(); err != nil {
		return fmt.Errorf("conditional action %s has an invalid 'else' action: %v", ca.Name, err)
	}

	return nil
}
