#!/usr/bin/env bash
# run_audit.sh - runs the official lem-in audit test cases

set -eo pipefail

echo "=========================================="
echo " Building lem-in executable..."
echo "=========================================="
go build -o lem-in ./cmd
if [ $? -ne 0 ]; then
    echo "Error: Failed to build lem-in"
    exit 1
fi
echo "Build successful."
echo ""

echo "=========================================="
echo " Running Audit Tests..."
echo "=========================================="

PASSED=0
FAILED=0

# Loop through all files in audit directory
for file in audit/*.txt; do
    filename=$(basename "$file")
    echo "------------------------------------------"
    echo "Testing: $filename"
    echo "------------------------------------------"
    
    # Run the executable and capture output and exit code
    # We disable exit on error temporarily to capture the non-zero status
    set +e
    output=$(./lem-in "$file" 2>&1)
    exit_code=$?
    set -e

    if [[ "$filename" =~ ^bad ]]; then
        # Bad examples: should exit with non-zero code and print error
        if [ $exit_code -ne 0 ]; then
            echo "RESULT: PASS (Returned error as expected)"
            echo "Output: $output"
            PASSED=$((PASSED+1))
        else
            echo "RESULT: FAIL (Expected error but exit code was 0)"
            FAILED=$((FAILED+1))
        fi
    else
        # Good examples: should exit with 0 and output moves
        if [ $exit_code -eq 0 ]; then
            echo "RESULT: PASS"
            echo "Output:"
            echo "$output"
            PASSED=$((PASSED+1))
        else
            echo "RESULT: FAIL (Exit code $exit_code)"
            echo "Output: $output"
            FAILED=$((FAILED+1))
        fi
    fi
    echo ""
done

echo "=========================================="
echo " Audit Summary"
echo "=========================================="
echo "Total Passed: $PASSED"
echo "Total Failed: $FAILED"

if [ $FAILED -ne 0 ]; then
    exit 1
fi
