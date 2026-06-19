# Lem-in

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://golang.org/)
[![Zone01](https://img.shields.io/badge/Zone01-Athens-8BC34A?style=for-the-badge&logo=codeforces&logoColor=white)](https://zone01.gr/)
[![License](https://img.shields.io/badge/License-MIT-D32F2F?style=for-the-badge&logo=opensourceinitiative&logoColor=white)](LICENSE)

A command-line pathfinding application in Go that simulates ants moving through a colony using optimal multi-path routing.  
Built as part of the Zone01 Athens curriculum.

## Table of Contents

- [Description](#description)
- [Features](#features)
- [Repository Structure](#repository-structure)
- [How to Run](#how-to-run)
- [Input](#input)
- [Requirements](#requirements)
- [Algorithm Flow](#algorithm-flow)
- [Documentation](#documentation)
- [License](#license)

## Description

Lem-in is a high‑performance Go implementation of the classic ant‑farm pathfinding challenge.
The program reads a colony description (rooms, coordinates, tunnels and ant count), validates the input, computes the optimal set of disjoint paths using a max‑flow algorithm, distributes ants intelligently across those paths and simulates their movement turn‑by‑turn - always producing the minimum number of turns required to move all ants from ##start to ##end.

## Features

- Input Parsing – Read ants, rooms, coordinates and tunnels
- Validation – Detect malformed rooms, links, duplicates and missing start/end
- Graph Construction – Build colony graph with room-splitting
- Pathfinding – Compute disjoint shortest paths using max-flow
- Path Selection – Choose optimal path set for minimum turns
- Ant Distribution – Assign ants to paths based on cost
- Simulation – Move ants turn-by-turn following all movement rules
- Output – Print original input and ant movements in Lx-room format

## Repository Structure

```
lem-in/
├── audit/                     # Official audit test files
├── cmd/
│   └── main.go                # Entry point
├── docs/
│   ├── phases/                # Development phase documentation
│   ├── phases_4_&_5_plans/    # Pathfinding & scheduling planning
│   └── workflow.md            # High-level program workflow
├── pkg/
│   ├── graph/                 # Graph building & room splitting
│   ├── models/                # Core data structures
│   ├── output/                # Output formatting
│   ├── parser/                # Input parsing
│   ├── pathfinder/            # BFS + Edmonds–Karp + path selection
│   ├── scheduler/             # Ant distribution & simulation
│   └── validator/             # Input validation
├── scripts/
│   ├── audit.sh               # Automated audit tests
│   ├── check.sh               # Quick developer checks
│   └── run_audit.sh           # Quick audit test runner
├── go.mod
├── LICENSE
└── README.md
```

## How to Run

### Build

```bash
go build -o lem-in ./cmd
```

### Run

```bash
./lem-in ./audit/example00.txt

OR

go run ./cmd ./audit/example00.txt
```

### Run using scripts

```bash
bash scripts/check.sh
bash scripts/audit.sh
```

## Input

The program requires a valid lem-in input file containing:
- Number of ants  
- Room definitions (`name x y`)  
- Tunnels (`room1-room2`)  
- `##start` and `##end` markers  

Example:
```
3                  # Number of ants

##start            # Start marker
A 1 2              # Room definition: name x y

B 3 4              # Room definition
C 5 6              # Room definition

##end              # End marker
D 7 8              # Room definition

A-B                # Tunnel: room1-room2
B-C
C-D
```

## Requirements

- Go 1.21+
- Only standard library packages allowed
- Input file must follow the lem‑in specification:
- Coordinates must be integers
- Room names must not start with L or #
- No duplicate rooms or tunnels
- No self‑links
- No unknown rooms in links

## Algorithm Flow

1. **Parse Input**  
   Read ants, rooms, coordinates, tunnels and detect `##start` / `##end`.

2. **Validate Data**  
   Check for invalid names, duplicates, absense of ants, malformed lines.

3. **Build Graph**  
   Create a graph of rooms and tunnels and apply room-splitting to enforce one ant per room.

4. **Find Paths**  
   Use Edmonds–Karp (BFS-based max-flow) to extract disjoint shortest paths.

5. **Select Optimal Paths**  
   Choose the path set that results in the minimum number of turns.

6. **Distribute Ants**  
   Assign ants to paths based on path length and load.

7. **Simulate**  
   Move ants turn-by-turn and print only the ants that moved.

## Documentation

- [Architecture](./docs/workflow.md) - high‑level stystem flow
- [Phases](./docs/phases/) - step‑by‑step development task breakdown
- [Phases 4 & 5](./docs/phases_4_&_5_plans/) - detailed algorithm explanations

## License

This project is licensed under the [MIT License](LICENSE).
