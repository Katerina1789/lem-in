package pathfinder

import "sort"

// SelectOptimalPaths picks the subset of allPaths that minimises the number of turns
// needed to move antCount ants from start to end.
func SelectOptimalPaths(allPaths [][]string, antCount int) [][]string {
	if len(allPaths) == 0 || antCount <= 0 {
		return nil
	}

	sort.Slice(allPaths, func(i, j int) bool {
		return len(allPaths[i]) < len(allPaths[j])
	})

	var bestPaths [][]string
	bestTurns := -1

	for k := 1; k <= len(allPaths); k++ {
		selected := allPaths[:k]
		turns := computeTurns(selected, antCount)
		if bestTurns == -1 || turns < bestTurns {
			bestTurns = turns
			bestPaths = selected
		}
	}

	return bestPaths
}

// computeTurns returns the minimum number of turns needed to move antCount ants
// through the given set of paths using binary search.
func computeTurns(paths [][]string, antCount int) int {
	lengths := make([]int, len(paths))
	maxL := 0
	for i, p := range paths {
		lengths[i] = len(p) - 1
		if lengths[i] > maxL {
			maxL = lengths[i]
		}
	}

	low := 0
	high := antCount + maxL

	for low < high {
		mid := (low + high) / 2
		capacity := 0
		for _, L := range lengths {
			if mid >= L {
				capacity += mid - L + 1 // ants that can finish path L within T=mid turns
			}
		}
		if capacity >= antCount {
			high = mid
		} else {
			low = mid + 1
		}
	}

	return low
}
