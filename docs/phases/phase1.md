# Phase 1 — Data Ingestion

| Task ID | Description | Deliverable |
|---------|-------------|-------------|
| **T1.1** | Implement `parser.ReadFile(path string) (string, error)` | Returns raw bytes as string |
| **T1.2** | Implement `parser.SplitLines(raw string) []string` | Split by `\n`, trim `\r` |
| **T1.3** | Implement `parser.FilterComments(lines []string) []string` | Remove lines starting with `#` that are NOT `##start` or `##end`. Lines starting with `##` that are not `##start`/`##end` are **ignored** (treated as comments) |
| **T1.4** | Implement `parser.TrimWhitespace(lines []string) []string` | Trim leading/trailing spaces from each line, remove empty lines |
| **T1.5** | Write unit tests for T1.1–T1.4 | `parser_test.go` |

**Edge cases to handle:**
- Windows line endings (`\r\n`)
- Trailing newline at EOF
- Empty file
- File with only comments
- Mixed `#comment` and `##start`/`##end`
- Comment lines with leading whitespace (e.g., `  #comment`)
