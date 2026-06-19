#!/usr/bin/env bash
set -eo pipefail

FAILED=0
BIN="./lem-in"
AUDIT_DIR="audit"

echo "[BUILD]"
go build -o lem-in ./cmd && echo "PASS" || { echo "FAIL"; exit 1; }

# Helpers
get_end_room() {
    grep -E "^##end" -A1 "$1" | tail -n1 | awk '{print $1}'
}

count_turns() {
    echo "$1" | grep "^L" | wc -l | tr -d ' '
}

all_moved_ants() {
    echo "$1" | grep -oE "L[0-9]+" | sed 's/L//' | sort -n | uniq
}

ants_reached_end() {
    endroom="$2"
    echo "$1" | grep -oE "L[0-9]+-$endroom" | sed "s/-$endroom//" | sed 's/L//' | sort -n | uniq
}

deterministic() {
    out1=$($BIN "$1")
    out2=$($BIN "$1")
    [[ "$(echo "$out1")" == "$(echo "$out2")" ]]
}

check_packages() {
    imports=$(go list -f '{{ join .Imports "\n" }}' ./cmd)
    while IFS= read -r pkg; do
        [[ "$pkg" != *.* ]] && continue
        go list -f '{{.Standard}}' "$pkg" 2>/dev/null | grep true >/dev/null && continue
        return 1
    done <<< "$imports"
    return 0
}

# Turn limits only for examples 00–05
declare -A MAX_TURNS=(
    ["example00.txt"]=6
    ["example01.txt"]=8
    ["example02.txt"]=11
    ["example03.txt"]=6
    ["example04.txt"]=6
    ["example05.txt"]=8
)

echo ""
echo "================ AUDIT ================"

for file in $AUDIT_DIR/*.txt; do
    name=$(basename "$file")
    echo ""
    echo "[$name]"

    set +e
    output=$($BIN "$file" 2>&1)
    code=$?
    set -e

    # BAD EXAMPLES
    if [[ "$name" =~ ^bad ]]; then
        if [ $code -ne 0 ]; then
            echo "PASS (error as expected)"
        else
            echo "FAIL (bad example should error)"
            FAILED=$((FAILED+1))
        fi
        continue
    fi

    # GOOD EXAMPLES
    if [ $code -eq 0 ]; then
        echo "PASS (exit code)"
    else
        echo "FAIL (exit code)"
        FAILED=$((FAILED+1))
        continue
    fi

    # Determinism
    if deterministic "$file"; then
        echo "PASS (deterministic)"
    else
        echo "FAIL (deterministic)"
        FAILED=$((FAILED+1))
    fi

    # Movement format
    if echo "$output" | grep -E "^L[0-9]+-[^[:space:]]+" >/dev/null; then
        echo "PASS (format)"
    else
        echo "FAIL (format)"
        FAILED=$((FAILED+1))
    fi

    # Turn count
    turns=$(count_turns "$output")
    max=${MAX_TURNS[$name]}

    if [ -n "$max" ]; then
        if [ "$turns" -le "$max" ]; then
            echo "PASS (turns $turns <= $max)"
        else
            echo "FAIL (turns $turns > $max)"
            FAILED=$((FAILED+1))
        fi
    else
        echo "PASS (no turn limit)"
    fi

    # End room detection
    endroom=$(get_end_room "$file")

    # Ant completion
    moved=$(all_moved_ants "$output")
    ended=$(ants_reached_end "$output" "$endroom")

    diff=$(comm -3 <(echo "$moved") <(echo "$ended") | wc -l | tr -d ' ')
    if [ "$diff" -eq 0 ]; then
        echo "PASS (all ants reached $endroom)"
    else
        echo "FAIL (ants missing)"
        FAILED=$((FAILED+1))
    fi
done

echo ""
echo "================ PERFORMANCE ================"

echo -n "example06: "
time $BIN audit/example06.txt >/dev/null

echo -n "example07: "
time $BIN audit/example07.txt >/dev/null

echo ""
echo "================ PACKAGES ================"

if check_packages; then
    echo "PASS (std only)"
else
    echo "FAIL (non-std imports)"
    FAILED=$((FAILED+1))
fi

echo ""
echo "================ FINAL RESULT ================"

if [ "$FAILED" -gt 0 ]; then
    echo -e "\e[31mOVERALL: FAIL ($FAILED failed tests)\e[0m"
    exit 1
else
    echo -e "\e[32mOVERALL: PASS (all tests passed)\e[0m"
    exit 0
fi
