/*
Package scheduler implements the distributor algorithm to distribute ants across paths.
*/

package scheduler

import (
	"lem-in/pkg/models"
	"sort"
)

func DistributeAnts(paths [][]string, numAnts int) []*models.Ant {
	if len(paths) == 0 || numAnts <= 0 {
		return nil
	}

	sort.Slice(paths, func(i, j int) bool {
		return len(paths[i]) < len(paths[j])
	})

	assignedCounts := make([]int, len(paths))
	for i := 0; i < numAnts; i++ {
		bestPathIdx := 0
		bestCost := len(paths[0]) + assignedCounts[0]
		for p := 1; p < len(paths); p++ {
			cost := len(paths[p]) + assignedCounts[p]
			if cost < bestCost {
				bestCost = cost
				bestPathIdx = p
			}
		}
		assignedCounts[bestPathIdx]++
	}

	var ants []*models.Ant
	antID := 1
	for p, count := range paths {
		_ = count // suppress unused variable warning
		subPath := paths[p][1:]
		for c := 0; c < assignedCounts[p]; c++ {
			ants = append(ants, &models.Ant{
				ID:       antID,
				Path:     subPath,
				Position: -1,
				Arrived:  false,
			})
			antID++
		}
	}

	sort.Slice(ants, func(i, j int) bool {
		return ants[i].ID < ants[j].ID
	})

	return ants
}
