package parser

import (
	"os"
	"strings"
)

// ReadFile reads the entire file at path and returns its contents as a string.
func ReadFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// SplitLines splits raw input by newlines and strips carriage returns from each line.
func SplitLines(raw string) []string {
	lines := strings.Split(raw, "\n")
	for i := range lines {
		lines[i] = strings.TrimSuffix(lines[i], "\r") // strip \r for Windows line endings
	}
	return lines
}

// FilterComments removes comment lines (starting with '#') while keeping ##start and ##end directives.
func FilterComments(lines []string) []string {
	var result []string
	for _, line := range lines {
		trimmed := strings.TrimLeft(line, " \t")
		if !strings.HasPrefix(trimmed, "#") { // not a comment, keep it
			result = append(result, line)
			continue
		}
		if trimmed == "##start" || trimmed == "##end" { // directives are not comments
			result = append(result, line)
		}
	}
	return result
}

// TrimWhitespace trims leading/trailing whitespace from each line and removes blank lines.
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
