package parser

import (
	"testing"
)

func TestReadFile(t *testing.T) {
	// This test would need a temp file, skipping detailed test for now
	// In real scenario, create a temp file and test
}

func TestSplitLines(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "Unix line endings",
			input:    "line1\nline2\nline3",
			expected: []string{"line1", "line2", "line3"},
		},
		{
			name:     "Windows line endings",
			input:    "line1\r\nline2\r\nline3",
			expected: []string{"line1", "line2", "line3"},
		},
		{
			name:     "Mixed line endings",
			input:    "line1\r\nline2\nline3",
			expected: []string{"line1", "line2", "line3"},
		},
		{
			name:     "Trailing newline",
			input:    "line1\nline2\n",
			expected: []string{"line1", "line2", ""},
		},
		{
			name:     "Empty string",
			input:    "",
			expected: []string{""},
		},
		{
			name:     "Single line",
			input:    "single",
			expected: []string{"single"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SplitLines(tt.input)
			if len(got) != len(tt.expected) {
				t.Fatalf("expected %d lines, got %d", len(tt.expected), len(got))
			}
			for i, line := range got {
				if line != tt.expected[i] {
					t.Errorf("line %d: expected %q, got %q", i, tt.expected[i], line)
				}
			}
		})
	}
}

func TestFilterComments(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name:     "Keep ##start and ##end",
			input:    []string{"##start", "room1", "##end"},
			expected: []string{"##start", "room1", "##end"},
		},
		{
			name:     "Remove # comments",
			input:    []string{"#comment", "room1", "# another comment"},
			expected: []string{"room1"},
		},
		{
			name:     "Remove unknown ## commands",
			input:    []string{"##unknown", "room1", "##other"},
			expected: []string{"room1"},
		},
		{
			name:     "Mixed",
			input:    []string{"#comment", "##start", "room1", "##unknown", "#another", "##end"},
			expected: []string{"##start", "room1", "##end"},
		},
		{
			name:     "Empty input",
			input:    []string{},
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FilterComments(tt.input)
			if len(got) != len(tt.expected) {
				t.Fatalf("expected %d items, got %d", len(tt.expected), len(got))
			}
			for i, item := range got {
				if item != tt.expected[i] {
					t.Errorf("item %d: expected %q, got %q", i, tt.expected[i], item)
				}
			}
		})
	}
}

func TestTrimWhitespace(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name:     "Remove leading/trailing spaces",
			input:    []string{"  room1  ", "  room2  "},
			expected: []string{"room1", "room2"},
		},
		{
			name:     "Remove empty lines",
			input:    []string{"room1", "", "room2", "  ", "room3"},
			expected: []string{"room1", "room2", "room3"},
		},
		{
			name:     "Preserve internal spaces",
			input:    []string{"  a b c  "},
			expected: []string{"a b c"},
		},
		{
			name:     "Empty input",
			input:    []string{},
			expected: []string{},
		},
		{
			name:     "All empty",
			input:    []string{"", "  ", "\t"},
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TrimWhitespace(tt.input)
			if len(got) != len(tt.expected) {
				t.Fatalf("expected %d items, got %d", len(tt.expected), len(got))
			}
			for i, item := range got {
				if item != tt.expected[i] {
					t.Errorf("item %d: expected %q, got %q", i, tt.expected[i], item)
				}
			}
		})
	}
}
