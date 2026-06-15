package validator

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"lem-in/pkg/models"
)

// Validate processes parsed lines and returns a populated Graph, ant count, or error
func Validate(lines []string) (*models.Graph, int, error) {
	// Stage 1: Validate ant count
	antCount, lineIdx, err := validateAntCount(lines)
	if err != nil {
		return nil, 0, err
	}

	// Stage 2: Validate rooms
	rooms, startName, endName, lineIdx, err := validateRooms(lines, lineIdx)
	if err != nil {
		return nil, 0, err
	}

	// Stage 3: Validate links
	links, err := validateLinks(lines, lineIdx, rooms)
	if err != nil {
		return nil, 0, err
	}

	// Build graph
	graph := buildGraph(rooms, links, startName, endName)

	return graph, antCount, nil
}

// Stage 1: Validate and extract ant count
func validateAntCount(lines []string) (int, int, error) {
	if len(lines) == 0 {
		return 0, 0, ErrNoAntCount()
	}

	// Find first non-empty line (should be ant count)
	antLine := lines[0]

	// Must be a valid integer
	antCount, err := strconv.Atoi(antLine)
	if err != nil {
		return 0, 0, ErrInvalidAntCount(fmt.Sprintf("not an integer: %s", antLine))
	}

	// Must be positive
	if antCount <= 0 {
		return 0, 0, ErrAntCountNotPositive()
	}

	// Check for leading zeros (e.g., "007" is invalid, but "0" is valid)
	if len(antLine) > 1 && antLine[0] == '0' {
		return 0, 0, ErrAntCountLeadingZeros()
	}

	return antCount, 1, nil
}

// Stage 2: Validate and extract rooms
func validateRooms(lines []string, startIdx int) (map[string]*models.Room, string, string, int, error) {
	rooms := make(map[string]*models.Room)
	coords := make(map[[2]int]bool) // track (x, y) pairs
	var startName, endName string
	var startCount, endCount int

	i := startIdx
	for i < len(lines) {
		line := lines[i]

		// Check for link section (contains exactly one dash)
		if strings.Contains(line, "-") && !strings.Contains(line, " ") {
			// This is a link, exit room section
			break
		}

		// Handle ##start and ##end markers
		if line == "##start" {
			startCount++
			if startCount > 1 {
				return nil, "", "", 0, ErrMultipleStartRooms()
			}
			i++
			if i >= len(lines) {
				return nil, "", "", 0, ErrStartEndNotFollowedByRoom()
			}
			// Next line should be the start room
			roomLine := lines[i]
			if strings.HasPrefix(roomLine, "##") || strings.Contains(roomLine, "-") {
				return nil, "", "", 0, ErrStartEndNotFollowedByRoom()
			}
			name, x, y, err := parseRoomLine(roomLine)
			if err != nil {
				return nil, "", "", 0, err
			}
			if _, exists := rooms[name]; exists {
				return nil, "", "", 0, ErrDuplicateRoom(name)
			}
			if coords[[2]int{x, y}] {
				return nil, "", "", 0, ErrDuplicateCoordinates(x, y)
			}
			startName = name
			rooms[name] = &models.Room{Name: name, CoordX: x, CoordY: y, IsStart: true}
			coords[[2]int{x, y}] = true
			i++
			continue
		}

		if line == "##end" {
			endCount++
			if endCount > 1 {
				return nil, "", "", 0, ErrMultipleEndRooms()
			}
			i++
			if i >= len(lines) {
				return nil, "", "", 0, ErrStartEndNotFollowedByRoom()
			}
			// Next line should be the end room
			roomLine := lines[i]
			if strings.HasPrefix(roomLine, "##") || strings.Contains(roomLine, "-") {
				return nil, "", "", 0, ErrStartEndNotFollowedByRoom()
			}
			name, x, y, err := parseRoomLine(roomLine)
			if err != nil {
				return nil, "", "", 0, err
			}
			if _, exists := rooms[name]; exists {
				return nil, "", "", 0, ErrDuplicateRoom(name)
			}
			if coords[[2]int{x, y}] {
				return nil, "", "", 0, ErrDuplicateCoordinates(x, y)
			}
			endName = name
			rooms[name] = &models.Room{Name: name, CoordX: x, CoordY: y, IsEnd: true}
			coords[[2]int{x, y}] = true
			i++
			continue
		}

		// Regular room line
		name, x, y, err := parseRoomLine(line)
		if err != nil {
			return nil, "", "", 0, err
		}
		if _, exists := rooms[name]; exists {
			return nil, "", "", 0, ErrDuplicateRoom(name)
		}
		if coords[[2]int{x, y}] {
			return nil, "", "", 0, ErrDuplicateCoordinates(x, y)
		}
		rooms[name] = &models.Room{Name: name, CoordX: x, CoordY: y}
		coords[[2]int{x, y}] = true
		i++
	}

	// Validate start and end exist
	if startCount == 0 {
		return nil, "", "", 0, ErrNoStartRoom()
	}
	if endCount == 0 {
		return nil, "", "", 0, ErrNoEndRoom()
	}

	// Validate start != end
	if startName == endName {
		return nil, "", "", 0, ErrStartEndSameRoom()
	}

	return rooms, startName, endName, i, nil
}

// Parse a single room line: "name x y"
func parseRoomLine(line string) (string, int, int, error) {
	parts := strings.Fields(line)
	if len(parts) != 3 {
		return "", 0, 0, ErrInvalidRoomFormat()
	}

	name := parts[0]

	// Validate room name
	if name == "" {
		return "", 0, 0, ErrInvalidRoomName("room name is empty")
	}
	if strings.HasPrefix(name, "L") {
		return "", 0, 0, ErrRoomNameStartsWithL()
	}
	if strings.HasPrefix(name, "#") {
		return "", 0, 0, ErrRoomNameStartsWithHash()
	}
	if strings.Contains(name, " ") {
		return "", 0, 0, ErrRoomNameContainsSpaces()
	}

	// Parse coordinates
	x, err := strconv.ParseInt(parts[1], 10, 32)
	if err != nil {
		return "", 0, 0, ErrInvalidCoordinates(fmt.Sprintf("x is not a valid integer: %s", parts[1]))
	}
	if x > math.MaxInt32 || x < math.MinInt32 {
		return "", 0, 0, ErrInvalidCoordinates(fmt.Sprintf("x overflow: %d", x))
	}

	y, err := strconv.ParseInt(parts[2], 10, 32)
	if err != nil {
		return "", 0, 0, ErrInvalidCoordinates(fmt.Sprintf("y is not a valid integer: %s", parts[2]))
	}
	if y > math.MaxInt32 || y < math.MinInt32 {
		return "", 0, 0, ErrInvalidCoordinates(fmt.Sprintf("y overflow: %d", y))
	}

	return name, int(x), int(y), nil
}

// Stage 3: Validate and extract links
func validateLinks(lines []string, startIdx int, rooms map[string]*models.Room) ([]models.Link, error) {
	var links []models.Link
	seenLinks := make(map[[2]string]bool)

	for i := startIdx; i < len(lines); i++ {
		line := lines[i]

		// Skip command lines (##start, ##end)
		if strings.HasPrefix(line, "##") {
			continue
		}

		// Check if this is a link
		if !strings.Contains(line, "-") {
			// Not a link, skip
			continue
		}

		// Links must not contain spaces
		if strings.Contains(line, " ") {
			return nil, ErrInvalidLinkFormat()
		}

		// Split on the single dash
		parts := strings.Split(line, "-")
		if len(parts) != 2 {
			return nil, ErrInvalidLinkFormat()
		}

		room1, room2 := parts[0], parts[1]

		if room1 == "" || room2 == "" {
			return nil, ErrInvalidLinkFormat()
		}

		// Check for self-link
		if room1 == room2 {
			return nil, ErrSelfLink(room1)
		}

		// Check that both rooms exist
		if _, exists := rooms[room1]; !exists {
			return nil, ErrUnknownRoomInLink(room1)
		}
		if _, exists := rooms[room2]; !exists {
			return nil, ErrUnknownRoomInLink(room2)
		}

		// Check for duplicates (normalize so smaller name comes first)
		key1 := [2]string{room1, room2}
		key2 := [2]string{room2, room1}
		if seenLinks[key1] || seenLinks[key2] {
			return nil, ErrDuplicateLink(room1, room2)
		}
		seenLinks[key1] = true

		links = append(links, models.Link{Room1: room1, Room2: room2})
	}

	if len(links) == 0 {
		return nil, ErrNoLinks()
	}

	return links, nil
}

// Build the final Graph struct
func buildGraph(rooms map[string]*models.Room, links []models.Link, startName, endName string) *models.Graph {
	adjacency := make(map[string][]string)

	for _, link := range links {
		adjacency[link.Room1] = append(adjacency[link.Room1], link.Room2)
		adjacency[link.Room2] = append(adjacency[link.Room2], link.Room1)
	}

	return &models.Graph{
		Rooms:     rooms,
		Links:     links,
		Adjacency: adjacency,
	}
}
