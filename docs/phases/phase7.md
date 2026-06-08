# Phase 7 — Integration & QA

| Task ID | Description |
|---------|-------------|
| **T7.1** | Integrate all packages in `main.go` |
| **T7.2** | Run all provided examples, verify output matches exactly |
| **T7.3** | Run all bad examples, verify error message |
| **T7.4** | Performance test: example06 with 100 ants (< 1.5 min) |
| **T7.5** | Performance test: example07 with 1000 ants (< 2.5 min) |
| **T7.6** | Fuzz test: generate random valid/invalid inputs |
| **T7.7** | Memory leak check with large inputs |
| **T7.8** | Write comprehensive test suite |

## Test Cases to Cover

### Valid cases:
- Single path, single ant
- Single path, multiple ants
- Multiple disjoint paths
- Overlapping paths (where optimal subset must be chosen)
- Start and end directly connected
- Large ant count on short path
- Rooms with non-numeric names
- Comments scattered throughout
- Extra `##` commands (ignored)

### Invalid cases (30+):
- Missing ant count
- Ant count = 0
- Ant count negative
- No `##start`
- No `##end`
- Duplicate `##start`
- Duplicate `##end`
- Room name with `L` prefix
- Room name with `#`
- Duplicate room name
- Duplicate coordinates
- Invalid coordinate (letters)
- Missing coordinate
- Link with unknown room
- Self-link
- Duplicate link
- No links
- No path from start to end
- Empty file
- `##start` is also `##end`
- Room line with extra spaces
- Link with spaces
- Link with multiple dashes
- Ant count line with trailing text
- Room after links section (out of order)
