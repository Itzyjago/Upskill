package algo

import (
	"reflect"
	"testing"
)

// 1 -> 2, 1 -> 3, 2 -> 4, 3 -> 4, 4 -> 5 (a diamond into a tail)
func buildDiamondGraph() *Graph {
	g := NewGraph()
	g.AddEdge(1, 2)
	g.AddEdge(1, 3)
	g.AddEdge(2, 4)
	g.AddEdge(3, 4)
	g.AddEdge(4, 5)
	return g
}

func TestGraphBFS(t *testing.T) {
	g := buildDiamondGraph()
	want := []int{1, 2, 3, 4, 5}
	if got := g.BFS(1); !reflect.DeepEqual(got, want) {
		t.Errorf("BFS(1) = %v, want %v", got, want)
	}
}

func TestGraphDFS(t *testing.T) {
	g := buildDiamondGraph()
	want := []int{1, 2, 4, 5, 3}
	if got := g.DFS(1); !reflect.DeepEqual(got, want) {
		t.Errorf("DFS(1) = %v, want %v", got, want)
	}
}

func TestGraphBFSVisitsEachNodeOnce(t *testing.T) {
	g := buildDiamondGraph()
	order := g.BFS(1)
	seen := map[int]bool{}
	for _, n := range order {
		if seen[n] {
			t.Fatalf("node %d visited twice in %v", n, order)
		}
		seen[n] = true
	}
}
