package pathfinder

import "lem-in/pkg/models"

func FindAllDisjointPaths(g *models.FlowGraph) [][]string {
	for {
		path := FindShortestPath(g, g.Start.Name, g.End.Name)
		if path == nil {
			break
		}
		for _, edge := range path {
			edge.Flow++
			edge.Reverse.Flow--
		}
	}

	var paths [][]string
	for _, edge := range g.Start.Edges {
		if edge.Flow > 0 {
			path := []string{g.Start.Original}
			curr := edge.To
			for curr != g.End {
				if curr.Original != "" && curr.IsIn {
					path = append(path, curr.Original)
				}
				var nextNode *models.FlowNode
				for _, nextEdge := range curr.Edges {
					if nextEdge.Flow > 0 {
						nextNode = nextEdge.To
						break
					}
				}
				if nextNode == nil {
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
