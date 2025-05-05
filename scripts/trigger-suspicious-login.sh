#! /bin/sh

SEVERITY="${1:-medium}"

curl -X POST \
  -H "Content-Type: application/json" \
  -d "{
    \"title\": \"Someone logged in\",
    \"severity\": \"$SEVERITY\",
    \"tags\": [\"cloud\", \"suspicious-login\"]
  }" \
  http://localhost:8080/event
