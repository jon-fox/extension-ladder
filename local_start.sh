#!/bin/bash
set -e

# Prerequisites: go, Chrome/Chromium (for headless browser fallback)
# Install: brew install go && brew install --cask google-chrome

echo "=== Running Tests ==="
go test ./... -count=1

echo ""
echo "=== Building ==="
go build -o extension-ladder cmd/main.go

# Kill any process running on port 8080
lsof -ti :8080 | xargs kill -9 2>/dev/null && echo "Killed existing process on port 8080" || echo "No process on port 8080"

echo "=== Starting Extension Ladder ==="
./extension-ladder -r "ruleset.yaml;rulesets/"
