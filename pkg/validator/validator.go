package validator

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"lem-in/pkg/models"
)

// Validate processes parsed lines and returns a populated Graph, ant count, or an error.
func Validate(lines []string) (*models.Graph, int, error) {
	antCount, lineIdx, err := validateAntCount(lines)
	if err != nil {
		return nil, 0, err
	}

	rooms, lineIdx, err := validateRooms(lines, lineIdx)
	if err != nil {
		return nil, 0, err
	}

	links, err := validateLinks(lines, lineIdx, rooms)
	if err != nil {
		return nil, 0, err
	}

	graph := buildGraph(rooms, links)
	return graph, antCount, nil
}

// validateAntCount extracts and validates the ant count from the first line.
func validateAntCount(lines []string) (int, int, error) {
	if len(lines) == 0 {
		return 0, 0, ErrNoAntCount()
	}

	antLine := lines[0]

	antCount, err := strconv.Atoi(antLine)
	if err != nil {
		return 0, 0, ErrInvalidAntCount(fmt.Sprintf("not an integer: %s", antLine))
	}

	if antCount <= 0 {
		return 0, 0, ErrAntCountNotPositive()
	}

	if len(antLine) > 1 && antLine[0] == '0' { // e.g. "007" is invalid
		return 0, 0, ErrAntCountLeadingZeros()
	}

	return antCount, 1, nil
}

// validateRooms parses room definitions and ##start/##end markers from lines[startIdx:].
func validateRooms(lines []string, startIdx int) (map[string]*models.Room, int, error) {
	rooms := make(map[string]*models.Room)
	coords := make(map[[2]int]bool)
	var startName, endName string
	var startCount, endCount int

	i := startIdx
	for i < len(lines) {
		line := lines[i]

		if strings.Contains(line, "-") && !strings.Contains(line, " ") { // reached the link section
			break
		}

		if line == "##start" {
			startCount++
			if startCount > 1 {
				return nil, 0, ErrMultipleStartRooms()
			}
			i++
			if i >= len(lines) {
				return nil, 0, ErrStartEndNotFollowedByRoom()
			}
			roomLine := lines[i]
			if strings.HasPrefix(roomLine, "##") || strings.Contains(roomLine, "-") {
				return nil, 0, ErrStartEndNotFollowedByRoom()
			}
			name, x, y, err := parseRoomLine(roomLine)
			if err != nil {
				return nil, 0, err
			}
			if _, exists := rooms[name]; exists {
				return nil, 0, ErrDuplicateRoom(name)
			}
			if coords[[2]int{x, y}] {
				return nil, 0, ErrDuplicateCoordinates(x, y)
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
				return nil, 0, ErrMultipleEndRooms()
			}
			i++
			if i >= len(lines) {
				return nil, 0, ErrStartEndNotFollowedByRoom()
			}
			roomLine := lines[i]
			if strings.HasPrefix(roomLine, "##") || strings.Contains(roomLine, "-") {
				return nil, 0, ErrStartEndNotFollowedByRoom()
			}
			name, x, y, err := parseRoomLine(roomLine)
			if err != nil {
				return nil, 0, err
			}
			if _, exists := rooms[name]; exists {
				return nil, 0, ErrDuplicateRoom(name)
			}
			if coords[[2]int{x, y}] {
				return nil, 0, ErrDuplicateCoordinates(x, y)
			}
			endName = name
			rooms[name] = &models.Room{Name: name, CoordX: x, CoordY: y, IsEnd: true}
			coords[[2]int{x, y}] = true
			i++
			continue
		}

		name, x, y, err := parseRoomLine(line)
		if err != nil {
			return nil, 0, err
		}
		if _, exists := rooms[name]; exists {
			return nil, 0, ErrDuplicateRoom(name)
		}
		if coords[[2]int{x, y}] {
			return nil, 0, ErrDuplicateCoordinates(x, y)
		}
		rooms[name] = &models.Room{Name: name, CoordX: x, CoordY: y}
		coords[[2]int{x, y}] = true
		i++
	}

	if startCount == 0 {
		return nil, 0, ErrNoStartRoom()
	}
	if endCount == 0 {
		return nil, 0, ErrNoEndRoom()
	}
	if startName == endName {
		return nil, 0, ErrStartEndSameRoom()
	}

	return rooms, i, nil
}

// parseRoomLine parses a single "name x y" room line and validates its fields.
func parseRoomLine(line string) (string, int, int, error) {
	parts := strings.Fields(line)
	if len(parts) != 3 {
		return "", 0, 0, ErrInvalidRoomFormat()
	}

	name := parts[0]

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

// validateLinks parses and validates all link lines from lines[startIdx:].
func validateLinks(lines []string, startIdx int, rooms map[string]*models.Room) ([]models.Link, error) {
	var links []models.Link
	seenLinks := make(map[[2]string]bool)

	for i := startIdx; i < len(lines); i++ {
		line := lines[i]

		if strings.HasPrefix(line, "##") { // skip ##start / ##end directives in the link section
			continue
		}

		if !strings.Contains(line, "-") {
			continue
		}

		if strings.Contains(line, " ") { // links must not contain spaces
			return nil, ErrInvalidLinkFormat()
		}

		parts := strings.Split(line, "-")
		if len(parts) != 2 {
			return nil, ErrInvalidLinkFormat()
		}

		room1, room2 := parts[0], parts[1]

		if room1 == "" || room2 == "" {
			return nil, ErrInvalidLinkFormat()
		}

		if room1 == room2 {
			return nil, ErrSelfLink(room1)
		}

		if _, exists := rooms[room1]; !exists {
			return nil, ErrUnknownRoomInLink(room1)
		}
		if _, exists := rooms[room2]; !exists {
			return nil, ErrUnknownRoomInLink(room2)
		}

		key1 := [2]string{room1, room2}
		key2 := [2]string{room2, room1}
		if seenLinks[key1] || seenLinks[key2] { // normalize order to detect both directions
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

// buildGraph constructs a Graph from the validated rooms and links.
func buildGraph(rooms map[string]*models.Room, links []models.Link) *models.Graph {
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
