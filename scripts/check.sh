#!/usr/bin/env sh
set -u

FAILED=0

echo "Running go fmt check..."
UNFORMATTED=$(gofmt -l .)
if [ -n "$UNFORMATTED" ]; then
    echo "Unformatted files:"
    echo "$UNFORMATTED"
    FAILED=$((FAILED+1))
else
    echo "All files formatted."
fi

echo "Running go vet..."
if ! go vet ./...; then
    echo "go vet reported issues."
    FAILED=$((FAILED+1))
else
    echo "go vet passed."
fi

echo "Running go test..."
if ! go test ./...; then
    echo "Tests failed."
    FAILED=$((FAILED+1))
else
    echo "Tests passed."
fi

echo "Checks complete!"

# Final exit code
if [ "$FAILED" -gt 0 ]; then
    exit 1
else
    exit 0
fi
