package graph

import "lem-in/pkg/models"

// BuildAdjacency fills Graph.Adjacency with bidirectional neighbor lists from Graph.Links
func BuildAdjacency(g *models.Graph) {
	g.Adjacency = make(map[string][]string)

	// each link A-B: add B as neighbor of A, and A as neighbor of B
	for _, link := range g.Links {
		g.Adjacency[link.Room1] = append(g.Adjacency[link.Room1], link.Room2)
		g.Adjacency[link.Room2] = append(g.Adjacency[link.Room2], link.Room1)
	}
}

// IsConnected returns true if end is reachable from start using BFS on Graph.Adjacency
func IsConnected(g *models.Graph, start, end string) bool {
	// mark start as visited and add it to the queue
	visited := map[string]bool{start: true}
	queue := []string{start}

	for len(queue) > 0 {
		// dequeue front
		curr := queue[0]
		queue = queue[1:]

		// reached end - a valid path exists
		if curr == end {
			return true
		}

		// enqueue each unvisited neighbor
		for _, neighbor := range g.Adjacency[curr] {
			if !visited[neighbor] {
				visited[neighbor] = true
				queue = append(queue, neighbor)
			}
		}
	}

	// queue exhausted without finding end
	return false
}
