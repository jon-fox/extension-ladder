#!/bin/bash

# Kill any process running on port 8080
lsof -ti :8080 | xargs kill -9 2>/dev/null && echo "Killed existing process on port 8080" || echo "No process on port 8080"

# Start the application
go run cmd/main.go -r ruleset.yaml
