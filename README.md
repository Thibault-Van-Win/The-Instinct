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

Currently, all configuration is done by using configuration files. For this purpose, a YAML rule loader was implemented. A config file should contain a list of so called **Reflexes** which consists of:

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

## Dependencies

- Go v1.22
