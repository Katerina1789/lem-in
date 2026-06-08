# Phase 3 — Graph Construction

| Task ID | Description |
|---------|-------------|
| **T3.1** | Build `map[string]*Room` from validated room list |
| **T3.2** | Build adjacency list `map[string][]string` from links (bidirectional) |
| **T3.3** | Implement `graph.IsConnected(start, end string) bool` — simple BFS to check reachability |
| **T3.4** | Implement `graph.SplitForFlow(g *Graph) *FlowGraph` — room splitting transformation |

**Room Splitting Detail:**

For each room `R` (except start/end), create two nodes:
- `R_in` — all incoming edges connect here
- `R_out` — all outgoing edges connect from here
- Edge `R_in → R_out` with capacity 1

For start: only `start_out` exists.  
For end: only `end_in` exists.

For each original link `A-B`:
- `A_out → B_in` (capacity 1)
- `B_out → A_in` (capacity 1)

This enforces the "one ant per room" constraint at the flow level.

```go
type FlowNode struct {
    Name     string
    Original string // original room name (empty for split nodes)
    IsIn     bool
    IsOut    bool
    Edges    []*FlowEdge
}

type FlowEdge struct {
    To       *FlowNode
    Capacity int
    Flow     int
    Reverse  *FlowEdge // for residual graph
}
```
