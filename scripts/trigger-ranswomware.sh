#! /bin/sh

curl -X POST \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Ransomware detected",
    "severity": "high",
    "tags": ["ransomware", "urgent"]
  }' \
  http://localhost:8080/event