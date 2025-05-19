package main

import (
	"fmt"
	"log"
	"os/exec"

	"github.com/Thibault-Van-Win/The-Instinct/pkg/action"
	"github.com/Thibault-Van-Win/The-Instinct/pkg/security_context"
	"github.com/hashicorp/go-plugin"
)

func main() {

	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
		Cmd:             exec.Command("./plugins/greeter"),
	})
	defer client.Kill()

	// Connect via RPC
	rpcClient, err := client.Client()
	if err != nil {
		log.Fatal(err)
	}

	// Request the plugin
	raw, err := rpcClient.Dispense("greeter")
	if err != nil {
		log.Fatal(err)
	}

	greeter := raw.(action.Action)
	fmt.Printf("Plugin name: %s\n", greeter.GetName())
	fmt.Printf("Plugin type: %s\n", greeter.GetType())
	fmt.Printf("Plugin validity: %s\n", greeter.Validate())

	ctx, err := security_context.New(
		security_context.WithEvent(map[string]any{
			"test": "test",
		}),
		security_context.WithVariable("key", "value"),
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Plugin action: %s\n", greeter.Execute(ctx))
	
}

var handshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "BASIC_PLUGIN",
	MagicCookieValue: "hello",
}

// pluginMap is the map of plugins we can dispense.
var pluginMap = map[string]plugin.Plugin{
	"greeter": &action.ActionPlugin{},
}