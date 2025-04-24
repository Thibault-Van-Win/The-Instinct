#! /bin/sh

curl -X POST \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Someone logged in",
		"severity": "medium",
		"tags": ["cloud", "suspicious-login"]
  }' \
  http://localhost:8080/event