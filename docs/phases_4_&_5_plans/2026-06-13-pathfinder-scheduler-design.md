# Design Document: Pathfinder & Scheduler Implementation (Bypassing Phase Gates)

**Date**: 2026-06-13  
**Status**: Approved

## 1. Overview
Due to project constraints, we are bypassing the sequential phase gates to implement **Phase 4 (Pathfinding)** and **Phase 5 (Ant Scheduling)** before **Phase 3 (Graph Construction)** is fully implemented. To achieve this, we will stub the required flow graph models and use mocked inputs within our unit tests to verify the correctness of the pathfinding and scheduling algorithms.

---

## 2. Shared Data Structures
We will define the flow graph and residual network structures directly in `pkg/models/models.go` to keep them globally accessible:

```go
// FlowGraph represents the split-node graph used for pathfinding
type FlowGraph struct {
	Nodes map[string]*FlowNode
	Start *FlowNode
	End   *FlowNode
}

// FlowNode represents a split node (in/out) or start/end node
type FlowNode struct {
	Name     string
	Original string      // original room name (e.g., "room1")
	IsIn     bool        // true if this is R_in
	IsOut    bool        // true if this is R_out
	Edges    []*FlowEdge // outgoing edges (including residual/back edges)
}

// FlowEdge represents a directed edge in the flow graph
type FlowEdge struct {
	From     *FlowNode
	To       *FlowNode
	Capacity int
	Flow     int
	Reverse  *FlowEdge // pointer to the reverse edge in the residual graph
}
```

---

## 3. Phase 4: Pathfinder Implementation (`pkg/pathfinder/`)
The pathfinder uses the Edmonds-Karp algorithm (BFS on a residual network) to find the maximum number of vertex-disjoint paths and selects the subset of paths that minimizes total execution turns.

### A. Residual BFS (`pkg/pathfinder/bfs.go`)
* Finds the shortest path in the residual graph from start to end.
* An edge is only traversed if `Capacity - Flow > 0`.
* Returns `[]*FlowEdge` (the path of edges) to facilitate flow augmentation.

### B. Edmonds-Karp (`pkg/pathfinder/edmonds_karp.go`)
* Repeatedly calls the residual BFS to find augmenting paths.
* For each augmenting path, increases flow along forward edges and decreases flow along reverse edges.
* Once no augmenting paths exist, traces paths using edges where `Flow > 0`.
* Collapses split nodes (e.g., `["start_out", "room_in", "room_out", "end_in"]` -> `["start", "room", "end"]`) and returns `[][]string`.

### C. Path Selector (`pkg/pathfinder/selector.go`)
* Evaluates subsets of paths $k \in [1, \text{len(allPaths)}]$ sorted by length.
* Uses binary search to find the minimum turns required to distribute all ants across the $k$ paths:
  * `low = 0, high = antCount + max(L_i)`
  * `capacity = sum(max(0, mid - L_i + 1))`
* Returns the subset of paths that minimizes the turn count.

---

## 4. Phase 5: Ant Scheduling Implementation (`pkg/scheduler/`)
Once the optimal paths are selected, we distribute ants and simulate their movement turn-by-turn.

### A. Ant Distributor (`pkg/scheduler/distributor.go`)
* Distributes ants greedily using the formula: assign the next ant to the path that minimizes `len(path) + assigned_ants`. Ties are broken by choosing the shorter path first.
* Creates `*models.Ant` objects with unique 1-based IDs, setting their `Path` to exclude the start room and initializing `Position = -1`.

### B. Turn Simulator (`pkg/scheduler/simulator.go`)
* Simulates moves turn-by-turn.
* Active ants are processed from front to back on each path, sliding them forward into the next room if it is unoccupied (or if it is the end room).
* New ants are injected into the first room of each path if the room is empty.
* Returns a slice of `models.Turn` detailing each ant's moves.
