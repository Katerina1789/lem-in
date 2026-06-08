# Phase 0 — Project Scaffolding

| Task ID | Description | Deliverable |
|---------|-------------|-------------|
| **T0.1** | Initialize Go module `lem-in` | `go.mod` |
| **T0.2** | Create directory structure as per file tree above | All dirs + placeholder files |
| **T0.3** | Create `main.go` skeleton: read `os.Args[1]`, call `os.ReadFile`, print content | `main.go` |
| **T0.4** | Create `pkg/models/models.go` with all shared structs | `models.go` |
| **T0.5** | Set up CI-like script: `go vet ./... && go test ./...` | `scripts/check.sh` |

**`pkg/models/models.go` specification:**

```go
package models

// Room represents a single room in the colony
type Room struct {
    Name      string
    CoordX    int
    CoordY    int
    IsStart   bool
    IsEnd     bool
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
```
