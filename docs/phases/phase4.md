# Phase 4 — Pathfinding

| Task ID | Description |
|---------|-------------|
| **T4.1** | Implement `bfs.FindShortestPath(g *FlowGraph, start, end string) []string` — unweighted BFS on residual graph |
| **T4.2** | Implement `edmonds_karp.FindAllDisjointPaths(g *FlowGraph) [][]string` — repeatedly find augmenting paths, collect all |
| **T4.3** | Convert flow paths back to original room names (collapse split nodes) |
| **T4.4** | Implement `selector.SelectOptimalPaths(allPaths [][]string, antCount int) [][]string` |
| **T4.5** | For each subset size k from 1 to len(allPaths), compute $T(k)$ and pick minimum |

**Optimal Path Selection Algorithm:**

```
Input: paths (sorted by length ascending), numAnts
Output: bestPaths

bestTurns = ∞
bestPaths = []

for k = 1 to len(paths):
    selected = paths[0:k]  // shortest k paths
    turns = computeTurns(selected, numAnts)
    if turns < bestTurns:
        bestTurns = turns
        bestPaths = selected

return bestPaths
```

**`computeTurns` algorithm:**

Given $k$ paths of lengths $L_1 \leq L_2 \leq \ldots \leq L_k$ (edge counts) and $N$ ants:

We want to find the minimal $T$ such that ants can be distributed across paths and all finish by turn $T$.

For path $i$, the maximum ants it can finish by turn $T$ is: $\max(0, T - L_i + 1)$.

Binary search for $T$:

```
low = 0, high = numAnts + max(L_i)

while low < high:
    mid = (low + high) / 2
    capacity = sum(max(0, mid - L_i + 1) for each path i)
    if capacity >= numAnts:
        high = mid
    else:
        low = mid + 1

return low
```

Then assign ants: path $i$ gets $\max(0, T_{opt} - L_i + 1)$ ants, but capped so total = $N$.
