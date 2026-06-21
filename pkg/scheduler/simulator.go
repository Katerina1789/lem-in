package scheduler

import (
	"lem-in/pkg/models"
	"sort"
)

// Simulate runs the ant movement simulation turn by turn until all ants reach the end room.
func Simulate(ants []*models.Ant, graph *models.Graph) []models.Turn {
	var turns []models.Turn

	for {
		if allAntsArrived(ants) {
			break
		}

		currentTurn := models.Turn{}
		roomOccupied := getOccupiedRooms(ants)
		pathList := getUniquePaths(ants)

		for _, path := range pathList {
			pathAnts := getSortedAntsOnPath(ants, path)

			activeMoves := moveActiveAnts(pathAnts, roomOccupied)
			currentTurn = append(currentTurn, activeMoves...)

			if launchMove := launchNewAnt(pathAnts, path, roomOccupied); launchMove != nil {
				currentTurn = append(currentTurn, *launchMove)
			}
		}

		if len(currentTurn) > 0 {
			sort.Slice(currentTurn, func(i, j int) bool {
				return currentTurn[i].AntID < currentTurn[j].AntID
			})
			turns = append(turns, currentTurn)
		} else {
			break
		}
	}

	return turns
}

// equalPaths reports whether two path slices contain the same rooms in the same order.
func equalPaths(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// allAntsArrived reports whether every ant has reached the end room.
func allAntsArrived(ants []*models.Ant) bool {
	for _, ant := range ants {
		if !ant.Arrived {
			return false
		}
	}
	return true
}

// getOccupiedRooms returns a set of room names currently occupied by ants that are
// in transit (not at the end room and not yet launched).
func getOccupiedRooms(ants []*models.Ant) map[string]bool {
	roomOccupied := make(map[string]bool)
	for _, ant := range ants {
		if !ant.Arrived && ant.Position >= 0 {
			currRoom := ant.Path[ant.Position]
			if currRoom != ant.Path[len(ant.Path)-1] { // exclude ants sitting at the end room
				roomOccupied[currRoom] = true
			}
		}
	}
	return roomOccupied
}

// getUniquePaths returns one entry per distinct path found among the given ants.
func getUniquePaths(ants []*models.Ant) [][]string {
	var pathList [][]string
	for _, ant := range ants {
		found := false
		for _, p := range pathList {
			if equalPaths(p, ant.Path) {
				found = true
				break
			}
		}
		if !found {
			pathList = append(pathList, ant.Path)
		}
	}
	return pathList
}

// getSortedAntsOnPath returns all ants assigned to path, sorted front-to-back by position
// so that leading ants are moved first and cannot be blocked by ants behind them.
func getSortedAntsOnPath(ants []*models.Ant, path []string) []*models.Ant {
	var pathAnts []*models.Ant
	for _, ant := range ants {
		if equalPaths(ant.Path, path) {
			pathAnts = append(pathAnts, ant)
		}
	}
	sort.Slice(pathAnts, func(i, j int) bool {
		return pathAnts[i].Position > pathAnts[j].Position // higher position = closer to end
	})
	return pathAnts
}

// moveActiveAnts advances each in-transit ant one step along its path if the next room is free.
func moveActiveAnts(pathAnts []*models.Ant, roomOccupied map[string]bool) []models.TurnMove {
	var moves []models.TurnMove
	for i := 0; i < len(pathAnts); i++ {
		ant := pathAnts[i]
		if ant.Arrived || ant.Position == -1 {
			continue
		}

		nextPos := ant.Position + 1
		nextRoom := ant.Path[nextPos]
		isEnd := nextPos == len(ant.Path)-1

		if isEnd || !roomOccupied[nextRoom] {
			currRoom := ant.Path[ant.Position]
			delete(roomOccupied, currRoom)

			ant.Position = nextPos
			if isEnd {
				ant.Arrived = true
			} else {
				roomOccupied[nextRoom] = true
			}

			moves = append(moves, models.TurnMove{AntID: ant.ID, Room: nextRoom})
		}
	}
	return moves
}

// launchNewAnt starts the first waiting ant on path if the path's first room is free.
func launchNewAnt(pathAnts []*models.Ant, path []string, roomOccupied map[string]bool) *models.TurnMove {
	firstRoom := path[0]
	isFirstRoomEnd := len(path) == 1 // path goes directly start→end
	if isFirstRoomEnd || !roomOccupied[firstRoom] {
		var inactiveAnt *models.Ant
		for _, ant := range pathAnts {
			if ant.Position == -1 {
				inactiveAnt = ant
				break
			}
		}

		if inactiveAnt != nil {
			inactiveAnt.Position = 0
			if isFirstRoomEnd {
				inactiveAnt.Arrived = true
			} else {
				roomOccupied[firstRoom] = true
			}

			return &models.TurnMove{AntID: inactiveAnt.ID, Room: firstRoom}
		}
	}
	return nil
}
