package output

import (
	"strings"
	"testing"

	"lem-in/pkg/models"
)

func TestFormatFileEcho(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "normalizes trailing newlines to one",
			input: "4\n##start\n0 0 3\n##end\n1 1 1\n0-1\n\n\n",
			want:  "4\n##start\n0 0 3\n##end\n1 1 1\n0-1\n",
		},
		{
			name:  "adds newline when missing",
			input: "4\n##start\n0 0 3\n##end\n1 1 1\n0-1",
			want:  "4\n##start\n0 0 3\n##end\n1 1 1\n0-1\n",
		},
		{
			name:  "preserves single trailing newline",
			input: "2\nstart 0 0\nend 1 1\nstart-end\n",
			want:  "2\nstart 0 0\nend 1 1\nstart-end\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatFileEcho(tt.input)
			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestFormatTurns(t *testing.T) {
	tests := []struct {
		name  string
		turns []models.Turn
		want  string
	}{
		{
			name: "single turn single ant",
			turns: []models.Turn{
				{{AntID: 1, Room: "room1"}},
			},
			want: "L1-room1\n",
		},
		{
			name: "single turn multiple ants",
			turns: []models.Turn{
				{{AntID: 1, Room: "a"}, {AntID: 2, Room: "b"}},
			},
			want: "L1-a L2-b\n",
		},
		{
			name: "multiple turns",
			turns: []models.Turn{
				{{AntID: 1, Room: "r1"}, {AntID: 2, Room: "r2"}},
				{{AntID: 1, Room: "end"}, {AntID: 3, Room: "r1"}},
			},
			want: "L1-r1 L2-r2\nL1-end L3-r1\n",
		},
		{
			name:  "empty turns",
			turns: []models.Turn{},
			want:  "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatTurns(tt.turns)
			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestFormatFull(t *testing.T) {
	input := "2\n##start\na 0 0\n##end\nb 1 1\na-b\n"
	turns := []models.Turn{
		{{AntID: 1, Room: "b"}},
		{{AntID: 2, Room: "b"}},
	}

	want := "2\n##start\na 0 0\n##end\nb 1 1\na-b\n\nL1-b\nL2-b\n"

	got := FormatFull(input, turns)
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFormatFull_NoTrailingSpace(t *testing.T) {
	input := "1\n##start\na 0 0\n##end\nb 1 1\na-b\n"
	turns := []models.Turn{
		{{AntID: 1, Room: "b"}, {AntID: 2, Room: "c"}},
	}

	got := FormatFull(input, turns)
	lines := strings.Split(strings.TrimRight(got, "\n"), "\n")
	for _, line := range lines {
		if len(line) > 0 && line[len(line)-1] == ' ' {
			t.Errorf("line has trailing space: %q", line)
		}
	}
}
