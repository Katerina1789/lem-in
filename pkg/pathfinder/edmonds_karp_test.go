package pathfinder

import (
	"testing"
)

func TestFindAllDisjointPaths(t *testing.T) {
	g := buildTestFlowGraph(
		[]string{"start", "A", "B", "end"},
		[][]string{
			{"start", "A"}, {"A", "end"},
			{"start", "B"}, {"B", "end"},
		},
		"start",
		"end",
	)

	paths := FindAllDisjointPaths(g)
	if len(paths) != 2 {
		t.Fatalf("expected 2 disjoint paths, got %d", len(paths))
	}

	t.Logf("Paths found: %v", paths)
}
