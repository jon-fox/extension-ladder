#!/bin/bash
set -e

case "${1:-tidy}" in
  tidy)
    echo "=== Tidying ==="
    go mod tidy
    go clean
    ;;
  lint)
    echo "=== Tidying ==="
    go mod tidy
    go clean
    echo "=== Linting ==="
    gofumpt -l -w .
    golangci-lint run --fix
    ;;
  install-linters)
    echo "=== Installing Linters ==="
    go install mvdan.cc/gofumpt@latest
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
    ;;
  *)
    echo "Usage: ./maintenance.sh [tidy|lint|install-linters]"
    exit 1
    ;;
esac

echo "Done."
