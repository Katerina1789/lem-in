package validator

import (
	"testing"

	"lem-in/pkg/models"
)

func TestValidateAntCount(t *testing.T) {
	tests := []struct {
		name      string
		lines     []string
		expectErr bool
		errType   string
	}{
		{
			name:      "Valid ant count",
			lines:     []string{"5", "##start"},
			expectErr: false,
		},
		{
			name:      "Ant count = 0",
			lines:     []string{"0"},
			expectErr: true,
			errType:   "ErrAntCountNotPositive",
		},
		{
			name:      "Negative ant count",
			lines:     []string{"-5"},
			expectErr: true,
		},
		{
			name:      "Non-integer",
			lines:     []string{"abc"},
			expectErr: true,
			errType:   "ErrInvalidAntCount",
		},
		{
			name:      "Decimal",
			lines:     []string{"5.5"},
			expectErr: true,
		},
		{
			name:      "Leading zeros",
			lines:     []string{"007"},
			expectErr: true,
			errType:   "ErrAntCountLeadingZeros",
		},
		{
			name:      "Empty lines",
			lines:     []string{},
			expectErr: true,
			errType:   "ErrNoAntCount",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			antCount, idx, err := validateAntCount(tt.lines)
			if tt.expectErr && err == nil {
				t.Errorf("expected error, got nil")
			}
			if !tt.expectErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if !tt.expectErr && antCount <= 0 {
				t.Errorf("invalid ant count: %d", antCount)
			}
			if !tt.expectErr && idx != 1 {
				t.Errorf("expected idx=1, got %d", idx)
			}
		})
	}
}

func TestValidateRooms(t *testing.T) {
	tests := []struct {
		name      string
		lines     []string
		expectErr bool
		errType   string
	}{
		{
			name:      "Valid rooms",
			lines:     []string{"##start", "room1 0 0", "##end", "room2 1 1"},
			expectErr: false,
		},
		{
			name:      "Missing ##start",
			lines:     []string{"##end", "room2 1 1"},
			expectErr: true,
			errType:   "ErrNoStartRoom",
		},
		{
			name:      "Missing ##end",
			lines:     []string{"##start", "room1 0 0"},
			expectErr: true,
			errType:   "ErrNoEndRoom",
		},
		{
			name:      "Multiple ##start",
			lines:     []string{"##start", "room1 0 0", "##start", "room1b 0 1", "##end", "room2 1 1"},
			expectErr: true,
			errType:   "ErrMultipleStartRooms",
		},
		{
			name:      "Multiple ##end",
			lines:     []string{"##start", "room1 0 0", "##end", "room2 1 1", "##end", "room3 2 2"},
			expectErr: true,
			errType:   "ErrMultipleEndRooms",
		},
		{
			name:      "Start and end are same room",
			lines:     []string{"##start", "room1 0 0", "##end", "room1 0 0"},
			expectErr: true,
			errType:   "ErrStartEndSameRoom",
		},
		{
			name:      "Duplicate room name",
			lines:     []string{"##start", "room1 0 0", "room1 1 1", "##end", "room2 2 2"},
			expectErr: true,
			errType:   "ErrDuplicateRoom",
		},
		{
			name:      "Duplicate coordinates",
			lines:     []string{"##start", "room1 0 0", "room2 0 0", "##end", "room3 1 1"},
			expectErr: true,
			errType:   "ErrDuplicateCoordinates",
		},
		{
			name:      "Room name starts with L",
			lines:     []string{"##start", "Lroom 0 0", "##end", "room2 1 1"},
			expectErr: true,
			errType:   "ErrRoomNameStartsWithL",
		},
		{
			name:      "Room name starts with #",
			lines:     []string{"##start", "#room 0 0", "##end", "room2 1 1"},
			expectErr: true,
			errType:   "ErrRoomNameStartsWithHash",
		},
		{
			name:      "Room name contains spaces",
			lines:     []string{"##start", "room 1 0 0", "##end", "room2 1 1"},
			expectErr: true,
			errType:   "ErrInvalidRoomFormat", // 4 tokens instead of 3
		},
		{
			name:      "Invalid room format (2 tokens)",
			lines:     []string{"##start", "room1 0", "##end", "room2 1 1"},
			expectErr: true,
			errType:   "ErrInvalidRoomFormat",
		},
		{
			name:      "Invalid room format (4 tokens)",
			lines:     []string{"##start", "room1 0 0 extra", "##end", "room2 1 1"},
			expectErr: true,
			errType:   "ErrInvalidRoomFormat",
		},
		{
			name:      "Invalid coordinate (non-integer X)",
			lines:     []string{"##start", "room1 abc 0", "##end", "room2 1 1"},
			expectErr: true,
			errType:   "ErrInvalidCoordinates",
		},
		{
			name:      "Invalid coordinate (non-integer Y)",
			lines:     []string{"##start", "room1 0 xyz", "##end", "room2 1 1"},
			expectErr: true,
			errType:   "ErrInvalidCoordinates",
		},
		{
			name:      "##start not followed by room",
			lines:     []string{"##start", "##end", "room2 1 1"},
			expectErr: true,
			errType:   "ErrStartEndNotFollowedByRoom",
		},
		{
			name:      "##end not followed by room",
			lines:     []string{"##start", "room1 0 0", "##end"},
			expectErr: true,
			errType:   "ErrStartEndNotFollowedByRoom",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rooms, _, err := validateRooms(tt.lines, 0)
			if tt.expectErr && err == nil {
				t.Errorf("expected error, got nil")
			}
			if !tt.expectErr {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				hasStart, hasEnd := false, false
				for _, r := range rooms {
					if r.IsStart {
						hasStart = true
					}
					if r.IsEnd {
						hasEnd = true
					}
				}
				if !hasStart || !hasEnd {
					t.Errorf("start or end room not set")
				}
				if len(rooms) < 2 {
					t.Errorf("expected at least 2 rooms, got %d", len(rooms))
				}
			}
		})
	}
}

func TestValidateLinks(t *testing.T) {
	// Helper to build room map
	makeRooms := func(names ...string) map[string]*models.Room {
		rooms := make(map[string]*models.Room)
		for i, name := range names {
			rooms[name] = &models.Room{Name: name, CoordX: i, CoordY: 0}
		}
		return rooms
	}

	tests := []struct {
		name      string
		lines     []string
		rooms     map[string]*models.Room
		startIdx  int
		expectErr bool
		errType   string
	}{
		{
			name:      "Valid link",
			lines:     []string{"room1-room2"},
			rooms:     makeRooms("room1", "room2"),
			startIdx:  0,
			expectErr: false,
		},
		{
			name:      "Unknown room in link",
			lines:     []string{"room1-unknown"},
			rooms:     makeRooms("room1", "room2"),
			startIdx:  0,
			expectErr: true,
			errType:   "ErrUnknownRoomInLink",
		},
		{
			name:      "Self-link",
			lines:     []string{"room1-room1"},
			rooms:     makeRooms("room1", "room2"),
			startIdx:  0,
			expectErr: true,
			errType:   "ErrSelfLink",
		},
		{
			name:      "Duplicate link",
			lines:     []string{"room1-room2", "room2-room1"},
			rooms:     makeRooms("room1", "room2", "room3"),
			startIdx:  0,
			expectErr: true,
			errType:   "ErrDuplicateLink",
		},
		{
			name:      "No links",
			lines:     []string{},
			rooms:     makeRooms("room1", "room2"),
			startIdx:  0,
			expectErr: true,
			errType:   "ErrNoLinks",
		},
		{
			name:      "Link with spaces",
			lines:     []string{"room1 - room2"},
			rooms:     makeRooms("room1", "room2"),
			startIdx:  0,
			expectErr: true,
			errType:   "ErrInvalidLinkFormat",
		},
		{
			name:      "Invalid link format (multiple dashes)",
			lines:     []string{"room1-room2-room3"},
			rooms:     makeRooms("room1", "room2", "room3"),
			startIdx:  0,
			expectErr: true,
			errType:   "ErrInvalidLinkFormat",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			links, err := validateLinks(tt.lines, tt.startIdx, tt.rooms)
			if tt.expectErr && err == nil {
				t.Errorf("expected error, got nil")
			}
			if !tt.expectErr {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if len(links) == 0 {
					t.Errorf("expected links, got none")
				}
			}
		})
	}
}

func TestValidateFullInput(t *testing.T) {
	tests := []struct {
		name      string
		lines     []string
		expectErr bool
	}{
		{
			name: "Valid complete input",
			lines: []string{
				"2",
				"##start",
				"start 0 0",
				"##end",
				"end 5 0",
				"middle 2 0",
				"start-middle",
				"middle-end",
			},
			expectErr: false,
		},
		{
			name: "Missing ant count",
			lines: []string{
				"##start",
				"start 0 0",
				"##end",
				"end 5 0",
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			graph, antCount, err := Validate(tt.lines)
			if tt.expectErr && err == nil {
				t.Errorf("expected error, got nil")
			}
			if !tt.expectErr {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if graph == nil {
					t.Errorf("expected graph, got nil")
				}
				if antCount <= 0 {
					t.Errorf("invalid ant count: %d", antCount)
				}
			}
		})
	}
}
