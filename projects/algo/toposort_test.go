package algo

import "testing"

func TestTopoSortValidOrdering(t *testing.T) {
	g := NewGraph()
	g.AddEdge(5, 2)
	g.AddEdge(5, 0)
	g.AddEdge(4, 0)
	g.AddEdge(4, 1)
	g.AddEdge(2, 3)
	g.AddEdge(3, 1)

	order, ok := g.TopoSort()
	if !ok {
		t.Fatal("TopoSort() ok = false, want true (this DAG has no cycle)")
	}
	pos := make(map[int]int, len(order))
	for i, n := range order {
		pos[n] = i
	}
	edges := [][2]int{{5, 2}, {5, 0}, {4, 0}, {4, 1}, {2, 3}, {3, 1}}
	for _, e := range edges {
		if pos[e[0]] >= pos[e[1]] {
			t.Errorf("edge %d->%d violated: pos[%d]=%d, pos[%d]=%d", e[0], e[1], e[0], pos[e[0]], e[1], pos[e[1]])
		}
	}
}

func TestTopoSortDetectsCycle(t *testing.T) {
	g := NewGraph()
	g.AddEdge(1, 2)
	g.AddEdge(2, 3)
	g.AddEdge(3, 1) // closes the cycle

	if _, ok := g.TopoSort(); ok {
		t.Error("TopoSort() ok = true, want false (graph has a cycle)")
	}
}
