package action

import "github.com/Thibault-Van-Win/The-Instinct/pkg/security_context"

// Does not need a type const as the name of
// the plugin will be used in the factory

type PluginActionDecorator struct {
	BaseAction
	Wrapped Action
	Config map[string]any
}

func NewPluginActionDecorator(wrapped Action, config map[string]any) (*PluginActionDecorator, error) {
	return &PluginActionDecorator{
		Wrapped: wrapped,
		Config:  config,
	}, nil
}

func (pad *PluginActionDecorator) Execute(ctx *security_context.SecurityContext) error {
	// Inject the config for the plugin
	ctx.Variables["config"] = pad.Config

	return pad.Wrapped.Execute(ctx)
}

func (pad *PluginActionDecorator) GetType() string {
	return pad.Wrapped.GetType()
}

func (pad *PluginActionDecorator) GetName() string {
	return pad.Wrapped.GetName()
}

func (pad *PluginActionDecorator) Validate() error {
	return pad.Wrapped.Validate()
}
