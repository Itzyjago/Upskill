package algo

// UnionFind (disjoint-set) tracks a partition over 0..n-1, with near-O(1)
// amortized Find/Union via path compression + union by rank.
type UnionFind struct {
	parent []int
	rank   []int
}

func NewUnionFind(n int) *UnionFind {
	uf := &UnionFind{parent: make([]int, n), rank: make([]int, n)}
	for i := range uf.parent {
		uf.parent[i] = i
	}
	return uf
}

// Find returns the representative of x's set, compressing the path as it goes.
func (uf *UnionFind) Find(x int) int {
	if uf.parent[x] != x {
		uf.parent[x] = uf.Find(uf.parent[x])
	}
	return uf.parent[x]
}

// Union merges the sets containing x and y. Returns false if they were
// already in the same set (useful for cycle detection).
func (uf *UnionFind) Union(x, y int) bool {
	rx, ry := uf.Find(x), uf.Find(y)
	if rx == ry {
		return false
	}
	switch {
	case uf.rank[rx] < uf.rank[ry]:
		rx, ry = ry, rx
	case uf.rank[rx] == uf.rank[ry]:
		uf.rank[rx]++
	}
	uf.parent[ry] = rx
	return true
}

func (uf *UnionFind) Connected(x, y int) bool {
	return uf.Find(x) == uf.Find(y)
}
