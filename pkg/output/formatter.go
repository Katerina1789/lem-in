package output

import (
	"fmt"
	"strings"

	"lem-in/pkg/models"
)

// FormatFileEcho returns the original file content normalized to end with a single newline.
func FormatFileEcho(input string) string {
	return strings.TrimRight(input, "\r\n") + "\n"
}

// FormatTurns formats a slice of turns as "Lx-room Ly-room\n" per turn.
// No trailing space per line; each turn ends with a newline.
func FormatTurns(turns []models.Turn) string {
	var sb strings.Builder
	for _, turn := range turns {
		parts := make([]string, len(turn))
		for i, move := range turn {
			parts[i] = fmt.Sprintf("L%d-%s", move.AntID, move.Room)
		}
		sb.WriteString(strings.Join(parts, " "))
		sb.WriteByte('\n')
	}
	return sb.String()
}
