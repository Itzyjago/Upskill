package algo

import "container/heap"

// GridPoint is a cell in a 2D grid.
type GridPoint struct{ Row, Col int }

type astarItem struct {
	point GridPoint
	fCost int
}

type astarHeap []astarItem

func (h astarHeap) Len() int            { return len(h) }
func (h astarHeap) Less(i, j int) bool  { return h[i].fCost < h[j].fCost }
func (h astarHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *astarHeap) Push(x interface{}) { *h = append(*h, x.(astarItem)) }
func (h *astarHeap) Pop() interface{} {
	old := *h
	n := len(old)
	item := old[n-1]
	*h = old[:n-1]
	return item
}

func manhattan(a, b GridPoint) int {
	dr, dc := a.Row-b.Row, a.Col-b.Col
	if dr < 0 {
		dr = -dr
	}
	if dc < 0 {
		dc = -dc
	}
	return dr + dc
}

// AStarGridPath finds the shortest path from start to goal on a grid where
// grid[r][c] == true means "blocked". Like Dijkstra but guided by the
// Manhattan-distance heuristic (admissible since moves cost 1 and the grid
// only allows up/down/left/right — the heuristic never overestimates), so
// it explores fewer nodes than a plain BFS/Dijkstra would on a large open
// grid. Returns nil if no path exists.
func AStarGridPath(grid [][]bool, start, goal GridPoint) []GridPoint {
	rows, cols := len(grid), len(grid[0])
	inBounds := func(p GridPoint) bool {
		return p.Row >= 0 && p.Row < rows && p.Col >= 0 && p.Col < cols
	}
	if !inBounds(start) || !inBounds(goal) || grid[start.Row][start.Col] || grid[goal.Row][goal.Col] {
		return nil
	}

	gScore := map[GridPoint]int{start: 0}
	cameFrom := map[GridPoint]GridPoint{}
	h := &astarHeap{{point: start, fCost: manhattan(start, goal)}}
	heap.Init(h)
	visited := map[GridPoint]bool{}

	dirs := []GridPoint{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}
	for h.Len() > 0 {
		cur := heap.Pop(h).(astarItem).point
		if cur == goal {
			return reconstructPath(cameFrom, cur)
		}
		if visited[cur] {
			continue
		}
		visited[cur] = true

		for _, d := range dirs {
			next := GridPoint{Row: cur.Row + d.Row, Col: cur.Col + d.Col}
			if !inBounds(next) || grid[next.Row][next.Col] || visited[next] {
				continue
			}
			tentative := gScore[cur] + 1
			if g, ok := gScore[next]; !ok || tentative < g {
				gScore[next] = tentative
				cameFrom[next] = cur
				heap.Push(h, astarItem{point: next, fCost: tentative + manhattan(next, goal)})
			}
		}
	}
	return nil
}

func reconstructPath(cameFrom map[GridPoint]GridPoint, end GridPoint) []GridPoint {
	path := []GridPoint{end}
	for {
		prev, ok := cameFrom[path[len(path)-1]]
		if !ok {
			break
		}
		path = append(path, prev)
	}
	for l, r := 0, len(path)-1; l < r; l, r = l+1, r-1 {
		path[l], path[r] = path[r], path[l]
	}
	return path
}
