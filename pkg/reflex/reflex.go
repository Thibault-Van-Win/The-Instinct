package reflex

import (
	"github.com/Thibault-Van-Win/The-Instinct/pkg/action"
	"github.com/Thibault-Van-Win/The-Instinct/pkg/rule"
)

type Reflex struct {
	Name    string
	Rule    rule.Rule
	Actions []action.Action
}

// ReflexConfig represents the structure of a reflex configuration
type ReflexConfig struct {
	Name       string                `yaml:"name" json:"name"`
	RuleConfig rule.RuleConfig       `yaml:"rule" json:"rule"`
	Actions    []action.ActionConfig `yaml:"actions" json:"actions"`
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
