package algo

import "testing"

func TestKruskalMST(t *testing.T) {
	// Classic textbook graph, 5 nodes; known-minimum total weight is 16.
	edges := []Edge{
		{0, 1, 2}, {0, 3, 6}, {1, 2, 3}, {1, 3, 8}, {1, 4, 5},
		{2, 4, 7}, {3, 4, 9},
	}
	mst := KruskalMST(5, edges)

	if got := len(mst); got != 4 {
		t.Fatalf("MST for 5 nodes should have 4 edges, got %d", got)
	}
	total := 0
	for _, e := range mst {
		total += e.Weight
	}
	if total != 16 {
		t.Errorf("MST total weight = %d, want 16", total)
	}

	uf := NewUnionFind(5)
	for _, e := range mst {
		if !uf.Union(e.From, e.To) {
			t.Errorf("MST edge %+v closes a cycle, not a valid tree", e)
		}
	}
}

func TestKruskalMSTDisconnectedGraph(t *testing.T) {
	// Two disjoint components: {0,1} and {2,3}. No spanning tree exists
	// over all 4 nodes, so the result covers each component separately.
	edges := []Edge{{0, 1, 1}, {2, 3, 1}}
	mst := KruskalMST(4, edges)
	if got := len(mst); got != 2 {
		t.Fatalf("expected 2 edges (one per component), got %d", got)
	}
}
