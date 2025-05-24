package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/Thibault-Van-Win/The-Instinct/pkg/action"
	"github.com/Thibault-Van-Win/The-Instinct/pkg/instinct"
	"github.com/Thibault-Van-Win/The-Instinct/pkg/loaders"
	"github.com/Thibault-Van-Win/The-Instinct/pkg/rule"
)

func main() {
	// Create the rule registry
	ruleRegistry := rule.NewRuleRegistry(
		rule.WithStandardRules(),
	)

	// Create action registry
	actionRegistry := action.NewActionRegistry(
		action.WithStandardActions(),
		action.WithPlugins(),
	)

	defer actionRegistry.Close()

	// Create a new instinct system
	system := instinct.New(ruleRegistry, actionRegistry)

	loader, err := loaders.NewLoaderFactory(ruleRegistry, actionRegistry).CreateLoader(loaders.YAML, map[string]any{
		"directory": "./config",
	})
	if err != nil {
		log.Fatalf("Failed to create reflex loader: %v", err)
	}

	// Load reflexes from YAML files
	if err := system.LoadReflexes(loader); err != nil {
		log.Fatalf("Failed to load reflexes: %v", err)
	}

	// Process an incoming alert
	alertJson := []byte(`{
		"title": "Ransomware detected",
		"severity": "high",
		"tags": ["ransomware", "urgent"]
	}`)

	var alert map[string]any
	if err := json.Unmarshal(alertJson, &alert); err != nil {
		log.Fatal(err)
	}

	// Process the event
	if err := system.ProcessEvent(alert); err != nil {
		fmt.Fprintf(os.Stderr, "Error processing event: %v\n", err)
		os.Exit(1)
	}

	// Process an incoming alert
	alertJson = []byte(`{
		"title": "Someone logged in",
		"severity": "medium",
		"tags": ["cloud", "suspicious-login"]
	}`)

	if err := json.Unmarshal(alertJson, &alert); err != nil {
		log.Fatal(err)
	}

	if err := system.ProcessEvent(alert); err != nil {
		fmt.Fprintf(os.Stderr, "Error processing event: %v\n", err)
		os.Exit(1)
	}
}
