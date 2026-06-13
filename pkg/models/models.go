package models

// Room represents a single room in the colony
type Room struct {
	Name    string
	CoordX  int
	CoordY  int
	IsStart bool
	IsEnd   bool
	// Internal use only — not from input
	Neighbors []string // adjacent room names
}

// Link represents a tunnel between two rooms
type Link struct {
	Room1 string
	Room2 string
}

// Ant represents an ant in transit
type Ant struct {
	ID       int
	Path     []string // sequence of room names (excluding start)
	Position int      // index in Path (0 = first room after start)
	Arrived  bool
}

// Path is a sequence of room names from start to end
type Path struct {
	Rooms  []string // includes start and end
	Length int      // number of edges (len(Rooms)-1)
}

// Graph is the full colony representation
type Graph struct {
	Rooms     map[string]*Room
	Links     []Link
	Adjacency map[string][]string // room name → neighbor names
}

// TurnMove records a single ant's movement in a turn
type TurnMove struct {
	AntID int
	Room  string
}

// Turn is a slice of moves for one turn
type Turn []TurnMove

// Solution holds the complete result
type Solution struct {
	Turns []Turn
}

// FlowGraph represents the split-node graph used for pathfinding
type FlowGraph struct {
	Nodes map[string]*FlowNode
	Start *FlowNode
	End   *FlowNode
}

// FlowNode represents a split node (in/out) or start/end node
type FlowNode struct {
	Name     string
	Original string      // original room name
	IsIn     bool        // true if this is R_in
	IsOut    bool        // true if this is R_out
	Edges    []*FlowEdge // outgoing edges
}

// FlowEdge represents a directed edge in the flow graph
type FlowEdge struct {
	From     *FlowNode
	To       *FlowNode
	Capacity int
	Flow     int
	Reverse  *FlowEdge // pointer to the reverse edge in the residual graph
}

