# Pathfinder & Scheduler Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use advancedSkills:executing-plans to implement this plan task-by-task.

**Goal:** Implement Phase 4 (Pathfinding) and Phase 5 (Ant Scheduling) using stubbed models and mocked flow graphs for unit testing, bypassing Phase 3 gating.

**Architecture:** We define the necessary flow graph structs (`FlowGraph`, `FlowNode`, `FlowEdge`) in `pkg/models/models.go`. We implement the residual BFS and Edmonds-Karp disjoint path finder in `pkg/pathfinder/` along with optimal path selection, and ant distribution/turn simulation in `pkg/scheduler/`.

**Tech Stack:** Go Standard Library (no external dependencies)

---

### Task 1: Define Flow Graph Structs

**Files:**
- Modify: `pkg/models/models.go`

**Step 1: Write the failing test**
Since this task only introduces new types to enable compilation of future tasks, we write a compilation validation check.
Create `pkg/models/models_compile_test.go`:
```go
package models

import "testing"

func TestCompilation(t *testing.T) {
	var _ FlowGraph
	var _ FlowNode
	var _ FlowEdge
	t.Log("Models compiled successfully")
}
```

**Step 2: Run test to verify it fails**
Run: `go test ./pkg/models/...`
Expected: Compile error: "undefined: FlowGraph", "undefined: FlowNode", "undefined: FlowEdge"

**Step 3: Write minimal implementation**
Append the following structs to `pkg/models/models.go`:
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
```

**Step 4: Run test to verify it passes**
Run: `go test ./pkg/models/...`
Expected: PASS

**Step 5: Commit**
```bash
git add pkg/models/models.go pkg/models/models_compile_test.go
git commit -m "feat: define FlowGraph, FlowNode, and FlowEdge structures"
```

---

### Task 2: Implement Residual BFS

**Files:**
- Create: `pkg/pathfinder/bfs.go`
- Create: `pkg/pathfinder/bfs_test.go`
- Create: `pkg/pathfinder/helper_test.go`

**Step 1: Write the failing test**
Create `pkg/pathfinder/helper_test.go` first:
```go
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
```

Create `pkg/pathfinder/bfs_test.go`:
```go
package pathfinder

import (
	"testing"
)

func TestFindShortestPath(t *testing.T) {
	g := buildTestFlowGraph(
		[]string{"start", "A", "end"},
		[][]string{{"start", "A"}, {"A", "end"}},
		"start",
		"end",
	)

	path := FindShortestPath(g, "start_out", "end_in")
	if len(path) != 3 { // start_out -> A_in, A_in -> A_out, A_out -> end_in
		t.Fatalf("expected path length 3, got %d", len(path))
	}

	if path[0].From.Name != "start_out" || path[2].To.Name != "end_in" {
		t.Errorf("invalid path sequence")
	}
}
```

**Step 2: Run test to verify it fails**
Run: `go test ./pkg/pathfinder/...`
Expected: Compile error: "undefined: FindShortestPath"

**Step 3: Write minimal implementation**
Create `pkg/pathfinder/bfs.go`:
```go
package pathfinder

import "lem-in/pkg/models"

func FindShortestPath(g *models.FlowGraph, startName, endName string) []*models.FlowEdge {
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
```

**Step 4: Run test to verify it passes**
Run: `go test ./pkg/pathfinder/...`
Expected: PASS

**Step 5: Commit**
```bash
git add pkg/pathfinder/bfs.go pkg/pathfinder/bfs_test.go pkg/pathfinder/helper_test.go
git commit -m "feat: implement residual BFS for FlowGraph pathfinding"
```

---

### Task 3: Implement Edmonds-Karp Disjoint Path Finder

**Files:**
- Create: `pkg/pathfinder/edmonds_karp.go`
- Create: `pkg/pathfinder/edmonds_karp_test.go`

**Step 1: Write the failing test**
Create `pkg/pathfinder/edmonds_karp_test.go`:
```go
package pathfinder

import (
	"testing"
)

func TestFindAllDisjointPaths(t *testing.T) {
	g := buildTestFlowGraph(
		[]string{"start", "A", "B", "end"},
		[][]string{
			{"start", "A"}, {"A", "end"},
			{"start", "B"}, {"B", "end"},
		},
		"start",
		"end",
	)

	paths := FindAllDisjointPaths(g)
	if len(paths) != 2 {
		t.Fatalf("expected 2 disjoint paths, got %d", len(paths))
	}

	t.Logf("Paths found: %v", paths)
}
```

**Step 2: Run test to verify it fails**
Run: `go test -run TestFindAllDisjointPaths ./pkg/pathfinder/...`
Expected: Compile error: "undefined: FindAllDisjointPaths"

**Step 3: Write minimal implementation**
Create `pkg/pathfinder/edmonds_karp.go`:
```go
package pathfinder

import "lem-in/pkg/models"

func FindAllDisjointPaths(g *models.FlowGraph) [][]string {
	for {
		path := FindShortestPath(g, g.Start.Name, g.End.Name)
		if path == nil {
			break
		}
		for _, edge := range path {
			edge.Flow++
			edge.Reverse.Flow--
		}
	}

	var paths [][]string
	for _, edge := range g.Start.Edges {
		if edge.Flow > 0 {
			path := []string{g.Start.Original}
			curr := edge.To
			for curr != g.End {
				if curr.Original != "" && curr.IsIn {
					path = append(path, curr.Original)
				}
				var nextNode *models.FlowNode
				for _, nextEdge := range curr.Edges {
					if nextEdge.Flow > 0 {
						nextNode = nextEdge.To
						break
					}
				}
				if nextNode == nil {
					break
				}
				curr = nextNode
			}
			path = append(path, g.End.Original)
			paths = append(paths, path)
		}
	}
	return paths
}
```

**Step 4: Run test to verify it passes**
Run: `go test -run TestFindAllDisjointPaths ./pkg/pathfinder/...`
Expected: PASS

**Step 5: Commit**
```bash
git add pkg/pathfinder/edmonds_karp.go pkg/pathfinder/edmonds_karp_test.go
git commit -m "feat: implement Edmonds-Karp max flow disjoint path finder"
```

---

### Task 4: Implement Path Selector

**Files:**
- Create: `pkg/pathfinder/selector.go`
- Create: `pkg/pathfinder/selector_test.go`

**Step 1: Write the failing test**
Create `pkg/pathfinder/selector_test.go`:
```go
package pathfinder

import (
	"testing"
)

func TestSelectOptimalPaths(t *testing.T) {
	paths := [][]string{
		{"start", "A", "end"},             // length 2
		{"start", "B", "C", "end"},        // length 3
	}

	// 5 ants.
	// If 1 path: 5 turns.
	// If 2 paths:
	// - path 1 gets 3 ants -> 2 + 3 - 1 = 4 turns.
	// - path 2 gets 2 ants -> 3 + 2 - 1 = 4 turns.
	// Max turns = 4. 2 paths is optimal.
	best := SelectOptimalPaths(paths, 5)
	if len(best) != 2 {
		t.Fatalf("expected 2 paths selected, got %d", len(best))
	}
}
```

**Step 2: Run test to verify it fails**
Run: `go test -run TestSelectOptimalPaths ./pkg/pathfinder/...`
Expected: Compile error: "undefined: SelectOptimalPaths"

**Step 3: Write minimal implementation**
Create `pkg/pathfinder/selector.go`:
```go
package pathfinder

import "sort"

func SelectOptimalPaths(allPaths [][]string, antCount int) [][]string {
	if len(allPaths) == 0 || antCount <= 0 {
		return nil
	}

	sort.Slice(allPaths, func(i, j int) bool {
		return len(allPaths[i]) < len(allPaths[j])
	})

	var bestPaths [][]string
	bestTurns := -1

	for k := 1; k <= len(allPaths); k++ {
		selected := allPaths[:k]
		turns := computeTurns(selected, antCount)
		if bestTurns == -1 || turns < bestTurns {
			bestTurns = turns
			bestPaths = selected
		}
	}

	return bestPaths
}

func computeTurns(paths [][]string, antCount int) int {
	lengths := make([]int, len(paths))
	maxL := 0
	for i, p := range paths {
		lengths[i] = len(p) - 1
		if lengths[i] > maxL {
			maxL = lengths[i]
		}
	}

	low := 0
	high := antCount + maxL

	for low < high {
		mid := (low + high) / 2
		capacity := 0
		for _, L := range lengths {
			if mid >= L {
				capacity += mid - L + 1
			}
		}
		if capacity >= antCount {
			high = mid
		} else {
			low = mid + 1
		}
	}

	return low
}
```

**Step 4: Run test to verify it passes**
Run: `go test -run TestSelectOptimalPaths ./pkg/pathfinder/...`
Expected: PASS

**Step 5: Commit**
```bash
git add pkg/pathfinder/selector.go pkg/pathfinder/selector_test.go
git commit -m "feat: implement optimal path selector using turn cost optimization"
```

---

### Task 5: Implement Ant Distributor

**Files:**
- Create: `pkg/scheduler/distributor.go`
- Create: `pkg/scheduler/distributor_test.go`

**Step 1: Write the failing test**
Create `pkg/scheduler/distributor_test.go`:
```go
package scheduler

import (
	"testing"
)

func TestDistributeAnts(t *testing.T) {
	paths := [][]string{
		{"start", "A", "end"},
		{"start", "B", "C", "end"},
	}

	ants := DistributeAnts(paths, 5)
	if len(ants) != 5 {
		t.Fatalf("expected 5 ants, got %d", len(ants))
	}

	// Verify ant 1 has path excluding start
	if len(ants[0].Path) != 2 || ants[0].Path[0] != "A" || ants[0].Path[1] != "end" {
		t.Errorf("expected path without start, got %v", ants[0].Path)
	}
}
```

**Step 2: Run test to verify it fails**
Run: `go test ./pkg/scheduler/...`
Expected: Compile error: "undefined: DistributeAnts" (or package scheduler doesn't exist)

**Step 3: Write minimal implementation**
Create `pkg/scheduler/distributor.go`:
```go
package scheduler

import (
	"lem-in/pkg/models"
	"sort"
)

func DistributeAnts(paths [][]string, numAnts int) []*models.Ant {
	if len(paths) == 0 || numAnts <= 0 {
		return nil
	}

	sort.Slice(paths, func(i, j int) bool {
		return len(paths[i]) < len(paths[j])
	})

	assignedCounts := make([]int, len(paths))
	for i := 0; i < numAnts; i++ {
		bestPathIdx := 0
		bestCost := len(paths[0]) + assignedCounts[0]
		for p := 1; p < len(paths); p++ {
			cost := len(paths[p]) + assignedCounts[p]
			if cost < bestCost {
				bestCost = cost
				bestPathIdx = p
			}
		}
		assignedCounts[bestPathIdx]++
	}

	var ants []*models.Ant
	antID := 1
	for p, count := range paths {
		subPath := paths[p][1:]
		for c := 0; c < assignedCounts[p]; c++ {
			ants = append(ants, &models.Ant{
				ID:       antID,
				Path:     subPath,
				Position: -1,
				Arrived:  false,
			})
			antID++
		}
	}

	sort.Slice(ants, func(i, j int) bool {
		return ants[i].ID < ants[j].ID
	})

	return ants
}
```

**Step 4: Run test to verify it passes**
Run: `go test ./pkg/scheduler/...`
Expected: PASS

**Step 5: Commit**
```bash
git add pkg/scheduler/distributor.go pkg/scheduler/distributor_test.go
git commit -m "feat: implement ant distribution across paths"
```

---

### Task 6: Implement Turn Simulator

**Files:**
- Create: `pkg/scheduler/simulator.go`
- Create: `pkg/scheduler/simulator_test.go`

**Step 1: Write the failing test**
Create `pkg/scheduler/simulator_test.go`:
```go
package scheduler

import (
	"lem-in/pkg/models"
	"testing"
)

func TestSimulate(t *testing.T) {
	paths := [][]string{
		{"start", "A", "end"},
		{"start", "B", "C", "end"},
	}

	ants := DistributeAnts(paths, 3)
	turns := Simulate(ants, &models.Graph{})

	if len(turns) == 0 {
		t.Fatalf("expected turns simulation, got 0 turns")
	}

	// Verify all ants reached end
	for _, ant := range ants {
		if !ant.Arrived {
			t.Errorf("ant %d did not arrive", ant.ID)
		}
	}
}
```

**Step 2: Run test to verify it fails**
Run: `go test -run TestSimulate ./pkg/scheduler/...`
Expected: Compile error: "undefined: Simulate"

**Step 3: Write minimal implementation**
Create `pkg/scheduler/simulator.go`:
```go
package scheduler

import (
	"lem-in/pkg/models"
	"sort"
)

func Simulate(ants []*models.Ant, graph *models.Graph) []models.Turn {
	var turns []models.Turn

	for {
		allArrived := true
		for _, ant := range ants {
			if !ant.Arrived {
				allArrived = false
				break
			}
		}
		if allArrived {
			break
		}

		currentTurn := models.Turn{}
		roomOccupied := make(map[string]bool)
		for _, ant := range ants {
			if !ant.Arrived && ant.Position >= 0 {
				currRoom := ant.Path[ant.Position]
				if currRoom != ant.Path[len(ant.Path)-1] {
					roomOccupied[currRoom] = true
				}
			}
		}

		var pathList [][]string
		for _, ant := range ants {
			found := false
			for _, p := range pathList {
				if equalPaths(p, ant.Path) {
					found = true
					break
				}
			}
			if !found {
				pathList = append(pathList, ant.Path)
			}
		}

		for _, path := range pathList {
			var pathAnts []*models.Ant
			for _, ant := range ants {
				if equalPaths(ant.Path, path) {
					pathAnts = append(pathAnts, ant)
				}
			}

			sort.Slice(pathAnts, func(i, j int) bool {
				return pathAnts[i].Position > pathAnts[j].Position
			})

			for i := 0; i < len(pathAnts); i++ {
				ant := pathAnts[i]
				if ant.Arrived || ant.Position == -1 {
					continue
				}

				nextPos := ant.Position + 1
				nextRoom := ant.Path[nextPos]
				isEnd := nextPos == len(ant.Path)-1

				if isEnd || !roomOccupied[nextRoom] {
					currRoom := ant.Path[ant.Position]
					delete(roomOccupied, currRoom)

					ant.Position = nextPos
					if isEnd {
						ant.Arrived = true
					} else {
						roomOccupied[nextRoom] = true
					}

					currentTurn = append(currentTurn, models.TurnMove{
						AntID: ant.ID,
						Room:  nextRoom,
					})
				}
			}

			firstRoom := path[0]
			isFirstRoomEnd := len(path) == 1
			if isFirstRoomEnd || !roomOccupied[firstRoom] {
				var inactiveAnt *models.Ant
				for _, ant := range pathAnts {
					if ant.Position == -1 {
						inactiveAnt = ant
						break
					}
				}

				if inactiveAnt != nil {
					inactiveAnt.Position = 0
					if isFirstRoomEnd {
						inactiveAnt.Arrived = true
					} else {
						roomOccupied[firstRoom] = true
					}

					currentTurn = append(currentTurn, models.TurnMove{
						AntID: inactiveAnt.ID,
						Room:  firstRoom,
					})
				}
			}
		}

		if len(currentTurn) > 0 {
			turns = append(turns, currentTurn)
		} else {
			break
		}
	}

	return turns
}

func equalPaths(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
```

**Step 4: Run test to verify it passes**
Run: `go test -run TestSimulate ./pkg/scheduler/...`
Expected: PASS

**Step 5: Commit**
```bash
git add pkg/scheduler/simulator.go pkg/scheduler/simulator_test.go
git commit -m "feat: implement turn-by-turn ant movement simulation"
```
