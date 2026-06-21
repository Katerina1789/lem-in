package validator

import "fmt"

// ValidationError is a structured validation failure with a code and optional message.
type ValidationError struct {
	Code    string
	Message string
}

// Error implements the error interface, formatting code and message together.
func (e *ValidationError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("%s: %s", e.Code, e.Message)
	}
	return e.Code
}

// ErrNoAntCount returns an error when the input contains no ant count line.
func ErrNoAntCount() error {
	return &ValidationError{Code: "ErrNoAntCount", Message: "no ant count line found"}
}

// ErrInvalidAntCount returns an error when the ant count cannot be parsed as an integer.
func ErrInvalidAntCount(reason string) error {
	return &ValidationError{Code: "ErrInvalidAntCount", Message: reason}
}

// ErrAntCountNotPositive returns an error when the ant count is zero or negative.
func ErrAntCountNotPositive() error {
	return &ValidationError{Code: "ErrAntCountNotPositive", Message: "ant count must be positive"}
}

// ErrAntCountLeadingZeros returns an error when the ant count has a leading zero (e.g. "007").
func ErrAntCountLeadingZeros() error {
	return &ValidationError{Code: "ErrAntCountLeadingZeros", Message: "ant count has leading zeros"}
}

// ErrNoStartRoom returns an error when no ##start marker is found.
func ErrNoStartRoom() error {
	return &ValidationError{Code: "ErrNoStartRoom", Message: "no ##start marker found"}
}

// ErrNoEndRoom returns an error when no ##end marker is found.
func ErrNoEndRoom() error {
	return &ValidationError{Code: "ErrNoEndRoom", Message: "no ##end marker found"}
}

// ErrMultipleStartRooms returns an error when more than one ##start marker is found.
func ErrMultipleStartRooms() error {
	return &ValidationError{Code: "ErrMultipleStartRooms", Message: "multiple ##start markers found"}
}

// ErrMultipleEndRooms returns an error when more than one ##end marker is found.
func ErrMultipleEndRooms() error {
	return &ValidationError{Code: "ErrMultipleEndRooms", Message: "multiple ##end markers found"}
}

// ErrStartEndSameRoom returns an error when start and end are the same room.
func ErrStartEndSameRoom() error {
	return &ValidationError{Code: "ErrStartEndSameRoom", Message: "start and end rooms are the same"}
}

// ErrInvalidRoomFormat returns an error when a room line does not have exactly three tokens.
func ErrInvalidRoomFormat() error {
	return &ValidationError{Code: "ErrInvalidRoomFormat", Message: "room must have exactly 3 tokens: name x y"}
}

// ErrInvalidRoomName returns an error when a room name fails a naming rule.
func ErrInvalidRoomName(reason string) error {
	return &ValidationError{Code: "ErrInvalidRoomName", Message: reason}
}

// ErrInvalidCoordinates returns an error when a coordinate value cannot be parsed or overflows.
func ErrInvalidCoordinates(reason string) error {
	return &ValidationError{Code: "ErrInvalidCoordinates", Message: reason}
}

// ErrDuplicateRoom returns an error when the same room name appears more than once.
func ErrDuplicateRoom(name string) error {
	return &ValidationError{Code: "ErrDuplicateRoom", Message: fmt.Sprintf("duplicate room name: %s", name)}
}

// ErrDuplicateCoordinates returns an error when two rooms share the same coordinates.
func ErrDuplicateCoordinates(x, y int) error {
	return &ValidationError{Code: "ErrDuplicateCoordinates", Message: fmt.Sprintf("duplicate coordinates: (%d, %d)", x, y)}
}

// ErrRoomNameStartsWithL returns an error when a room name begins with 'L' (reserved for ant IDs).
func ErrRoomNameStartsWithL() error {
	return &ValidationError{Code: "ErrRoomNameStartsWithL", Message: "room name cannot start with 'L'"}
}

// ErrRoomNameStartsWithHash returns an error when a room name begins with '#'.
func ErrRoomNameStartsWithHash() error {
	return &ValidationError{Code: "ErrRoomNameStartsWithHash", Message: "room name cannot start with '#'"}
}

// ErrRoomNameContainsSpaces returns an error when a room name contains a space character.
func ErrRoomNameContainsSpaces() error {
	return &ValidationError{Code: "ErrRoomNameContainsSpaces", Message: "room name cannot contain spaces"}
}

// ErrStartEndNotFollowedByRoom returns an error when ##start or ##end is not immediately followed by a room line.
func ErrStartEndNotFollowedByRoom() error {
	return &ValidationError{Code: "ErrStartEndNotFollowedByRoom", Message: "##start or ##end not followed by a room line"}
}

// ErrUnknownRoomInLink returns an error when a link references a room that was not declared.
func ErrUnknownRoomInLink(roomName string) error {
	return &ValidationError{Code: "ErrUnknownRoomInLink", Message: fmt.Sprintf("link references unknown room: %s", roomName)}
}

// ErrDuplicateLink returns an error when the same link appears more than once.
func ErrDuplicateLink(room1, room2 string) error {
	return &ValidationError{Code: "ErrDuplicateLink", Message: fmt.Sprintf("duplicate link: %s-%s", room1, room2)}
}

// ErrSelfLink returns an error when a link connects a room to itself.
func ErrSelfLink(roomName string) error {
	return &ValidationError{Code: "ErrSelfLink", Message: fmt.Sprintf("self-referencing link: %s-%s", roomName, roomName)}
}

// ErrInvalidLinkFormat returns an error when a link line does not match the "room1-room2" format.
func ErrInvalidLinkFormat() error {
	return &ValidationError{Code: "ErrInvalidLinkFormat", Message: "link must have format: room1-room2"}
}

// ErrNoLinks returns an error when the input defines no links at all.
func ErrNoLinks() error {
	return &ValidationError{Code: "ErrNoLinks", Message: "no links defined"}
}
