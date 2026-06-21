package scheduler

import (
	"lem-in/pkg/models"
	"sort"
)

// DistributeAnts assigns each of the numAnts ants to one of the given paths using a
// greedy min-cost strategy: each ant goes to the path with the lowest current cost
// (path length + number of ants already assigned to it).
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
	for i := 1; i <= numAnts; i++ {
		ants = append(ants, &models.Ant{
			ID:       i,
			Position: -1,
			Arrived:  false,
		})
	}

	antIdx := 0
	for antIdx < numAnts {
		for p := range paths {
			if assignedCounts[p] > 0 {
				ants[antIdx].Path = paths[p][1:] // store path without the start room
				assignedCounts[p]--
				antIdx++
				if antIdx >= numAnts {
					break
				}
			}
		}
	}

	return ants
}
