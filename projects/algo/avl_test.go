package algo

import (
	"reflect"
	"testing"
)

func buildAVL(vals []int) *AVLNode {
	var root *AVLNode
	for _, v := range vals {
		root = root.Insert(v)
	}
	return root
}

func TestAVLInOrderIsSorted(t *testing.T) {
	root := buildAVL([]int{5, 3, 8, 1, 4, 7, 9, 2, 6})
	want := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	if got := root.InOrder(); !reflect.DeepEqual(got, want) {
		t.Errorf("InOrder() = %v, want %v", got, want)
	}
}

func TestAVLStaysBalancedOnSortedInserts(t *testing.T) {
	// A plain BST degrades to a linked list (height n) on sorted-order
	// inserts. AVL must keep height within the O(log n) bound: for n
	// nodes, height <= ~1.44*log2(n+2) (a known AVL bound); check a much
	// looser but still meaningful bound of 2*log2(n+1) to avoid coupling
	// the test to the exact constant.
	var root *AVLNode
	n := 100
	for i := 0; i < n; i++ {
		root = root.Insert(i)
	}
	h := height(root)
	maxAcceptable := 0
	for x := n + 1; x > 1; x >>= 1 {
		maxAcceptable++
	}
	maxAcceptable *= 2
	if h > maxAcceptable {
		t.Errorf("height = %d after %d sorted inserts, want <= %d (O(log n), not degraded to a list)", h, n, maxAcceptable)
	}
}

func TestAVLDuplicateInsertIgnored(t *testing.T) {
	root := buildAVL([]int{5, 3, 5, 3})
	want := []int{3, 5}
	if got := root.InOrder(); !reflect.DeepEqual(got, want) {
		t.Errorf("InOrder() = %v, want %v", got, want)
	}
}
