package parser

import (
	"os"
	"strings"
)

// ReadFile reads the entire file and returns its content as a string
func ReadFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// SplitLines splits raw content by \n and trims \r (handles Windows line endings)
func SplitLines(raw string) []string {
	lines := strings.Split(raw, "\n")
	for i := range lines {
		lines[i] = strings.TrimRight(lines[i], "\r")
	}
	return lines
}

// FilterComments removes lines starting with # (except ##start and ##end)
// Also removes lines starting with ## that are not ##start or ##end (treated as comments)
func FilterComments(lines []string) []string {
	var result []string
	for _, line := range lines {
		// Keep ##start and ##end
		if line == "##start" || line == "##end" {
			result = append(result, line)
		} else if strings.HasPrefix(line, "#") {
			// Skip all other # lines (comments)
			continue
		} else {
			result = append(result, line)
		}
	}
	return result
}

// TrimWhitespace trims leading/trailing spaces from each line and removes empty lines
func TrimWhitespace(lines []string) []string {
	var result []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}
