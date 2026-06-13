package scheduler

import (
	"lem-in/pkg/models"
	"testing"
)

func TestSimulate(t *testing.T) {
	paths := [][]string{
		{"start", "A", "end"},
		{"start", "B", "C", "end"},
	}

	ants := DistributeAnts(paths, 3)
	turns := Simulate(ants, &models.Graph{})

	if len(turns) == 0 {
		t.Fatalf("expected turns simulation, got 0 turns")
	}

	// Verify all ants reached end
	for _, ant := range ants {
		if !ant.Arrived {
			t.Errorf("ant %d did not arrive", ant.ID)
		}
	}
}
