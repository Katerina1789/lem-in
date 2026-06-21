package graph

import "lem-in/pkg/models"

// SplitForFlow converts a validated Graph into a FlowGraph using node splitting.
// Each internal room R becomes R_in and R_out connected by a capacity-1 edge, which
// lets the max-flow algorithm enforce the one-ant-per-room constraint automatically.
func SplitForFlow(g *models.Graph) *models.FlowGraph {
	nodes := make(map[string]*models.FlowNode)

	for name, room := range g.Rooms {
		switch {
		case room.IsStart:
			nodes[name+"_out"] = &models.FlowNode{Name: name + "_out", Original: name, IsOut: true} // ants only leave start
		case room.IsEnd:
			nodes[name+"_in"] = &models.FlowNode{Name: name + "_in", Original: name, IsIn: true} // ants only arrive at end
		default:
			inNode := &models.FlowNode{Name: name + "_in", Original: name, IsIn: true}
			outNode := &models.FlowNode{Name: name + "_out", Original: name, IsOut: true}
			nodes[name+"_in"] = inNode
			nodes[name+"_out"] = outNode
			addEdgePair(inNode, outNode, 1) // capacity-1 internal edge enforces one-ant-per-room
		}
	}

	for _, link := range g.Links {
		r1Out := nodes[link.Room1+"_out"]
		r2In := nodes[link.Room2+"_in"]
		r2Out := nodes[link.Room2+"_out"]
		r1In := nodes[link.Room1+"_in"]

		if r1Out != nil && r2In != nil {
			addEdgePair(r1Out, r2In, 1) // tunnel: Room1 → Room2
		}
		if r2Out != nil && r1In != nil {
			addEdgePair(r2Out, r1In, 1) // tunnel: Room2 → Room1 (tunnels are undirected)
		}
	}

	var startName, endName string
	for name, room := range g.Rooms {
		if room.IsStart {
			startName = name
		}
		if room.IsEnd {
			endName = name
		}
	}

	return &models.FlowGraph{
		Nodes: nodes,
		Start: nodes[startName+"_out"],
		End:   nodes[endName+"_in"],
	}
}

// addEdgePair creates a forward edge (capacity cap) and its residual reverse edge (capacity 0).
// The reverse edge lets Edmonds-Karp cancel previously committed flow when rerouting paths.
func addEdgePair(from, to *models.FlowNode, cap int) {
	forward := &models.FlowEdge{From: from, To: to, Capacity: cap, Flow: 0}
	reverse := &models.FlowEdge{From: to, To: from, Capacity: 0, Flow: 0, Reverse: forward} // residual only
	forward.Reverse = reverse

	from.Edges = append(from.Edges, forward)
	to.Edges = append(to.Edges, reverse) // reverse on to-node so BFS can traverse it backwards
}
