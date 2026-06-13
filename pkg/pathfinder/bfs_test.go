package pathfinder

import (
	"testing"
)

func TestFindShortestPath(t *testing.T) {
	g := buildTestFlowGraph(
		[]string{"start", "A", "end"},
		[][]string{{"start", "A"}, {"A", "end"}},
		"start",
		"end",
	)

	path := FindShortestPath(g, "start_out", "end_in")
	if len(path) != 3 { // start_out -> A_in, A_in -> A_out, A_out -> end_in
		t.Fatalf("expected path length 3, got %d", len(path))
	}

	if path[0].From.Name != "start_out" || path[2].To.Name != "end_in" {
		t.Errorf("invalid path sequence")
	}
}
