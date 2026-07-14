package algo

import "testing"

func TestDijkstraShortestPaths(t *testing.T) {
	g := NewWeightedGraph()
	// 1 -(4)-> 2, 1 -(1)-> 3, 3 -(2)-> 2, 2 -(5)-> 4, 3 -(8)-> 4
	// shortest 1->2 is via 3 (1+2=3), not direct (4).
	g.AddEdge(1, 2, 4)
	g.AddEdge(1, 3, 1)
	g.AddEdge(3, 2, 2)
	g.AddEdge(2, 4, 5)
	g.AddEdge(3, 4, 8)

	dist := g.ShortestPaths(1)
	want := map[int]int{1: 0, 2: 3, 3: 1, 4: 8}
	for node, wantDist := range want {
		if got, ok := dist[node]; !ok || got != wantDist {
			t.Errorf("dist[%d] = %d, %v, want %d", node, got, ok, wantDist)
		}
	}
}

func TestDijkstraUnreachableNodeOmitted(t *testing.T) {
	g := NewWeightedGraph()
	g.AddEdge(1, 2, 1)
	g.AddEdge(3, 4, 1) // disconnected component

	dist := g.ShortestPaths(1)
	if _, ok := dist[3]; ok {
		t.Error("node 3 is unreachable from 1, should not appear in dist")
	}
	if _, ok := dist[4]; ok {
		t.Error("node 4 is unreachable from 1, should not appear in dist")
	}
}
