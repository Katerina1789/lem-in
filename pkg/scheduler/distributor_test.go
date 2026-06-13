package scheduler

import (
	"testing"
)

func TestDistributeAnts(t *testing.T) {
	paths := [][]string{
		{"start", "A", "end"},
		{"start", "B", "C", "end"},
	}

	ants := DistributeAnts(paths, 5)
	if len(ants) != 5 {
		t.Fatalf("expected 5 ants, got %d", len(ants))
	}

	// Verify ant 1 has path excluding start
	if len(ants[0].Path) != 2 || ants[0].Path[0] != "A" || ants[0].Path[1] != "end" {
		t.Errorf("expected path without start, got %v", ants[0].Path)
	}
}
