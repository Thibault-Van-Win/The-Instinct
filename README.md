# The Instinct

> _“Guiding the pack with instinct and precision.”_  

## Introduction

Your tireless rule-based assistant within **The Den**. **The Instinct** helps triage alerts, auto-escalates threats, and recommends or takes actions based on predefined logic — freeing up analysts to focus on the unknown.

## How to run

To access the detection system, a webserver is provided which accepts events. To start run the following command:

```sh
go run cmd/webserver/main.go
```

Next, an example event can be created with curl:

```sh
curl -X POST \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Ransomware detected",
    "severity": "high",
    "tags": ["ransomware", "urgent"]
  }' \
  http://localhost:8080/event
```

For testing purposes, some of these commands have received their own script which can be found under `./scripts/`

## Configuration

Currently, there are two options to configure The Instinct, which are for now hardcoded. Each of the reflex config options comes with its own loader:

## YAML

The first option is to use YAML files in combination with the **YAML file loader**. A config file should contain a list of so called **Reflexes** which consists of:

1. Name: The name of the reflex
2. Expression: CEL syntax for when a reflex should lead to action
3. Actions: a list of actions for an expression match on an incoming event

Example:

```yaml
- name: ransomware-detector
  expression: "severity == 'high' && tags.exists(tag, tag == 'ransomware')"
  actions:
    - type: print
      params:
        message: "Ransomware detected! Initiating response..."
```

The following code shows how to load rules from a YAML file:

```go
system := instinct.New(ruleRegistry, actionRegistry)
// Load reflexes from YAML files
if err := system.LoadReflexes(loaders.YAML, map[string]any{
    "directory": "./config",
}); err != nil {
    log.Fatalf("Failed to load reflexes from YAML: %v", err)
}
```

The following config options can be provided:

1. `directory`: directory containing a collection of yaml files with reflex configs. (required)

## MongoDB

The second option uses MongoDB for storing reflex configuration objects. This comes with a **MongoDB loader**. Currently, the webserver has no endpoint for configuring these, but this will be provided in the future. For now, rules can be configure this way by using the `seeddb` utility which can be found under `/tools/seeddb/main.go`. The loader can be used in the following way:

```go
system = instinct.New(ruleRegistry, actionRegistry)
// Load reflexes from MongoDB instance
if err := system.LoadReflexes(loaders.MongoDB, map[string]any{
    "uri": "mongodb://user:secret@localhost:27017",
    "database": "instinct",
    "collection": "reflexes",
}); err != nil {
    log.Fatalf("Failed to load reflexes: %v", err)
}
```

The following configuration options can be provided:

1. `uri`: connection string for a MongoDB instance. E.g., `"mongodb://user:secret@localhost:27017"`. (required)
2. `database`: database to be used for the reflex configurations. (required)
3. `collection`: collection to be used for the reflex configurations. (required)

## Dependencies

- Go v1.22
- MongoDB v8.0.8
