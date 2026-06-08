# AGENTS.md — Lem-in Project Orchestration

## Project Overview

**lem-in** is an ant colony pathfinding simulator. Given a graph of rooms and tunnels with `N` ants starting at `##start` and needing to reach `##end`, the program finds the optimal set of vertex-disjoint paths and simulates ant movement turn-by-turn.

**Hard Constraints:**
- One ant per room (except `##start`/`##end`)
- One use per tunnel per turn
- Go standard library only
- Must handle all malformed input gracefully

---

## Agent Roles & Responsibilities

| Agent ID | Role | Focus |
|----------|------|-------|
| **A0** | Orchestrator | Oversees integration, runs tests, approves merges |
| **A1** | Parser Agent | File reading, line splitting, comment filtering |
| **A2** | Validator Agent | Stage 1–3 validation, error message generation |
| **A3** | Graph Builder Agent | Room/Link struct population, adjacency list construction |
| **A4** | Pathfinder Agent | Bidirectional BFS, disjoint path discovery, path selection |
| **A5** | Scheduler Agent | Ant-to-path assignment, turn-by-turn simulation |
| **A6** | Output Formatter Agent | Correct `Lx-y` formatting, final output assembly |
| **A7** | QA/Test Agent | Unit tests, integration tests, edge case files |

---

## Phase Execution Order

```
Phase 0: Project Scaffolding (A0)
Phase 1: Data Ingestion     (A1) ─┐
Phase 2: Validation         (A2) ─┤ (can partially overlap)
Phase 3: Graph Construction (A3) ─┘
Phase 4: Pathfinding        (A4)
Phase 5: Ant Scheduling     (A5)
Phase 6: Output Formatting  (A6)
Phase 7: Integration & QA   (A0 + A7)
```

Each phase is **gated** — the next phase cannot begin until the current phase passes all unit tests.

---

## Key Architectural Decisions

1. **Room Node Splitting**: To enforce the "one ant per room" constraint during pathfinding, each room (except start/end) is split into `room_in → room_out` with a capacity-1 edge between them. This is an *internal representation only*; output uses original room names.

2. **Pathfinding Algorithm**: 
   - Use **BFS** (Edmonds-Karp style) on the split-graph to find augmenting paths.
   - Collect **all** vertex-disjoint paths.
   - For each `k` from 1 to `maxPaths`, calculate total turns:  
     $T(k) = \max_{i \in [1..k]} \left( L_i + \left\lceil \frac{a_i}{1} \right\rceil - 1 \right)$  
     where $L_i$ = length of path $i$ (in edges), and $a_i$ are ants assigned to path $i$.
   - Select $k$ that minimizes $T(k)$.

3. **Ant Distribution**: Given $k$ paths of lengths $L_1 \leq L_2 \leq \ldots \leq L_k$, assign ants greedily: each path $i$ gets $\max(0, T_{target} - L_i + 1)$ ants, where $T_{target}$ is the target turn count, iteratively adjusted.

4. **Turn Simulation**: Queue-based. Each path maintains a queue of ants. At turn $t$, ant $j$ on path $i$ moves from room $p$ to room $p+1$ if room $p+1$ is empty. Start injects a new ant when the first room of a path is free.

---

## Communication Protocol

- All agents commit to feature branches: `feat/<agent>-<task>`
- PRs must pass `go vet`, `golint`, and all existing tests
- A0 reviews and merges
- Shared data structures live in `pkg/models/` and are **immutable** after Phase 3

---

## File Structure

```
lem-in/
├── main.go                  # Entry point
├── AGENTS.md                # This file
├── go.mod
├── pkg/
│   ├── models/
│   │   └── models.go        # Room, Link, Ant, Graph, Path structs
│   ├── parser/
│   │   ├── parser.go        # File → raw lines
│   │   └── parser_test.go
│   ├── validator/
│   │   ├── validator.go     # 3-stage validation
│   │   ├── errors.go        # Error types & messages
│   │   └── validator_test.go
│   ├── graph/
│   │   ├── builder.go       # Populate Room/Link slices
│   │   ├── split.go         # Room splitting for flow
│   │   └── graph_test.go
│   ├── pathfinder/
│   │   ├── bfs.go           # Bidirectional BFS
│   │   ├── edmonds_karp.go  # Augmenting paths
│   │   ├── selector.go      # Optimal path subset selection
│   │   └── pathfinder_test.go
│   ├── scheduler/
│   │   ├── distributor.go   # Ant → path assignment
│   │   ├── simulator.go     # Turn-by-turn simulation
│   │   └── scheduler_test.go
│   └── output/
│       ├── formatter.go     # Lx-y formatting
│       └── formatter_test.go
├── testdata/
│   ├── example00.txt
│   ├── example01.txt
│   ├── ...
│   ├── badexample00.txt
│   └── badexample01.txt
└── visualizer/              # Bonus
    └── main.go
```

---

## Gating Criteria Per Phase

| Phase | Entry Criteria | Exit Criteria |
|-------|---------------|---------------|
| 0 | — | `main.go` compiles, reads arg, prints file content |
| 1 | Phase 0 ✅ | All lines returned as `[]string`, `#` comments filtered |
| 2 | Phase 1 ✅ | 20+ error cases handled, valid input passes silently |
| 3 | Phase 2 ✅ | Graph built, adjacency lists correct, no leaks |
| 4 | Phase 3 ✅ | Shortest path found for examples 00–05 |
| 5 | Phase 4 ✅ | Turn counts match spec for all examples |
| 6 | Phase 5 ✅ | Output format matches spec byte-for-byte |
| 7 | All above ✅ | All example tests pass, time limits met |
```

---