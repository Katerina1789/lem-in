package graph

import "lem-in/pkg/models"

// SplitForFlow converts a validated Graph into a FlowGraph using node splitting
// Node splitting turns the "one ant per room" constraint into a per-edge capacity limit which lets the flow algorithm enforce it automatically during pathfinding.
func SplitForFlow(g *models.Graph) *models.FlowGraph {
	nodes := make(map[string]*models.FlowNode)

	// create flow nodes for every room - start/end get one node, others get two
	for name, room := range g.Rooms {
		switch {
		case room.IsStart:
			// start only needs _out - ants leave but never re-enter start
			nodes[name+"_out"] = &models.FlowNode{Name: name + "_out", Original: name, IsOut: true}

		case room.IsEnd:
			// end only needs _in - ants arrive but we never route out of end
			nodes[name+"_in"] = &models.FlowNode{Name: name + "_in", Original: name, IsIn: true}

		default:
			// internal room: split into R_in and R_out connected by a capacity-1 edge
			inNode := &models.FlowNode{Name: name + "_in", Original: name, IsIn: true}
			outNode := &models.FlowNode{Name: name + "_out", Original: name, IsOut: true}
			nodes[name+"_in"] = inNode
			nodes[name+"_out"] = outNode

			// capacity-1 internal edge enforces one-ant-per-room at the flow level
			addEdgePair(inNode, outNode, 1)
		}
	}

	// add directed tunnel edges for every original link A-B (both directions)
	for _, link := range g.Links {
		r1Out := nodes[link.Room1+"_out"]
		r2In := nodes[link.Room2+"_in"]
		r2Out := nodes[link.Room2+"_out"]
		r1In := nodes[link.Room1+"_in"]

		// A_out -> B_in: ants travel from A to B through this tunnel
		if r1Out != nil && r2In != nil {
			addEdgePair(r1Out, r2In, 1)
		}

		// B_out -> A_in: tunnels are undirected so add the reverse direction too
		if r2Out != nil && r1In != nil {
			addEdgePair(r2Out, r1In, 1)
		}
	}

	// locate start and end room names to wire FlowGraph.Start and FlowGraph.End
	var startName, endName string
	for name, room := range g.Rooms {
		if room.IsStart {
			startName = name
		}
		if room.IsEnd {
			endName = name
		}
	}

	// return the completed flow graph with Start and End pointers set
	return &models.FlowGraph{
		Nodes: nodes,
		Start: nodes[startName+"_out"],
		End:   nodes[endName+"_in"],
	}
}

// addEdgePair creates a forward edge (capacity cap) and a residual reverse edge (capacity 0).
// The reverse edge allows Edmonds-Karp to cancel previously committed flow when rerouting.
func addEdgePair(from, to *models.FlowNode, cap int) {
	// build the forward edge from -> to
	forward := &models.FlowEdge{From: from, To: to, Capacity: cap, Flow: 0}

	// build the reverse edge to -> from with capacity 0 (residual only, not a real tunnel)
	reverse := &models.FlowEdge{From: to, To: from, Capacity: 0, Flow: 0, Reverse: forward}
	forward.Reverse = reverse

	// attach forward to the from-node's outgoing list
	from.Edges = append(from.Edges, forward)

	// attach reverse to the to-node's outgoing list so BFS can traverse it backwards
	to.Edges = append(to.Edges, reverse)
}
