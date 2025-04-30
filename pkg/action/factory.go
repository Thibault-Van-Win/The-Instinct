package action

import "fmt"

// ActionRegistry is a registry of action factories
type ActionRegistry struct {
	factories map[string]ActionFactory
}

type ActionRegistryOption func(*ActionRegistry)

// ActionFactory is a function that creates an action from parameters
type ActionFactory func(params map[string]any) (Action, error)

// NewActionRegistry creates a new action registry
func NewActionRegistry(opts ...ActionRegistryOption) *ActionRegistry {

	instance := &ActionRegistry{
		factories: make(map[string]ActionFactory),
	}

	for _, opt := range opts {
		opt(instance)
	}

	return instance
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
		return NewPrintAction(message), nil
	})
}

func WithStandardActions() ActionRegistryOption {
	return func(ar *ActionRegistry) {
		ar.RegisterStandardActions()
	}
}

func WithActionFactory(name string, factory ActionFactory) ActionRegistryOption {
	return func(ar *ActionRegistry) {
		ar.Register(name, factory)
	}
}
