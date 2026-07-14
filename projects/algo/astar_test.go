package algo

import "testing"

func gridFromStrings(rows []string) [][]bool {
	grid := make([][]bool, len(rows))
	for i, row := range rows {
		grid[i] = make([]bool, len(row))
		for j, ch := range row {
			grid[i][j] = ch == '#'
		}
	}
	return grid
}

func TestAStarGridPathFindsShortestPath(t *testing.T) {
	grid := gridFromStrings([]string{
		".....",
		".###.",
		".....",
	})
	path := AStarGridPath(grid, GridPoint{0, 0}, GridPoint{2, 4})
	if path == nil {
		t.Fatal("expected a path to exist")
	}
	// Manhattan distance is 6, and this grid allows a path of exactly that
	// length (go around the wall along a shortest route) — path length in
	// nodes is steps+1.
	if got := len(path) - 1; got != 6 {
		t.Errorf("path length (steps) = %d, want 6 (shortest possible)", got)
	}
	if path[0] != (GridPoint{0, 0}) || path[len(path)-1] != (GridPoint{2, 4}) {
		t.Errorf("path = %v, should start at (0,0) and end at (2,4)", path)
	}
	for i := 1; i < len(path); i++ {
		if manhattan(path[i-1], path[i]) != 1 {
			t.Errorf("path %v is not made of single grid steps between %v and %v", path, path[i-1], path[i])
		}
	}
}

func TestAStarGridPathNoPathExists(t *testing.T) {
	grid := gridFromStrings([]string{
		".#.",
		".#.",
		".#.",
	})
	if path := AStarGridPath(grid, GridPoint{0, 0}, GridPoint{0, 2}); path != nil {
		t.Errorf("AStarGridPath = %v, want nil (wall fully separates start and goal)", path)
	}
}

func TestAStarGridPathStartEqualsGoal(t *testing.T) {
	grid := gridFromStrings([]string{"."})
	path := AStarGridPath(grid, GridPoint{0, 0}, GridPoint{0, 0})
	if len(path) != 1 {
		t.Errorf("AStarGridPath(start==goal) = %v, want a single-point path", path)
	}
}

func TestAStarGridPathBlockedStart(t *testing.T) {
	grid := gridFromStrings([]string{"#."})
	if path := AStarGridPath(grid, GridPoint{0, 0}, GridPoint{0, 1}); path != nil {
		t.Errorf("AStarGridPath with blocked start = %v, want nil", path)
	}
}
