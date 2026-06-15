package validator

import "fmt"

// Custom error types for validation failures

type ValidationError struct {
	Code    string
	Message string
}

func (e *ValidationError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("%s: %s", e.Code, e.Message)
	}
	return e.Code
}

// Error constructors

func ErrNoAntCount() error {
	return &ValidationError{Code: "ErrNoAntCount", Message: "no ant count line found"}
}

func ErrInvalidAntCount(reason string) error {
	return &ValidationError{Code: "ErrInvalidAntCount", Message: reason}
}

func ErrAntCountNotPositive() error {
	return &ValidationError{Code: "ErrAntCountNotPositive", Message: "ant count must be positive"}
}

func ErrAntCountLeadingZeros() error {
	return &ValidationError{Code: "ErrAntCountLeadingZeros", Message: "ant count has leading zeros"}
}

func ErrNoStartRoom() error {
	return &ValidationError{Code: "ErrNoStartRoom", Message: "no ##start marker found"}
}

func ErrNoEndRoom() error {
	return &ValidationError{Code: "ErrNoEndRoom", Message: "no ##end marker found"}
}

func ErrMultipleStartRooms() error {
	return &ValidationError{Code: "ErrMultipleStartRooms", Message: "multiple ##start markers found"}
}

func ErrMultipleEndRooms() error {
	return &ValidationError{Code: "ErrMultipleEndRooms", Message: "multiple ##end markers found"}
}

func ErrStartEndSameRoom() error {
	return &ValidationError{Code: "ErrStartEndSameRoom", Message: "start and end rooms are the same"}
}

func ErrInvalidRoomFormat() error {
	return &ValidationError{Code: "ErrInvalidRoomFormat", Message: "room must have exactly 3 tokens: name x y"}
}

func ErrInvalidRoomName(reason string) error {
	return &ValidationError{Code: "ErrInvalidRoomName", Message: reason}
}

func ErrInvalidCoordinates(reason string) error {
	return &ValidationError{Code: "ErrInvalidCoordinates", Message: reason}
}

func ErrDuplicateRoom(name string) error {
	return &ValidationError{Code: "ErrDuplicateRoom", Message: fmt.Sprintf("duplicate room name: %s", name)}
}

func ErrDuplicateCoordinates(x, y int) error {
	return &ValidationError{Code: "ErrDuplicateCoordinates", Message: fmt.Sprintf("duplicate coordinates: (%d, %d)", x, y)}
}

func ErrRoomNameStartsWithL() error {
	return &ValidationError{Code: "ErrRoomNameStartsWithL", Message: "room name cannot start with 'L'"}
}

func ErrRoomNameStartsWithHash() error {
	return &ValidationError{Code: "ErrRoomNameStartsWithHash", Message: "room name cannot start with '#'"}
}

func ErrRoomNameContainsSpaces() error {
	return &ValidationError{Code: "ErrRoomNameContainsSpaces", Message: "room name cannot contain spaces"}
}

func ErrStartEndNotFollowedByRoom() error {
	return &ValidationError{Code: "ErrStartEndNotFollowedByRoom", Message: "##start or ##end not followed by a room line"}
}

func ErrUnknownRoomInLink(roomName string) error {
	return &ValidationError{Code: "ErrUnknownRoomInLink", Message: fmt.Sprintf("link references unknown room: %s", roomName)}
}

func ErrDuplicateLink(room1, room2 string) error {
	return &ValidationError{Code: "ErrDuplicateLink", Message: fmt.Sprintf("duplicate link: %s-%s", room1, room2)}
}

func ErrSelfLink(roomName string) error {
	return &ValidationError{Code: "ErrSelfLink", Message: fmt.Sprintf("self-referencing link: %s-%s", roomName, roomName)}
}

func ErrInvalidLinkFormat() error {
	return &ValidationError{Code: "ErrInvalidLinkFormat", Message: "link must have format: room1-room2"}
}

func ErrNoLinks() error {
	return &ValidationError{Code: "ErrNoLinks", Message: "no links defined"}
}
