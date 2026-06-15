package scheduler

import (
	"lem-in/pkg/models"
	"sort"
)

/*Simulate distributes ants across paths and simulates the movement of ants.

	Args:
		ants: A slice of pointers to Ant structs representing the ants to be distributed.
		graph: A pointer to a Graph struct representing the graph.

	Returns:

	A slice of Turn structs representing the turns taken by the ants.
*/

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

			firstRoom := path[0]
			isFirstRoomEnd := len(path) == 1
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

					currentTurn = append(currentTurn, models.TurnMove{
						AntID: inactiveAnt.ID,
						Room:  firstRoom,
					})
				}
			}
		}

		if len(currentTurn) > 0 {
			turns = append(turns, currentTurn)
		} else {
			break
		}
	}

	return turns
}

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

func allAntsArrived(ants []*models.Ant) bool {
	for _, ant := range ants {
		if !ant.Arrived {
			return false
		}
	}
	return true
}

func getOccupiedRooms(ants []*models.Ant) map[string]bool {
	roomOccupied := make(map[string]bool)
	for _, ant := range ants {
		if !ant.Arrived && ant.Position >= 0 {
			currRoom := ant.Path[ant.Position]
			if currRoom != ant.Path[len(ant.Path)-1] {
				roomOccupied[currRoom] = true
			}
		}
	}
	return roomOccupied
}

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

func getSortedAntsOnPath(ants []*models.Ant, path []string) []*models.Ant {
	var pathAnts []*models.Ant
	for _, ant := range ants {
		if equalPaths(ant.Path, path) {
			pathAnts = append(pathAnts, ant)
		}
	}
	sort.Slice(pathAnts, func(i, j int) bool {
		return pathAnts[i].Position > pathAnts[j].Position
	})
	return pathAnts
}

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

			moves = append(moves, models.TurnMove{
				AntID: ant.ID,
				Room:  nextRoom,
			})
		}
	}
	return moves
}
