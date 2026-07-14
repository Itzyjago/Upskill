package algo

import "testing"

func TestUnionFindConnected(t *testing.T) {
	uf := NewUnionFind(6)
	uf.Union(0, 1)
	uf.Union(1, 2)
	uf.Union(3, 4)

	if !uf.Connected(0, 2) {
		t.Error("0 and 2 should be connected via 1")
	}
	if uf.Connected(0, 3) {
		t.Error("0 and 3 should not be connected")
	}
	if uf.Connected(4, 5) {
		t.Error("4 and 5 should not be connected")
	}
	uf.Union(2, 3)
	if !uf.Connected(0, 4) {
		t.Error("0 and 4 should be connected after merging the two groups")
	}
}

func TestUnionFindUnionReturnsFalseWhenAlreadyConnected(t *testing.T) {
	uf := NewUnionFind(3)
	if !uf.Union(0, 1) {
		t.Error("first Union(0, 1) should return true")
	}
	if uf.Union(0, 1) {
		t.Error("second Union(0, 1) should return false, already connected")
	}
}

func TestUnionFindCycleDetection(t *testing.T) {
	// Classic use: adding edges to a graph, Union returns false exactly
	// when an edge would close a cycle.
	uf := NewUnionFind(4)
	edges := [][2]int{{0, 1}, {1, 2}, {2, 3}, {3, 0}}
	var cycleEdges int
	for _, e := range edges {
		if !uf.Union(e[0], e[1]) {
			cycleEdges++
		}
	}
	if cycleEdges != 1 {
		t.Errorf("cycleEdges = %d, want 1 (the closing edge 3-0)", cycleEdges)
	}
}
