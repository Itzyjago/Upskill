package algo

import "testing"

func TestBellmanFordNegativeWeights(t *testing.T) {
	// 0 -> 1 (4), 0 -> 2 (5), 1 -> 2 (-3), 2 -> 3 (4)
	// shortest 0->2 is via 1 (4 + -3 = 1), not direct (5).
	edges := []Edge{
		{0, 1, 4}, {0, 2, 5}, {1, 2, -3}, {2, 3, 4},
	}
	dist, negCycle := BellmanFord(4, edges, 0)
	if negCycle {
		t.Fatal("negativeCycle = true, want false")
	}
	want := []int{0, 4, 1, 5}
	for i, w := range want {
		if dist[i] != w {
			t.Errorf("dist[%d] = %d, want %d", i, dist[i], w)
		}
	}
}

func TestBellmanFordDetectsNegativeCycle(t *testing.T) {
	// 0 -> 1 (1), 1 -> 2 (-1), 2 -> 1 (-1): the 1<->2 edges form a cycle
	// that keeps getting cheaper forever.
	edges := []Edge{
		{0, 1, 1}, {1, 2, -1}, {2, 1, -1},
	}
	_, negCycle := BellmanFord(3, edges, 0)
	if !negCycle {
		t.Error("negativeCycle = false, want true")
	}
}

func TestBellmanFordAgreesWithDijkstraWhenNonNegative(t *testing.T) {
	edges := []Edge{
		{1, 2, 4}, {1, 3, 1}, {3, 2, 2}, {2, 4, 5}, {3, 4, 8},
	}
	dist, negCycle := BellmanFord(5, edges, 1)
	if negCycle {
		t.Fatal("negativeCycle = true, want false")
	}

	g := NewWeightedGraph()
	for _, e := range edges {
		g.AddEdge(e.From, e.To, e.Weight)
	}
	dijkstraDist := g.ShortestPaths(1)
	for node, want := range dijkstraDist {
		if got := dist[node]; got != want {
			t.Errorf("BellmanFord dist[%d] = %d, Dijkstra says %d", node, got, want)
		}
	}
}
