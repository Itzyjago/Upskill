package algo

// FenwickTree (Binary Indexed Tree) supports point updates and prefix-sum
// queries over a fixed-size array, both in O(log n) — the trick being each
// index i stores a partial sum covering a range determined by i's lowest
// set bit, so both operations only ever touch O(log n) indices.
type FenwickTree struct {
	tree []int // 1-indexed internally
}

func NewFenwickTree(n int) *FenwickTree {
	return &FenwickTree{tree: make([]int, n+1)}
}

// Add adds delta to the value at index i (0-indexed).
func (f *FenwickTree) Add(i, delta int) {
	for i++; i < len(f.tree); i += i & (-i) {
		f.tree[i] += delta
	}
}

// PrefixSum returns the sum of elements in [0, i] (0-indexed, inclusive).
func (f *FenwickTree) PrefixSum(i int) int {
	sum := 0
	for i++; i > 0; i -= i & (-i) {
		sum += f.tree[i]
	}
	return sum
}

// RangeSum returns the sum of elements in [lo, hi] (0-indexed, inclusive).
func (f *FenwickTree) RangeSum(lo, hi int) int {
	if lo == 0 {
		return f.PrefixSum(hi)
	}
	return f.PrefixSum(hi) - f.PrefixSum(lo-1)
}
