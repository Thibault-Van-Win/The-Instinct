package reflex

import (
	"fmt"

	"github.com/Thibault-Van-Win/The-Instinct/pkg/action"
	"github.com/Thibault-Van-Win/The-Instinct/pkg/rule"
)

// Domain model for a reflex
type Reflex struct {
	Name    string
	Rule    rule.Rule
	Actions []action.Action
}

// Domain model for a reflex configuration
type ReflexConfig struct {
	Name          string                `yaml:"name" json:"name"`
	RuleConfig    rule.RuleConfig       `yaml:"rule" json:"rule"`
	ActionConfigs []action.ActionConfig `yaml:"actions" json:"actions"`
}

func (r *Reflex) Match(data map[string]any) (bool, error) {
	return r.Rule.Match(data)
}

func (r *Reflex) Do() error {
	for _, act := range r.Actions {
		act.Do()
	}

	return nil
}

func NewReflex(name string, rule rule.Rule, actions []action.Action) *Reflex {
	return &Reflex{
		Name:    name,
		Rule:    rule,
		Actions: actions,
	}
}

func ReflexFromConfig(config ReflexConfig, ruleReg *rule.RuleRegistry, actionReg *action.ActionRegistry) (*Reflex, error) {
	ruleInstance, err := ruleReg.Create(config.RuleConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create rule: %w", err)
	}

	// Create the actions
	actions := make([]action.Action, 0, len(config.ActionConfigs))
	for _, actionConfig := range config.ActionConfigs {
		actionInstance, err := actionReg.Create(actionConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to create action for reflex %s: %w", config.Name, err)
		}
		actions = append(actions, actionInstance)
	}

	return NewReflex(config.Name, ruleInstance, actions), nil
}
