#!/usr/bin/env sh
set -eu

echo "Running go vet..."
go vet ./... || true

echo "Running go test..."
go test ./...

echo "Checks complete."
