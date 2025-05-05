package reflex

import (
	"fmt"

	"github.com/Thibault-Van-Win/The-Instinct/pkg/action"
	"github.com/Thibault-Van-Win/The-Instinct/pkg/rule"
	"github.com/Thibault-Van-Win/The-Instinct/pkg/security_context"
)

// Domain model for a reflex
type Reflex struct {
	Name   string        `json:"name"`
	Rule   rule.Rule     `json:"rule"`
	Action action.Action `json:"action"`
}

func (r *Reflex) Match(ctx *security_context.SecurityContext) (bool, error) {
	return r.Rule.Match(ctx)
}

func (r *Reflex) Execute(ctx *security_context.SecurityContext) error {
	return r.Action.Execute(ctx)
}

func NewReflex(name string, rule rule.Rule, act action.Action) *Reflex {
	return &Reflex{
		Name:   name,
		Rule:   rule,
		Action: act,
	}
}

func ReflexFromConfig(config ReflexConfig, ruleReg *rule.RuleRegistry, actionReg *action.ActionRegistry) (*Reflex, error) {
	ruleInstance, err := ruleReg.Create(config.RuleConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create rule: %v", err)
	}

	actionInstance, err := actionReg.Create(config.ActionConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create action: %v", err)
	}

	return NewReflex(config.Name, ruleInstance, actionInstance), nil
}
