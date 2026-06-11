package parser

import (
	"os"
	"strings"
)

// ReadFile reads the entire file and returns its contents as a string.
func ReadFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// SplitLines splits the raw string by newlines and removes carriage returns.
func SplitLines(raw string) []string {
	// Split by \n
	lines := strings.Split(raw, "\n")
	// Trim \r from each line
	for i := range lines {
		lines[i] = strings.TrimSuffix(lines[i], "\r")
	}
	return lines
}

// FilterComments removes lines starting with '#' that are NOT '##start' or '##end'.
// Lines starting with '##' that are not '##start' or '##end' are also ignored.
func FilterComments(lines []string) []string {
	var result []string
	for _, line := range lines {
		trimmed := strings.TrimLeft(line, " \t")
		// Keep lines that don't start with '#'
		if !strings.HasPrefix(trimmed, "#") {
			result = append(result, line)
			continue
		}
		// If it starts with '#', only keep ##start or ##end
		if trimmed == "##start" || trimmed == "##end" {
			result = append(result, line)
		}
		// Otherwise, skip it (it's a comment)
	}
	return result
}

// TrimWhitespace trims leading/trailing whitespace from each line and removes empty lines.
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
