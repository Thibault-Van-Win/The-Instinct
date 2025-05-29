package main

import (
	"fmt"
	"log"

	"github.com/Thibault-Van-Win/The-Instinct/pkg/security_context"
)

func (th *TheHive) Execute(ctx *security_context.SecurityContext) error {
	cfg := ctx.Variables["config"].(map[string]any)
	log.Printf("Hello from the TheHive plugin, here is my config %v", cfg)
	log.Printf("Sending something to %s\n", th.HiveUrl)
	return nil
}

func (th *TheHive) GetType() string {
	return "thehive"
}

func (th *TheHive) GetName() string {
	return "A plugin for some Hive instance"
}

func (th *TheHive) Validate() error {
	if th.HiveSkipTls {
		log.Println("Warning: TLS verification is disabled for The Hive plugin")
	}

	if th.HiveUrl == "" {
		return fmt.Errorf("TheHive plugin should have a non-empty URL, please set the THE_HIVE_URL environment variable")
	}

	if th.HiveApiKey == "" {
		return fmt.Errorf("TheHive plugin should have an API key, please set the THE_HIVE_API_KEY environment variable")
	}

	return nil
}
