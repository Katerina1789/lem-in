/*
Package pathfinder implements the BFS algorithm to find the shortest path in a flow graph.
*/

package pathfinder

import "lem-in/pkg/models"

func FindShortestPath(g *models.FlowGraph, startName, endName string) []*models.FlowEdge {
	// retrieves the start and end nodes
	startNode := g.Nodes[startName]
	endNode := g.Nodes[endName]
	if startNode == nil || endNode == nil {
		return nil
	}

	queue := []*models.FlowNode{startNode}
	visited := map[*models.FlowNode]bool{startNode: true}
	parentEdge := map[*models.FlowNode]*models.FlowEdge{}

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]

		if curr == endNode {
			break
		}

		for _, edge := range curr.Edges {
			residualCap := edge.Capacity - edge.Flow
			if residualCap > 0 && !visited[edge.To] {
				visited[edge.To] = true
				parentEdge[edge.To] = edge
				queue = append(queue, edge.To)
			}
		}
	}

	if !visited[endNode] {
		return nil
	}

	var path []*models.FlowEdge
	curr := endNode
	for curr != startNode {
		edge := parentEdge[curr]
		path = append([]*models.FlowEdge{edge}, path...)
		curr = edge.From
	}
	return path
}
