# Phase 2 — Validation

This is the most critical phase. Must handle **30+ error cases**.

## Stage 1: Ant Count Validation

| Task ID | Description |
|---------|-------------|
| **T2.1** | Find first non-comment, non-empty line → must be a single positive integer |
| **T2.2** | Validate: integer > 0, no trailing characters, no leading zeros (except "0" alone) |
| **T2.3** | Store ant count. Return error if not found before `##start` |

**Error cases:**
- No ant count line
- Ant count = 0
- Ant count negative
- Ant count with decimals
- Ant count with letters
- Multiple integers on ant count line
- Empty file before `##start`

## Stage 2: Room Validation

| Task ID | Description |
|---------|-------------|
| **T2.4** | Scan for `##start` and `##end` markers. Record the *next* line as the designated room |
| **T2.5** | Validate each room line format: exactly 3 space-separated tokens |
| **T2.6** | Validate room name: no `L` prefix, no `#` prefix, no spaces within name |
| **T2.7** | Validate coordinates: both are valid integers (no overflow) |
| **T2.8** | Validate no duplicate room names |
| **T2.9** | Validate no duplicate coordinate pairs `(x, y)` |
| **T2.10** | Validate exactly one `##start` and one `##end` exist |
| **T2.11** | Validate start and end rooms are different rooms |
| **T2.12** | Validate room name is not empty |

**Error cases:**
- Missing `##start`
- Missing `##end`
- Multiple `##start`
- Multiple `##end`
- Room name starts with `L`
- Room name starts with `#`
- Room name contains spaces
- Duplicate room name
- Duplicate coordinates
- Invalid coordinate (non-integer)
- Overflow coordinate (> math.MaxInt32)
- Room line with 2 tokens (missing coordinate)
- Room line with 4+ tokens
- `##start`/`##end` not followed by a room line
- `##start`/`##end` followed by another command
- Same room is both start and end

## Stage 3: Link Validation

| Task ID | Description |
|---------|-------------|
| **T2.13** | Identify link lines: contain exactly one `-` and no spaces |
| **T2.14** | Split on `-`, validate both room names exist in room set |
| **T2.15** | Validate no duplicate links (A-B same as B-A) |
| **T2.16** | Validate no self-links (A-A) |
| **T2.17** | Validate at least one link exists |
| **T2.18** | Validate there is at least one path from start to end *(deferred to Phase 4)* |

**Error cases:**
- Link with unknown room name
- Duplicate link (forward/reverse)
- Self-referencing link
- Link with spaces
- Link with multiple `-`
- Link with empty room name on either side
- No links at all
- `##start`/`##end` appearing in link section

## Validation Output

| Task ID | Description |
|---------|-------------|
| **T2.19** | Create error types: `ErrInvalidAntCount`, `ErrNoStartRoom`, `ErrNoEndRoom`, `ErrDuplicateRoom`, `ErrInvalidCoords`, `ErrUnknownRoomInLink`, `ErrDuplicateLink`, `ErrSelfLink`, etc. |
| **T2.20** | Implement `validator.Validate(lines []string) (*models.Graph, int, error)` — returns populated graph (rooms + links) and ant count, or error |
| **T2.21** | In `main.go`, on error: print `ERROR: invalid data format` to stderr and exit with code 1. Optionally include detail message on separate line. |
