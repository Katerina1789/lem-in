package pathfinder

import "lem-in/pkg/models"

// FindShortestPath returns the shortest augmenting path from startName to endName in the residual
// flow graph using BFS, or nil if no such path exists.
func FindShortestPath(g *models.FlowGraph, startName, endName string) []*models.FlowEdge {
	startNode := g.Nodes[startName]
	endNode := g.Nodes[endName]
	if startNode == nil || endNode == nil {
		return nil
	}

	queue := []*models.FlowNode{startNode}
	visited := map[*models.FlowNode]bool{startNode: true} // tracks nodes already enqueued
	parentEdge := map[*models.FlowNode]*models.FlowEdge{} // records which edge led to each node

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]

		if curr == endNode {
			break
		}

		for _, edge := range curr.Edges {
			residualCap := edge.Capacity - edge.Flow // remaining capacity (also works for reverse edges)
			if residualCap > 0 && !visited[edge.To] {
				visited[edge.To] = true
				parentEdge[edge.To] = edge // remember how we reached edge.To
				queue = append(queue, edge.To)
			}
		}
	}

	if !visited[endNode] { // BFS exhausted without reaching end — no augmenting path
		return nil
	}

	// reconstruct path by walking parentEdge backwards from end to start
	var path []*models.FlowEdge
	curr := endNode
	for curr != startNode {
		edge := parentEdge[curr]
		path = append([]*models.FlowEdge{edge}, path...) // prepend to keep start→end order
		curr = edge.From
	}
	return path
}
