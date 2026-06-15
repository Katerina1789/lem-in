package main

import (
	"fmt"
	"os"

	"lem-in/pkg/parser"
	"lem-in/pkg/validator"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: lem-in <input-file>")
		os.Exit(1)
	}

	path := os.Args[1]

	// Phase 1: Parse file
	raw, err := parser.ReadFile(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, "ERROR: invalid data format")
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	lines := parser.SplitLines(raw)
	lines = parser.FilterComments(lines)
	lines = parser.TrimWhitespace(lines)

	// Phase 2: Validate
	graph, antCount, err := validator.Validate(lines)
	if err != nil {
		fmt.Fprintln(os.Stderr, "ERROR: invalid data format")
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// For now, just acknowledge success (Phase 3+ will build on this)
	_ = graph
	_ = antCount
}
