// Package for a TheHive plugin
//
// Accepts Action calls where the config can be found under ctx.Variables["config"]
// This should contain a sub action alongside some params for it
// Credentials for TheHive should be placed in the environment
package main

import (
	"log"

	"github.com/hashicorp/go-plugin"

	"github.com/Thibault-Van-Win/The-Instinct/pkg/action"
	"github.com/Thibault-Van-Win/The-Instinct/pkg/security_context"
)

type TheHive struct {
}

func (th *TheHive) Execute(ctx *security_context.SecurityContext) error {
	cfg := ctx.Variables["config"].(map[string]any)
	log.Printf("Hello from the TheHive plugin, here is my config %v", cfg)
	return nil
}

func (th *TheHive) GetType() string {
	return "thehive"
}

func (th *TheHive) GetName() string {
	return "A plugin for some Hive instance"
}

func (th *TheHive) Validate() error {
	return nil
}

var handshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "BASIC_PLUGIN",
	MagicCookieValue: "hello",
}

func main() {
	instance := &TheHive{}

	var pluginMap = map[string]plugin.Plugin{
		"thehive": &action.ActionPlugin{Impl: instance},
	}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: handshakeConfig,
		Plugins: pluginMap,
	})
}
