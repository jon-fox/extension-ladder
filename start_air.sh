#!/bin/bash
set -e

GOBIN=$(go env GOPATH)/bin

# Install Air if not present
if ! command -v air &> /dev/null && [ ! -f "$GOBIN/air" ]; then
    echo "=== Installing Air ==="
    go install github.com/air-verse/air@latest
fi

# Kill any process running on port 8080
lsof -ti :8080 | xargs kill -9 2>/dev/null && echo "Killed existing process on port 8080" || echo "No process on port 8080"

echo "=== Starting Air (hot-reload) ==="
"$GOBIN/air"
