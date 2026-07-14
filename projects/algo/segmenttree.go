package algo

// SegmentTree supports range-sum queries and point updates over a fixed
// array, both O(log n) — an alternative to FenwickTree that generalizes
// more easily to range min/max (swapping the merge function), at the cost
// of ~4x the array's memory for the tree itself.
type SegmentTree struct {
	n    int
	tree []int
}

func NewSegmentTree(nums []int) *SegmentTree {
	n := len(nums)
	st := &SegmentTree{n: n, tree: make([]int, 4*n)}
	if n > 0 {
		st.build(nums, 0, 0, n-1)
	}
	return st
}

func (st *SegmentTree) build(nums []int, node, lo, hi int) {
	if lo == hi {
		st.tree[node] = nums[lo]
		return
	}
	mid := lo + (hi-lo)/2
	left, right := 2*node+1, 2*node+2
	st.build(nums, left, lo, mid)
	st.build(nums, right, mid+1, hi)
	st.tree[node] = st.tree[left] + st.tree[right]
}

// Update sets the value at index i to v.
func (st *SegmentTree) Update(i, v int) {
	st.update(0, 0, st.n-1, i, v)
}

func (st *SegmentTree) update(node, lo, hi, i, v int) {
	if lo == hi {
		st.tree[node] = v
		return
	}
	mid := lo + (hi-lo)/2
	left, right := 2*node+1, 2*node+2
	if i <= mid {
		st.update(left, lo, mid, i, v)
	} else {
		st.update(right, mid+1, hi, i, v)
	}
	st.tree[node] = st.tree[left] + st.tree[right]
}

// Query returns the sum of elements in [qlo, qhi] (inclusive).
func (st *SegmentTree) Query(qlo, qhi int) int {
	return st.query(0, 0, st.n-1, qlo, qhi)
}

func (st *SegmentTree) query(node, lo, hi, qlo, qhi int) int {
	if qhi < lo || hi < qlo {
		return 0 // no overlap
	}
	if qlo <= lo && hi <= qhi {
		return st.tree[node] // fully covered
	}
	mid := lo + (hi-lo)/2
	return st.query(2*node+1, lo, mid, qlo, qhi) + st.query(2*node+2, mid+1, hi, qlo, qhi)
}
