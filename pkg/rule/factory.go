package rule

import "fmt"

type RuleRegistry struct {
	factories map[string]RuleFactory
}

type RuleRegistryOption func(*RuleRegistry)

// RuleFactory is a function that creates a rule from parameters
type RuleFactory func(params map[string]any) (Rule, error)

// NewRuleRegistry creates a new rule registry
func NewRuleRegistry(opts ...RuleRegistryOption) *RuleRegistry {

	instance := &RuleRegistry{
		factories: make(map[string]RuleFactory),
	}

	for _, opt := range opts {
		opt(instance)
	}

	return instance
}

// Register registers an rule factory
func (r *RuleRegistry) Register(name string, factory RuleFactory) {
	r.factories[name] = factory
}

// Create a Rule from the config
func (r *RuleRegistry) Create(config RuleConfig) (Rule, error) {
	factory, ok := r.factories[config.Type]
	if !ok {
		return nil, fmt.Errorf("unknown rule type: %s", config.Type)
	}
	return factory(config.Params)
}

func (r *RuleRegistry) RegisterStandardRules() {
	// Cel rule
	r.Register(RuleTypeCel, func(params map[string]any) (Rule, error) {
		return NewCelRule(params)
	})
}

func WithStandardRules() RuleRegistryOption {
	return func(rr *RuleRegistry) {
		rr.RegisterStandardRules()
	}
}

func WithRuleFactory(name string, factory RuleFactory) RuleRegistryOption {
	return func(rr *RuleRegistry) {
		rr.Register(name, factory)
	}
}
