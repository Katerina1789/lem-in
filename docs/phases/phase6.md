# Phase 6 — Output Formatting

| Task ID | Description |
|---------|-------------|
| **T6.1** | Implement `formatter.FormatFileEcho(input string) string` — echoes the original file content (ant count, rooms, links) |
| **T6.2** | Implement `formatter.FormatTurns(turns []Turn) string` — formats each turn as `Lx-y Lz-w ...\n` |
| **T6.3** | Implement `formatter.FormatFull(input string, turns []Turn) string` — combines echo + blank line + turn output |
| **T6.4** | Ensure exact format: no trailing spaces, no extra newlines |

**Output format specification:**

```
<ant_count>
<room_lines>
<link_lines>

L<ant>-<room> L<ant>-<room> ...
L<ant>-<room> ...
```

- Each turn on its own line
- Multiple moves on same line separated by single space
- No trailing space at end of line
- Final newline at end of file is acceptable
- The blank line between file echo and first turn line is **required**
