package algo

import "sort"

// Edge is an undirected weighted edge between two nodes, used for MST
// algorithms (direction doesn't matter for a spanning tree).
type Edge struct {
	From, To, Weight int
}

// KruskalMST returns the edges of a minimum spanning tree over the given
// edges and node count n (nodes 0..n-1), using union-find to reject edges
// that would close a cycle. O(E log E) time, dominated by the sort.
func KruskalMST(n int, edges []Edge) []Edge {
	sorted := make([]Edge, len(edges))
	copy(sorted, edges)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].Weight < sorted[j].Weight })

	uf := NewUnionFind(n)
	var mst []Edge
	for _, e := range sorted {
		if uf.Union(e.From, e.To) {
			mst = append(mst, e)
		}
	}
	return mst
}
