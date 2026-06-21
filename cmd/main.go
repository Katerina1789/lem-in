package main

import (
	"fmt"
	"os"

	"lem-in/pkg/graph"
	"lem-in/pkg/output"
	"lem-in/pkg/parser"
	"lem-in/pkg/pathfinder"
	"lem-in/pkg/scheduler"
	"lem-in/pkg/validator"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: lem-in <input-file>")
		os.Exit(1)
	}

	path := os.Args[1]

	raw, err := parser.ReadFile(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, "ERROR: invalid data format")
		os.Exit(1)
	}

	lines := parser.SplitLines(raw)
	lines = parser.FilterComments(lines)
	lines = parser.TrimWhitespace(lines)

	g, antCount, err := validator.Validate(lines)
	if err != nil {
		fmt.Fprintln(os.Stderr, "ERROR: invalid data format")
		os.Exit(1)
	}

	var startName, endName string
	for name, room := range g.Rooms {
		if room.IsStart {
			startName = name
		}
		if room.IsEnd {
			endName = name
		}
	}
	if !graph.IsConnected(g, startName, endName) {
		fmt.Fprintln(os.Stderr, "ERROR: invalid data format")
		os.Exit(1)
	}

	flowGraph := graph.SplitForFlow(g)
	allPaths := pathfinder.FindAllDisjointPaths(flowGraph)
	optimalPaths := pathfinder.SelectOptimalPaths(allPaths, antCount)
	if len(optimalPaths) == 0 {
		fmt.Fprintln(os.Stderr, "ERROR: invalid data format")
		os.Exit(1)
	}

	ants := scheduler.DistributeAnts(optimalPaths, antCount)
	turns := scheduler.Simulate(ants, g)

	fmt.Print(output.FormatFull(raw, turns))
}
