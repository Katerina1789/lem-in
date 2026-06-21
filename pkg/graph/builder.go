package graph

import "lem-in/pkg/models"

// BuildAdjacency fills Graph.Adjacency with bidirectional neighbor lists derived from Graph.Links.
func BuildAdjacency(g *models.Graph) {
	g.Adjacency = make(map[string][]string)

	for _, link := range g.Links {
		g.Adjacency[link.Room1] = append(g.Adjacency[link.Room1], link.Room2)
		g.Adjacency[link.Room2] = append(g.Adjacency[link.Room2], link.Room1)
	}
}

// IsConnected reports whether end is reachable from start via BFS on Graph.Adjacency.
func IsConnected(g *models.Graph, start, end string) bool {
	visited := map[string]bool{start: true}
	queue := []string{start}

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]

		if curr == end {
			return true
		}

		for _, neighbor := range g.Adjacency[curr] {
			if !visited[neighbor] {
				visited[neighbor] = true
				queue = append(queue, neighbor)
			}
		}
	}

	return false
}
