package pathfinder

import "lem-in/pkg/models"

func buildTestFlowGraph(rooms []string, links [][]string, start, end string) *models.FlowGraph {
	nodes := make(map[string]*models.FlowNode)
	
	for _, room := range rooms {
		if room == start {
			nodes[room+"_out"] = &models.FlowNode{Name: room + "_out", Original: room, IsOut: true}
		} else if room == end {
			nodes[room+"_in"] = &models.FlowNode{Name: room + "_in", Original: room, IsIn: true}
		} else {
			nodes[room+"_in"] = &models.FlowNode{Name: room + "_in", Original: room, IsIn: true}
			nodes[room+"_out"] = &models.FlowNode{Name: room + "_out", Original: room, IsOut: true}
			
			inNode := nodes[room+"_in"]
			outNode := nodes[room+"_out"]
			edge := &models.FlowEdge{From: inNode, To: outNode, Capacity: 1}
			revEdge := &models.FlowEdge{From: outNode, To: inNode, Capacity: 0, Reverse: edge}
			edge.Reverse = revEdge
			inNode.Edges = append(inNode.Edges, edge)
			outNode.Edges = append(outNode.Edges, revEdge)
		}
	}
	
	for _, link := range links {
		r1, r2 := link[0], link[1]
		
		var out1, in1, out2, in2 *models.FlowNode
		if r1 == start {
			out1 = nodes[r1+"_out"]
		} else if r1 == end {
			in1 = nodes[r1+"_in"]
		} else {
			in1 = nodes[r1+"_in"]
			out1 = nodes[r1+"_out"]
		}
		
		if r2 == start {
			out2 = nodes[r2+"_out"]
		} else if r2 == end {
			in2 = nodes[r2+"_in"]
		} else {
			in2 = nodes[r2+"_in"]
			out2 = nodes[r2+"_out"]
		}
		
		if out1 != nil && in2 != nil {
			edge := &models.FlowEdge{From: out1, To: in2, Capacity: 1}
			rev := &models.FlowEdge{From: in2, To: out1, Capacity: 0, Reverse: edge}
			edge.Reverse = rev
			out1.Edges = append(out1.Edges, edge)
			in2.Edges = append(in2.Edges, rev)
		}
		if out2 != nil && in1 != nil {
			edge := &models.FlowEdge{From: out2, To: in1, Capacity: 1}
			rev := &models.FlowEdge{From: in1, To: out2, Capacity: 0, Reverse: edge}
			edge.Reverse = rev
			out2.Edges = append(out2.Edges, edge)
			in1.Edges = append(in1.Edges, rev)
		}
	}
	
	return &models.FlowGraph{
		Nodes: nodes,
		Start: nodes[start+"_out"],
		End:   nodes[end+"_in"],
	}
}
