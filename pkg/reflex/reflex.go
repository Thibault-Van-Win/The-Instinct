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

func (r *Reflex) Match(data map[string]any) (bool, error) {
	return r.Rule.Match(data)
}

func (r *Reflex) Do() error {
	for _, action := range r.Actions {
		action.Do()
	}

	return nil
}
