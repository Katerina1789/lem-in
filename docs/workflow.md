# Lem-in — Program Workflow

> Visualization of the program's complete logic and data flow.

---

## What the program does

Reads an **ant farm** (a graph of rooms connected by tunnels), moves N ants from `##start` to `##end` in the **minimum number of turns**. Rule: each room holds **only one** ant per turn (except start/end).

---

## Full Pipeline

```
┌─────────────────────────────────────────────────────────────────────┐
│                           INPUT FILE                                │
│                                                                     │
│  4              ← number of ants                                    │
│  ##start                                                            │
│  a 0 0          ← room "a" at coordinates (0,0)                     │
│  ##end                                                              │
│  d 3 0                                                              │
│  b 1 0                                                              │
│  c 2 0                                                              │
│  a-b             ← link (tunnel)                                    │
│  b-c                                                                │
│  c-d                                                                │
└──────────────────────────┬──────────────────────────────────────────┘
                           │  raw string
                           ▼
╔═════════════════════════════════════════════════════════════════════╗
║  PHASE 1 — PARSER              pkg/parser/parser.go                 ║
║                                                                     ║
║  ReadFile(path)      → string          (reads the file)             ║
║  SplitLines(raw)     → []string        (splits into lines)          ║
║  FilterComments()    → []string        (removes # comments,         ║
║                                         keeps ##start / ##end)      ║
║  TrimWhitespace()    → []string        (removes spaces/blank lines) ║
╚══════════════════════════════╤══════════════════════════════════════╝
                               │  []string  (clean lines)
                               ▼
╔═════════════════════════════════════════════════════════════════════╗
║  PHASE 2 — VALIDATOR           (pending implementation)             ║
║                                                                     ║
║  Validate(lines) → (*Graph, int, error)                             ║
║                                                                     ║
║  Checks:                                                            ║
║  • Number of ants > 0                                               ║
║  • Exactly one ##start and one ##end                                ║
║  • Valid room names (no L, # prefix)                                ║
║  • Unique names and coordinates                                     ║
║  • Links between known rooms, no duplicates                         ║
╚══════════════════════════════╤══════════════════════════════════════╝
                               │
                    ┌──────────┴─────────────┐
                    │                        │
              *models.Graph               int
              ┌──────────────┐          (antCount = 4)
              │ Rooms:       │
              │  "a" → Room  │
              │  "b" → Room  │
              │  "c" → Room  │
              │  "d" → Room  │
              │ Links:       │
              │  [a-b, b-c,  │
              │   c-d]       │
              │ Adjacency:   │
              │  "a"→["b"]   │
              │  "b"→["a","c"]│
              │  "c"→["b","d"]│
              │  "d"→["c"]   │
              └──────┬───────┘
                     │
                     ▼
╔═════════════════════════════════════════════════════════════════════╗
║  PHASE 3 — GRAPH CONSTRUCTION                                       ║
║                                                                     ║
║  SplitForFlow(g *Graph) → *FlowGraph                                ║
║                                                                     ║
║  Each room R is split into two nodes:                               ║
║                                                                     ║
║   R_in  ──[cap:1]──→  R_out                                         ║
║                                                                     ║
║  This enforces the capacity of 1 per room.                          ║
║                                                                     ║
║  Example for a-b-c-d:                                               ║
║                                                                     ║
║  start_out ──1──→ b_in ──1──→ b_out ──1──→ c_in ──1──→ c_out        ║
║                                                    ──1──→ end_in    ║
║  (start and end are not split: they only have _out and _in)         ║
╚══════════════════════════════╤══════════════════════════════════════╝
                               │  *models.FlowGraph
                               ▼
╔═════════════════════════════════════════════════════════════════════╗
║  PHASE 4 — PATHFINDING         pkg/pathfinder/                      ║
║                                                                     ║
║  ┌─────────────────────────────────────────────────────────────┐    ║
║  │  ALGORITHM: Edmonds-Karp (Max-Flow with BFS)                │    ║
║  │                                                             │    ║
║  │  loop:                                                      │    ║
║  │    path = FindShortestPath(g, start, end)  ← BFS            │    ║
║  │    if path == nil → stop                                    │    ║
║  │    augment Flow by 1 on each edge of the path               │    ║
║  │    (and decrease the reverse edge for the residual graph)   │    ║
║  │                                                             │    ║
║  │  Result: each edge with flow=1 belongs to a path            │    ║
║  └─────────────────────────────────────────────────────────────┘    ║
║                                                                     ║
║  FindAllDisjointPaths(g) → [][]string                               ║
║                                                                     ║
║  Returns e.g.:  [["a","b","c","d"]]  (only 1 path here)             ║
║                                                                     ║
║  ┌─────────────────────────────────────────────────────────────┐    ║
║  │  SelectOptimalPaths(allPaths, antCount) → [][]string        │    ║
║  │                                                             │    ║
║  │  Tries k = 1, 2, 3... paths:                                │    ║
║  │    computeTurns(selected_k, N):                             │    ║
║  │      binary search for min T such that:                     │    ║
║  │      Σ max(0, T - len(path_i) + 1) >= N                     │    ║
║  │    → keeps the k with the min T                             │    ║
║  └─────────────────────────────────────────────────────────────┘    ║
╚══════════════════════════════╤══════════════════════════════════════╝
                               │  [][]string  (optimal paths)
                               ▼
╔═════════════════════════════════════════════════════════════════════╗
║  PHASE 5 — DISTRIBUTION        pkg/scheduler/distributor.go         ║
║                                                                     ║
║  DistributeAnts(paths, numAnts) → []*models.Ant                     ║
║                                                                     ║
║  For each new ant:                                                  ║
║    cost(path_i) = len(path_i) + assignedCount(path_i)               ║
║    → assign it to the path with min cost                            ║
║                                                                     ║
║  Each Ant:                                                          ║
║  ┌────────────────────────────┐                                     ║
║  │ ID       = 1,2,3,4         │                                     ║
║  │ Path     = ["b","c","d"]   │  (without start)                    ║
║  │ Position = -1              │  (-1 = not yet started)             ║
║  │ Arrived  = false           │                                     ║
║  └────────────────────────────┘                                     ║
╚══════════════════════════════╤══════════════════════════════════════╝
                               │  []*models.Ant
                               ▼
╔═════════════════════════════════════════════════════════════════════╗
║  PHASE 5 — SIMULATION          pkg/scheduler/simulator.go           ║
║                                                                     ║
║  Simulate(ants, graph) → []models.Turn                              ║
║                                                                     ║
║  Each turn:                                                         ║
║  1. roomOccupied = map[string]bool  (which rooms are occupied)      ║
║  2. For each path, moves front ants first (avoids deadlock —        ║
║     the rear ant cannot block the front one)                        ║
║  3. If nextRoom is free (or is end) → move                          ║
║  4. If Position == -1 and firstRoom is free → start                 ║
║  5. If ant reaches end → Arrived = true                             ║
║  6. Stops when all Arrived == true                                  ║
╚══════════════════════════════╤══════════════════════════════════════╝
                               │  []models.Turn
                               ▼
╔═════════════════════════════════════════════════════════════════════╗
║  PHASE 6 — OUTPUT FORMATTER    pkg/output/formatter.go              ║
║                                                                     ║
║  FormatFileEcho(input)       → echo of the file (1 trailing \n)     ║
║  FormatTurns(turns)          → "L{id}-{room} ...\n" per Turn        ║
║  FormatFull(input, turns)    → FileEcho + "\n" + Turns              ║
╚══════════════════════════════╤══════════════════════════════════════╝
                               │  string
                               ▼
┌─────────────────────────────────────────────────────────────────────┐
│                           STDOUT                                    │
│                                                                     │
│  4                    ← echo of input                               │
│  ##start                                                            │
│  a 0 0                                                              │
│  ##end                                                              │
│  d 3 0                                                              │
│  b 1 0                                                              │
│  c 2 0                                                              │
│  a-b                                                                │
│  b-c                                                                │
│  c-d                                                                │
│                       ← blank line (required)                       │
│  L1-b                 ← Turn 1                                      │
│  L1-c L2-b            ← Turn 2                                      │
│  L1-d L2-c L3-b       ← Turn 3                                      │
│  L2-d L3-c L4-b       ← Turn 4                                      │
│  L3-d L4-c            ← Turn 5                                      │
│  L4-d                 ← Turn 6                                      │
└─────────────────────────────────────────────────────────────────────┘
```

---

## Data Types — What is passed between phases

```
string          ──parser──►     []string           (file lines)
[]string        ──validator──►  *Graph  +  int     (graph + N ants)
*Graph          ──splitFlow──►  *FlowGraph         (split-node graph)
*FlowGraph      ──edmondskarp►  [][]string         (all disjoint paths)
[][]string+int  ──selector──►   [][]string         (optimal paths)
[][]string+int  ──distribute──► []*Ant             (ants with assigned paths)
[]*Ant          ──simulate──►   []Turn             (moves per turn)
string+[]Turn   ──format──►     string             (final output)
```

---

## Data Structures (models.go)

```
models.Graph                        models.FlowGraph
┌────────────────────────┐          ┌──────────────────────────┐
│ Rooms  map[string]*Room│          │Nodes map[string]*FlowNode│
│ Links  []Link          │          │Start *FlowNode           │
│ Adjacency map[string]  │          │End   *FlowNode           │
│          []string      │          └──────────────────────────┘
└────────────────────────┘
                                    models.FlowNode
models.Room                         ┌──────────────────────────┐
┌────────────────────────┐          │Name     string           │
│ Name    string         │          │Original string           │
│ CoordX  int            │          │IsIn     bool             │
│ CoordY  int            │          │IsOut    bool             │
│ IsStart bool           │          │Edges    []*FlowEdge      │
│ IsEnd   bool           │          └──────────────────────────┘
│ Neighbors []string     │
└────────────────────────┘          models.FlowEdge
                                    ┌──────────────────────────┐
models.Ant                          │From     *FlowNode        │
┌────────────────────────┐          │To       *FlowNode        │
│ ID       int           │          │Capacity int              │
│ Path     []string      │          │Flow     int              │
│ Position int (-1=start)│          │Reverse  *FlowEdge        │
│ Arrived  bool          │          └──────────────────────────┘
└────────────────────────┘

models.TurnMove                     models.Turn
┌────────────────────────┐          ┌──────────────────────────┐
│ AntID  int             │          │ []TurnMove               │
│ Room   string          │          │ (all moves               │
└────────────────────────┘          │  in a single turn)       │
                                    └──────────────────────────┘
```

---

## Room Splitting — How the FlowGraph works

The core idea: to enforce "one ant per room", each room `R` becomes **two nodes** with capacity 1 between them.

```
Original graph:          After split:

  a ──── b ──── c         start_out ──1──► b_in ──1──► b_out ──1──► c_in ──1──► c_out
                                                                                   │
                                                                                  1
                                                                                   ▼
                                                                                end_in

  Each edge capacity = 1
  Reverse edges (for residual graph): same structure with Flow in the opposite direction
```

---

## Edmonds-Karp Algorithm — Step by step

```
Iteration 1:
  BFS finds: start → b_in → b_out → c_in → c_out → end
  → Flow +1 on each edge of this path

Iteration 2:
  BFS searches for another path in the residual graph
  → None found (all edges saturated or blocked)
  → Stop

Result: 1 disjoint path = ["a", "b", "c", "d"]
```

---

## Optimal Path Selection — computeTurns

```
We have 4 ants, 1 path of length 3 (a→b→c→d = 3 edges):

Binary search for min T such that: max(0, T - 3 + 1) >= 4
  T=4: max(0, 4-3+1) = 2  < 4  ✗
  T=6: max(0, 6-3+1) = 4  >= 4 ✓
  T=5: max(0, 5-3+1) = 3  < 4  ✗
  T=6: ✓ → answer = 6 turns

If we had 2 paths of length 3 and 5:
  T=5: max(0,5-3+1) + max(0,5-5+1) = 3+1 = 4 >= 4 ✓ → 5 turns
  → 2 paths is better than 1!
```

---

## Package Dependencies

```
main.go
  └── pkg/parser
  └── pkg/output
        └── pkg/models

pkg/parser
  └── (stdlib only)

pkg/models
  └── (no dependencies)

pkg/pathfinder
  └── pkg/models

pkg/scheduler
  └── pkg/models

pkg/output
  └── pkg/models
```
