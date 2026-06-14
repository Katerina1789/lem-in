package graph

import (
	"testing"

	"lem-in/pkg/models"
)

// makeGraph builds a *models.Graph from room names and link pairs.
// First room is start, last room is end.
func makeGraph(rooms []string, links [][2]string) *models.Graph {
	roomMap := make(map[string]*models.Room)
	for i, name := range rooms {
		roomMap[name] = &models.Room{
			Name:    name,
			IsStart: i == 0,
			IsEnd:   i == len(rooms)-1,
		}
	}

	var linkSlice []models.Link
	for _, l := range links {
		linkSlice = append(linkSlice, models.Link{Room1: l[0], Room2: l[1]})
	}

	return &models.Graph{Rooms: roomMap, Links: linkSlice}
}

// --- BuildAdjacency ---

// Verifies both directions of a link are recorded in the adjacency list.
func TestBuildAdjacency_Basic(t *testing.T) {
	g := makeGraph(
		[]string{"start", "A", "end"},
		[][2]string{{"start", "A"}, {"A", "end"}},
	)
	BuildAdjacency(g)

	// start must list A as neighbor
	if !containsStr(g.Adjacency["start"], "A") {
		t.Errorf("expected start→A, got %v", g.Adjacency["start"])
	}
	// A must list both neighbors (bidirectional)
	if !containsStr(g.Adjacency["A"], "start") {
		t.Errorf("expected A→start, got %v", g.Adjacency["A"])
	}
	if !containsStr(g.Adjacency["A"], "end") {
		t.Errorf("expected A→end, got %v", g.Adjacency["A"])
	}
	// end must list A as neighbor
	if !containsStr(g.Adjacency["end"], "A") {
		t.Errorf("expected end→A, got %v", g.Adjacency["end"])
	}
}

// Verifies a room with two links has two neighbors.
func TestBuildAdjacency_MultipleLinks(t *testing.T) {
	g := makeGraph(
		[]string{"start", "A", "B", "end"},
		[][2]string{{"start", "A"}, {"start", "B"}, {"A", "end"}, {"B", "end"}},
	)
	BuildAdjacency(g)

	// start connects to both A and B
	if len(g.Adjacency["start"]) != 2 {
		t.Errorf("expected 2 neighbors for start, got %d", len(g.Adjacency["start"]))
	}
}

// Verifies no neighbors are added when there are no links.
func TestBuildAdjacency_EmptyLinks(t *testing.T) {
	g := makeGraph([]string{"start", "end"}, nil)
	BuildAdjacency(g)

	// No links means no neighbors
	if len(g.Adjacency["start"]) != 0 {
		t.Errorf("expected 0 neighbors, got %d", len(g.Adjacency["start"]))
	}
}

// --- IsConnected ---

// Verifies a direct link between start and end returns true.
func TestIsConnected_DirectLink(t *testing.T) {
	g := makeGraph(
		[]string{"start", "end"},
		[][2]string{{"start", "end"}},
	)
	BuildAdjacency(g)

	if !IsConnected(g, "start", "end") {
		t.Error("expected start→end to be connected")
	}
}

// Verifies a path through intermediate rooms is found.
func TestIsConnected_ChainPath(t *testing.T) {
	g := makeGraph(
		[]string{"start", "A", "B", "end"},
		[][2]string{{"start", "A"}, {"A", "B"}, {"B", "end"}},
	)
	BuildAdjacency(g)

	if !IsConnected(g, "start", "end") {
		t.Error("expected path through A→B to be found")
	}
}

// Verifies false is returned when end is not reachable.
func TestIsConnected_NoPath(t *testing.T) {
	g := makeGraph(
		[]string{"start", "A", "end"},
		[][2]string{{"start", "A"}},
	)
	BuildAdjacency(g)

	// A is a dead end — end is unreachable
	if IsConnected(g, "start", "end") {
		t.Error("expected no path from start to end")
	}
}

// Verifies false is returned when the graph has two disconnected halves.
func TestIsConnected_DisconnectedGraph(t *testing.T) {
	g := makeGraph(
		[]string{"start", "A", "B", "end"},
		[][2]string{{"start", "A"}, {"B", "end"}},
	)
	BuildAdjacency(g)

	// start-A and B-end are two separate components
	if IsConnected(g, "start", "end") {
		t.Error("expected disconnected graph to return false")
	}
}

// --- SplitForFlow ---

// Verifies start only gets _out and end only gets _in, with correct FlowGraph pointers.
func TestSplitForFlow_StartEndNodes(t *testing.T) {
	g := makeGraph(
		[]string{"start", "end"},
		[][2]string{{"start", "end"}},
	)
	fg := SplitForFlow(g)

	// Start must have _out only
	if fg.Nodes["start_out"] == nil {
		t.Error("expected start_out node")
	}
	if fg.Nodes["start_in"] != nil {
		t.Error("start should not have _in node")
	}

	// End must have _in only
	if fg.Nodes["end_in"] == nil {
		t.Error("expected end_in node")
	}
	if fg.Nodes["end_out"] != nil {
		t.Error("end should not have _out node")
	}

	// FlowGraph.Start and .End must point to the correct nodes
	if fg.Start != fg.Nodes["start_out"] {
		t.Error("FlowGraph.Start should point to start_out")
	}
	if fg.End != fg.Nodes["end_in"] {
		t.Error("FlowGraph.End should point to end_in")
	}
}

// Verifies an internal room is split into _in and _out with a capacity-1 internal edge.
func TestSplitForFlow_InternalRoomSplit(t *testing.T) {
	g := makeGraph(
		[]string{"start", "A", "end"},
		[][2]string{{"start", "A"}, {"A", "end"}},
	)
	fg := SplitForFlow(g)

	// Both split nodes must exist
	aIn := fg.Nodes["A_in"]
	aOut := fg.Nodes["A_out"]
	if aIn == nil || aOut == nil {
		t.Fatal("expected both A_in and A_out nodes")
	}

	// Flags and Original field must be set correctly
	if !aIn.IsIn || aIn.Original != "A" {
		t.Error("A_in: wrong flags or Original")
	}
	if !aOut.IsOut || aOut.Original != "A" {
		t.Error("A_out: wrong flags or Original")
	}

	// Internal capacity-1 edge A_in → A_out must exist
	internalEdge := findEdgeTo(aIn.Edges, aOut)
	if internalEdge == nil {
		t.Fatal("expected internal edge A_in → A_out")
	}
	if internalEdge.Capacity != 1 {
		t.Errorf("internal edge capacity: expected 1, got %d", internalEdge.Capacity)
	}

	// Reverse edge must be paired and have capacity 0
	if internalEdge.Reverse == nil {
		t.Error("internal edge missing reverse")
	}
	if internalEdge.Reverse.Capacity != 0 {
		t.Errorf("reverse edge capacity: expected 0, got %d", internalEdge.Reverse.Capacity)
	}
}

// Verifies tunnel edges are added in both directions between split nodes.
func TestSplitForFlow_TunnelEdges(t *testing.T) {
	g := makeGraph(
		[]string{"start", "A", "end"},
		[][2]string{{"start", "A"}, {"A", "end"}},
	)
	fg := SplitForFlow(g)

	startOut := fg.Nodes["start_out"]
	aIn := fg.Nodes["A_in"]
	aOut := fg.Nodes["A_out"]
	endIn := fg.Nodes["end_in"]

	// Forward tunnel: start_out → A_in
	if findEdgeTo(startOut.Edges, aIn) == nil {
		t.Error("expected tunnel edge start_out → A_in")
	}
	// Forward tunnel: A_out → end_in
	if findEdgeTo(aOut.Edges, endIn) == nil {
		t.Error("expected tunnel edge A_out → end_in")
	}
}

// Verifies each tunnel edge has a residual reverse edge on the destination node.
func TestSplitForFlow_ResidualEdgesExist(t *testing.T) {
	g := makeGraph(
		[]string{"start", "A", "end"},
		[][2]string{{"start", "A"}, {"A", "end"}},
	)
	fg := SplitForFlow(g)

	startOut := fg.Nodes["start_out"]
	aIn := fg.Nodes["A_in"]

	// A_in must carry the residual reverse of the start_out → A_in tunnel
	if findEdgeTo(aIn.Edges, startOut) == nil {
		t.Error("expected residual reverse edge A_in → start_out on A_in.Edges")
	}
}

// Verifies two parallel paths produce four internal split nodes.
func TestSplitForFlow_TwoPaths(t *testing.T) {
	g := makeGraph(
		[]string{"start", "A", "B", "end"},
		[][2]string{{"start", "A"}, {"A", "end"}, {"start", "B"}, {"B", "end"}},
	)
	fg := SplitForFlow(g)

	// A and B must each be split into _in and _out
	if fg.Nodes["A_in"] == nil || fg.Nodes["A_out"] == nil {
		t.Error("missing A split nodes")
	}
	if fg.Nodes["B_in"] == nil || fg.Nodes["B_out"] == nil {
		t.Error("missing B split nodes")
	}
}

// Verifies the total node count matches the expected split layout.
func TestSplitForFlow_NodeCount(t *testing.T) {
	g := makeGraph(
		[]string{"start", "A", "end"},
		[][2]string{{"start", "A"}, {"A", "end"}},
	)
	fg := SplitForFlow(g)

	// start_out + A_in + A_out + end_in = 4 nodes
	if len(fg.Nodes) != 4 {
		t.Errorf("expected 4 nodes, got %d", len(fg.Nodes))
	}
}

// --- test helpers ---

// containsStr returns true if target is in slice.
func containsStr(slice []string, target string) bool {
	for _, s := range slice {
		if s == target {
			return true
		}
	}
	return false
}

// findEdgeTo returns the first edge in edges whose To is target, or nil.
func findEdgeTo(edges []*models.FlowEdge, target *models.FlowNode) *models.FlowEdge {
	for _, e := range edges {
		if e.To == target {
			return e
		}
	}
	return nil
}
