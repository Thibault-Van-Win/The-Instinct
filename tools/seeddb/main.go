package main

import (
	"context"
	"log"
	"os"

	"gopkg.in/yaml.v3"

	"github.com/Thibault-Van-Win/The-Instinct/internal/config"
	"github.com/Thibault-Van-Win/The-Instinct/internal/factory"
	"github.com/Thibault-Van-Win/The-Instinct/pkg/action"
	"github.com/Thibault-Van-Win/The-Instinct/pkg/reflex"
	"github.com/Thibault-Van-Win/The-Instinct/pkg/rule"
)

var (
	configDir = "./config/security_reflexes.yaml"
)

func main() {
	conf, err := config.Instance()
	if err != nil {
		log.Fatalf("failed to retrieve config: %v", err)
	}

	log.Println("Seeding MongoDB database")

	// Load the yaml rules
	data, err := os.ReadFile(configDir)
	if err != nil {
		log.Fatalf("Failed to read the config file %s: %v", configDir, err)
	}

	// Parse the YAML
	var configs []reflex.ReflexConfig
	if err := yaml.Unmarshal(data, &configs); err != nil {
		log.Fatalf("Failed to parse YAML in file %s: %v", configDir, err)
	}

	// Create repo, add all items	
	// Create the rule registry
	ruleRegistry := rule.NewRuleRegistry(
		rule.WithStandardRules(),
	)

	// Create action registry
	actionRegistry := action.NewActionRegistry(
		action.WithStandardActions(),
	)

	// Initialize repository and service (dependency injection)
	repository, err := factory.NewReflexRepository(&conf.DbConfig, ruleRegistry, actionRegistry)
	if err != nil {
		log.Fatalf("Failed to create reflex repository: %v", err)
	}
	service := reflex.NewReflexService(repository)
	defer service.Close(context.Background())

	for _, reflexConfig := range configs {
		_, err := service.CreateReflex(context.Background(), reflexConfig)
		if err != nil {
			log.Fatalf("Failed to create reflex %s: %v", reflexConfig.Name, err)
		}
	}

	log.Println("Database seeding completed successfully")
}
