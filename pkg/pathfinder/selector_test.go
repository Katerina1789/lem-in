package pathfinder

import (
	"testing"
)

func TestSelectOptimalPaths(t *testing.T) {
	paths := [][]string{
		{"start", "A", "end"},             // length 2
		{"start", "B", "C", "end"},        // length 3
	}

	// 5 ants.
	// If 1 path: 5 turns.
	// If 2 paths:
	// - path 1 gets 3 ants -> 2 + 3 - 1 = 4 turns.
	// - path 2 gets 2 ants -> 3 + 2 - 1 = 4 turns.
	// Max turns = 4. 2 paths is optimal.
	best := SelectOptimalPaths(paths, 5)
	if len(best) != 2 {
		t.Fatalf("expected 2 paths selected, got %d", len(best))
	}
}
