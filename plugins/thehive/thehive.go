// Package for a TheHive plugin
//
// Accepts Action calls where the config can be found under ctx.Variables["config"]
// This should contain a sub action alongside some params for it
// Credentials for TheHive should be placed in the environment
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Thibault-Van-Win/thehive4go"
	"github.com/hashicorp/go-plugin"
	"github.com/sethvargo/go-envconfig"

	"github.com/Thibault-Van-Win/The-Instinct/pkg/action"
)

type TheHive struct {
	HiveUrl     string `env:"THE_HIVE_URL"`
	HiveApiKey  string `env:"THE_HIVE_API_KEY"`
	HiveSkipTls bool   `env:"THE_HIVE_SKIP_TLS"`

	apiClient *thehive4go.APIClient
}

func newFromEnv() (*TheHive, error) {
	var th TheHive

	if err := envconfig.Process(context.Background(), &th); err != nil {
		return nil, fmt.Errorf("failed to process env variables: %v", err)
	}

	th.apiClient = thehive4go.NewAPIClient(thehive4go.Config{
		SkipTLSVerification: th.HiveSkipTls,
		URL: th.HiveUrl,
		APIKey: th.HiveApiKey,
	})

	return &th, nil
}

var handshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "BASIC_PLUGIN",
	MagicCookieValue: "hello",
}

func main() {
	instance, err := newFromEnv()
	if err != nil {
		log.Fatalf("Failed to create thehive plugin from env: %v", err)
	}

	var pluginMap = map[string]plugin.Plugin{
		"thehive": &action.ActionPlugin{Impl: instance},
	}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
	})
}
