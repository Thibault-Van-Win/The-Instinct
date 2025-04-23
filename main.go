package main

import (
	"encoding/json"
	"log"

	"github.com/Thibault-Van-Win/The-Instinct/pkg/action"
	"github.com/Thibault-Van-Win/The-Instinct/pkg/reflex"
	"github.com/Thibault-Van-Win/The-Instinct/pkg/rule"
)

func main() {
	// Parse a reflex
	reflex := &reflex.Reflex{
		Name: "test",	
		Rule: &rule.CelRule{},
		Actions: []action.Action{action.PrintAction()},
	}

	// Mock some input data
	alertJson := []byte(`{
		"title": "Ransomware detected",
		"severity": "high",
		"tags": ["ransomware", "urgent"]
	}`)

	var alert map[string]any
	if err := json.Unmarshal(alertJson, &alert); err != nil {
		log.Fatal(err)
	}

	// Check if the Reflex is needed
	ok, _ := reflex.Match(alert)
	if  ok {
		reflex.Do()	
	}
}
