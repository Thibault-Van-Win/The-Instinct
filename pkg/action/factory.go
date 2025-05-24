package action

import (
	"fmt"
	"log"
	"os/exec"

	"github.com/hashicorp/go-plugin"
)

// ActionRegistry is a registry of action factories
type ActionRegistry struct {
	factories map[string]ActionFactory
	clients   []*plugin.Client
}

type ActionRegistryOption func(*ActionRegistry)

// ActionFactory is a function that creates an action from parameters
type ActionFactory func(params map[string]any) (Action, error)

// NewActionRegistry creates a new action registry
// This newly created action registry needs to be closed
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

func (r *ActionRegistry) Close() {
	for _, c := range r.clients {
		c.Kill()
	}
}

func WithStandardActions() ActionRegistryOption {
	return func(ar *ActionRegistry) {
		ar.registerStandardActions()
	}
}

func WithPlugins() ActionRegistryOption {
	return func(ar *ActionRegistry) {
		ar.registerPlugins()
	}
}

func WithActionFactory(name string, factory ActionFactory) ActionRegistryOption {
	return func(ar *ActionRegistry) {
		ar.Register(name, factory)
	}
}

func (r *ActionRegistry) registerStandardActions() {
	r.Register(ActionTypePrint, func(params map[string]any) (Action, error) {
		return NewPrintAction(params)
	})

	r.Register(ActionTypeSequential, func(params map[string]any) (Action, error) {
		return NewSequentialAction(params, r)
	})

	r.Register(ActionTypeParallel, func(params map[string]any) (Action, error) {
		return NewParallelAction(params, r)
	})

	r.Register(ActionTypeConditional, func(params map[string]any) (Action, error) {
		return NewConditionalAction(params, r)
	})

	r.Register(ActionTypeIterator, func(params map[string]any) (Action, error) {
		return NewIteratorAction(params, r)
	})
}

var handshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "BASIC_PLUGIN",
	MagicCookieValue: "hello",
}

// pluginMap is the map of plugins we can dispense.
var pluginMap = map[string]plugin.Plugin{
	"greeter": &ActionPlugin{},
}

func (r *ActionRegistry) registerPlugins() {
	action, client, err := loadPlugin("greeter")
	if err != nil {
		log.Printf("Failed to load greeter plugin: %v", err)
	}

	// Add the client so it's lifetime can be managed
	r.clients = append(r.clients, client)

	r.Register(
		action.GetType(),
		func(params map[string]any) (Action, error) {
			return NewPluginActionDecorator(action, params)
		},
	)
}

func loadPlugin(name string) (Action, *plugin.Client, error) {
	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
		Cmd:             exec.Command(fmt.Sprintf("./plugins/%s/%s", name, name)),
	})

	// Connect via RPC
	rpcClient, err := client.Client()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create rpc client for %s: %v", name, err)
	}

	// Request the plugin
	raw, err := rpcClient.Dispense(name)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to dispense plugin %s: %v", name, err)
	}

	action := raw.(Action)

	return action, client, nil
}
