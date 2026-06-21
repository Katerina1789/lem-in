package pathfinder

import "lem-in/pkg/models"

// FindAllDisjointPaths runs Edmonds-Karp max-flow on g and returns all vertex-disjoint paths
// from start to end as slices of original room names.
func FindAllDisjointPaths(g *models.FlowGraph) [][]string {
	// augmentation phase: keep finding shortest augmenting paths until none remain
	for {
		path := FindShortestPath(g, g.Start.Name, g.End.Name)
		if path == nil {
			break
		}
		for _, edge := range path {
			edge.Flow++         // commit one unit of flow on the forward edge
			edge.Reverse.Flow-- // decrease the reverse edge so future BFS can cancel this flow if needed
		}
	}

	// extraction phase: each start edge with flow=1 is the beginning of a distinct path
	var paths [][]string
	for _, edge := range g.Start.Edges {
		if edge.Flow > 0 {
			path := []string{g.Start.Original}
			curr := edge.To
			for curr != g.End {
				if curr.Original != "" && curr.IsIn { // collect room name only at _in nodes to avoid duplicates from the R_in→R_out split
					path = append(path, curr.Original)
				}
				var nextNode *models.FlowNode
				for _, nextEdge := range curr.Edges {
					if nextEdge.Flow > 0 { // follow the edge that carries flow
						nextNode = nextEdge.To
						break
					}
				}
				if nextNode == nil { // broken path (should not happen on a valid flow)
					break
				}
				curr = nextNode
			}
			path = append(path, g.End.Original)
			paths = append(paths, path)
		}
	}
	return paths
}
