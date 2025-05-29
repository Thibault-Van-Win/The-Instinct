package main

import (
	"fmt"
	"log"

	"github.com/Thibault-Van-Win/thehive4go/alert"
	"github.com/mitchellh/mapstructure"

	"github.com/Thibault-Van-Win/The-Instinct/pkg/security_context"
)

// Action implementation
func (th *TheHive) Execute(ctx *security_context.SecurityContext) error {
	cfg, ok := ctx.Variables["config"].(map[string]any)
	if !ok {
		return fmt.Errorf("expect config struct as a variable, go %v", cfg)
	}
	log.Printf("Hello from the TheHive plugin, here is my config %v", cfg)

	subaction, ok := cfg["subaction"].(string)
	if !ok {
		return fmt.Errorf("missing or invalid subaction")
	}

	switch subaction {
	case "create":
		return th.handleCreate(ctx)
	case "update":
		return th.handleUpdate(ctx)
	case "delete":
		return th.handleDelete(ctx)
	default:
		return fmt.Errorf("unsupported subaction: %s", subaction)
	}
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

// Subaction implementations
func (th *TheHive) handleCreate(ctx *security_context.SecurityContext) error {
	cfg, ok := ctx.Variables["config"].(map[string]any)
	if !ok {
		return fmt.Errorf("failed to retrieve config, either missing or malformed")
	}

	fields, ok := cfg["fields"].(map[string]any)
	if !ok {
		return fmt.Errorf("failed to retrieve fields for new alert: missing or malformed")
	}

	var createAlertRequest alert.CreateAlertRequest
	err := mapstructure.Decode(fields, createAlertRequest)
	if err != nil {
		return fmt.Errorf("failed to decode fields: %v", err)
	}

	_, err = th.apiClient.Alerts.Create(&createAlertRequest)
	return err
}

func (th *TheHive) handleUpdate(ctx *security_context.SecurityContext) error {
	return fmt.Errorf("update subaction is not implemented yet")
}

func (th *TheHive) handleDelete(ctx *security_context.SecurityContext) error {
	return fmt.Errorf("delete subaction is not implemented yet")
}
