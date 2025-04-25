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
	// Create action registry
	actionRegistry := action.NewActionRegistry()
	actionRegistry.RegisterStandardActions()

	// Create the rule registry
	ruleRegistry := rule.NewRuleRegistry()
	ruleRegistry.RegisterStandardRules()

	// Create a new instinct system
	system := instinct.New(ruleRegistry, actionRegistry)

	// Load reflexes from YAML files
	if err := system.LoadReflexes(loaders.YAMLLoader, map[string]any{
		"directory": "./config",
	}); err != nil {
		log.Fatalf("Failed to load reflexes from YAML: %v", err)
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
