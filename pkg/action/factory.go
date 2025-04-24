package action

import "fmt"

// ActionConfig represents the structure of an action configuration
type ActionConfig struct {
	Type   string         `yaml:"type" json:"type"`
	Params map[string]any `yaml:"params" json:"params"`
}

// ActionRegistry is a registry of action factories
type ActionRegistry struct {
	factories map[string]ActionFactory
}

// ActionFactory is a function that creates an action from parameters
type ActionFactory func(params map[string]any) (Action, error)

// NewActionRegistry creates a new action registry
func NewActionRegistry() *ActionRegistry {
	return &ActionRegistry{
		factories: make(map[string]ActionFactory),
	}
}

// Register registers an action factory
func (r *ActionRegistry) Register(name string, factory ActionFactory) {
	r.factories[name] = factory
}

// Create creates an action from a configuration
func (r *ActionRegistry) Create(config ActionConfig) (Action, error) {
	factory, ok := r.factories[config.Type]
	if !ok {
		return nil, fmt.Errorf("unknown action type: %s", config.Type)
	}
	return factory(config.Params)
}

func (r *ActionRegistry) RegisterStandardActions() {
	// Print action
	r.Register("print", func(params map[string]any) (Action, error) {
		message, ok := params["message"].(string)
		if !ok {
			return nil, fmt.Errorf("print action requires a message parameter")
		}
		return PrintAction(message), nil
	})
}
