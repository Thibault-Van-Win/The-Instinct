package main

import (
	"log"

	"github.com/hashicorp/go-plugin"

	"github.com/Thibault-Van-Win/The-Instinct/pkg/action"
	"github.com/Thibault-Van-Win/The-Instinct/pkg/security_context"
)

type Greeter struct {
}

// Implement the action interface
func (g *Greeter) Execute(ctx *security_context.SecurityContext) error {
	log.Printf("Hello from the plugin, this is my config: %v\n", ctx)
	return nil
}

func (g *Greeter) GetType() string {
	return "greeter"
}

func (g *Greeter) GetName() string {
	return "Plugin Greeter"
}

func (g *Greeter) Validate() error {
	return nil
}

var handshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "BASIC_PLUGIN",
	MagicCookieValue: "hello",
}

func main() {
	greeter := &Greeter{}

	// TODO: action plugin will need to be extracted from the action package
	// TODO: otherwise, this will lead to an import cycle in the factory
	// TODO: this suffice for some testing. Either plugin needs its own package or
	// TODO: the factories need to go to the internal package
	var pluginMap = map[string]plugin.Plugin{
		"greeter": &action.ActionPlugin{Impl: greeter},
	}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
	})
}
