# Phase 5 — Ant Scheduling

| Task ID | Description |
|---------|-------------|
| **T5.1** | Implement `distributor.DistributeAnts(paths [][]string, numAnts int) []*Ant` — assigns each ant to a path |
| **T5.2** | Implement `simulator.Simulate(ants []*Ant, graph *Graph) []Turn` — turn-by-turn simulation |
| **T5.3** | Handle the constraint: one ant per room, one use per tunnel per turn |

**Simulation Algorithm (detailed):**

```
// Data structures
antQueues[pathIndex] = queue of ants assigned to that path
antPosition[antID] = (pathIndex, stepIndex)  // step 0 = first room after start
roomOccupied[roomName] = false  // reset each turn
tunnelUsed[(room1, room2)] = false  // reset each turn

turns = []

while not all ants have arrived:
    currentTurn = []
    
    // Process each path
    for each path p:
        for each ant a on path p (from front of path to back):
            currentRoom = antPosition[a].room
            nextRoom = path[p].rooms[antPosition[a].stepIndex + 1]
            
            // Check if ant can move
            canMove = false
            if nextRoom is end:
                canMove = true  // end can hold unlimited
            else if not roomOccupied[nextRoom]:
                canMove = true
            
            if canMove and not tunnelUsed[(currentRoom, nextRoom)]:
                // Move ant
                antPosition[a].stepIndex++
                roomOccupied[currentRoom] = false
                if nextRoom != end:
                    roomOccupied[nextRoom] = true
                tunnelUsed[(currentRoom, nextRoom)] = true
                currentTurn.append((a.id, nextRoom))
                
                if nextRoom == end:
                    mark ant as arrived
    
    // Inject new ants at start of each path if room is free
    for each path p:
        firstRoom = path[p].rooms[1]  // first room after start
        if antQueues[p] not empty and not roomOccupied[firstRoom]:
            ant = antQueues[p].dequeue()
            if firstRoom == end:
                mark ant arrived
                currentTurn.append((ant.id, end))
            else:
                roomOccupied[firstRoom] = true
                antPosition[ant.id] = (p, 1)
                currentTurn.append((ant.id, firstRoom))
    
    turns.append(currentTurn)

return turns
```
