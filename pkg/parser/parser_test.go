package parser

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadFile(t *testing.T) {
	// Test valid file
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	content := "line1\nline2\nline3"
	err := os.WriteFile(testFile, []byte(content), 0o644)
	if err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	result, err := ReadFile(testFile)
	if err != nil {
		t.Fatalf("ReadFile failed: %v", err)
	}
	if result != content {
		t.Errorf("Expected %q, got %q", content, result)
	}

	// Test non-existent file
	_, err = ReadFile(filepath.Join(tmpDir, "nonexistent.txt"))
	if err == nil {
		t.Error("Expected error for non-existent file")
	}

	// Test empty file
	emptyFile := filepath.Join(tmpDir, "empty.txt")
	err = os.WriteFile(emptyFile, []byte(""), 0o644)
	if err != nil {
		t.Fatalf("Failed to write empty file: %v", err)
	}
	result, err = ReadFile(emptyFile)
	if err != nil {
		t.Fatalf("ReadFile failed: %v", err)
	}
	if result != "" {
		t.Errorf("Expected empty string, got %q", result)
	}
}

func TestSplitLines(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "basic lines with LF",
			input:    "line1\nline2\nline3",
			expected: []string{"line1", "line2", "line3"},
		},
		{
			name:     "lines with CRLF",
			input:    "line1\r\nline2\r\nline3",
			expected: []string{"line1", "line2", "line3"},
		},
		{
			name:     "trailing newline",
			input:    "line1\nline2\n",
			expected: []string{"line1", "line2", ""},
		},
		{
			name:     "trailing CRLF",
			input:    "line1\r\nline2\r\n",
			expected: []string{"line1", "line2", ""},
		},
		{
			name:     "empty string",
			input:    "",
			expected: []string{""},
		},
		{
			name:     "single line",
			input:    "hello",
			expected: []string{"hello"},
		},
		{
			name:     "mixed line endings",
			input:    "line1\r\nline2\nline3\r\n",
			expected: []string{"line1", "line2", "line3", ""},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SplitLines(tt.input)
			if len(result) != len(tt.expected) {
				t.Errorf("Length mismatch: expected %d, got %d", len(tt.expected), len(result))
				return
			}
			for i, line := range result {
				if line != tt.expected[i] {
					t.Errorf("Line %d: expected %q, got %q", i, tt.expected[i], line)
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
			name:     "no comments",
			input:    []string{"room1", "room2", "room3"},
			expected: []string{"room1", "room2", "room3"},
		},
		{
			name:     "single hash comments",
			input:    []string{"room1", "# this is a comment", "room2"},
			expected: []string{"room1", "room2"},
		},
		{
			name:     "keep ##start",
			input:    []string{"##start", "room1"},
			expected: []string{"##start", "room1"},
		},
		{
			name:     "keep ##end",
			input:    []string{"room1", "##end"},
			expected: []string{"room1", "##end"},
		},
		{
			name:     "filter other ## lines",
			input:    []string{"room1", "##other", "room2"},
			expected: []string{"room1", "room2"},
		},
		{
			name:     "comment with leading whitespace",
			input:    []string{"room1", "  # comment", "room2"},
			expected: []string{"room1", "room2"},
		},
		{
			name:     "##start with leading whitespace",
			input:    []string{"  ##start", "room1"},
			expected: []string{"  ##start", "room1"},
		},
		{
			name:     "##end with leading whitespace",
			input:    []string{"room1", "  ##end"},
			expected: []string{"room1", "  ##end"},
		},
		{
			name:     "mixed comments and directives",
			input:    []string{"##start", "# comment1", "room1", "room2", "# comment2", "##end"},
			expected: []string{"##start", "room1", "room2", "##end"},
		},
		{
			name:     "only comments",
			input:    []string{"# comment1", "# comment2"},
			expected: []string{},
		},
		{
			name:     "empty input",
			input:    []string{},
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FilterComments(tt.input)
			if len(result) != len(tt.expected) {
				t.Errorf("Length mismatch: expected %d, got %d", len(tt.expected), len(result))
				return
			}
			for i, line := range result {
				if line != tt.expected[i] {
					t.Errorf("Line %d: expected %q, got %q", i, tt.expected[i], line)
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
			name:     "no whitespace",
			input:    []string{"room1", "room2"},
			expected: []string{"room1", "room2"},
		},
		{
			name:     "leading/trailing spaces",
			input:    []string{"  room1  ", "   room2   "},
			expected: []string{"room1", "room2"},
		},
		{
			name:     "empty lines removed",
			input:    []string{"room1", "", "room2", "   ", "room3"},
			expected: []string{"room1", "room2", "room3"},
		},
		{
			name:     "tabs and spaces",
			input:    []string{"\t\troom1\t\t", "  room2  "},
			expected: []string{"room1", "room2"},
		},
		{
			name:     "mixed whitespace and content",
			input:    []string{"  ##start  ", "room1", "  ", "room2", "  ##end  "},
			expected: []string{"##start", "room1", "room2", "##end"},
		},
		{
			name:     "only empty lines",
			input:    []string{"", "   ", "\t"},
			expected: []string{},
		},
		{
			name:     "empty input",
			input:    []string{},
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := TrimWhitespace(tt.input)
			if len(result) != len(tt.expected) {
				t.Errorf("Length mismatch: expected %d, got %d", len(tt.expected), len(result))
				return
			}
			for i, line := range result {
				if line != tt.expected[i] {
					t.Errorf("Line %d: expected %q, got %q", i, tt.expected[i], line)
				}
			}
		})
	}
}

// Integration test: full pipeline
func TestFullPipeline(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	content := `##start
room1 0 0
# This is a comment
room2 1 1
  # Another comment with spaces
##end
room3 2 2
  
# Final comment
`
	err := os.WriteFile(testFile, []byte(content), 0o644)
	if err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	// Pipeline: ReadFile -> SplitLines -> FilterComments -> TrimWhitespace
	raw, err := ReadFile(testFile)
	if err != nil {
		t.Fatalf("ReadFile failed: %v", err)
	}

	lines := SplitLines(raw)
	lines = FilterComments(lines)
	lines = TrimWhitespace(lines)

	expected := []string{"##start", "room1 0 0", "room2 1 1", "##end", "room3 2 2"}
	if len(lines) != len(expected) {
		t.Errorf("Length mismatch: expected %d, got %d", len(expected), len(lines))
		return
	}
	for i, line := range lines {
		if line != expected[i] {
			t.Errorf("Line %d: expected %q, got %q", i, expected[i], line)
		}
	}
}
