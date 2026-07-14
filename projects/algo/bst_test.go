package algo

import (
	"reflect"
	"testing"
)

func buildBST(vals []int) *BSTNode {
	var root *BSTNode
	for _, v := range vals {
		root = root.Insert(v)
	}
	return root
}

func TestBSTInOrderIsSorted(t *testing.T) {
	root := buildBST([]int{5, 3, 8, 1, 4, 7, 9})
	want := []int{1, 3, 4, 5, 7, 8, 9}
	if got := root.InOrder(); !reflect.DeepEqual(got, want) {
		t.Errorf("InOrder() = %v, want %v", got, want)
	}
}

func TestBSTContains(t *testing.T) {
	root := buildBST([]int{5, 3, 8, 1, 4, 7, 9})
	for _, v := range []int{5, 1, 9, 4} {
		if !root.Contains(v) {
			t.Errorf("Contains(%d) = false, want true", v)
		}
	}
	for _, v := range []int{0, 6, 100} {
		if root.Contains(v) {
			t.Errorf("Contains(%d) = true, want false", v)
		}
	}
}

func TestBSTDuplicateInsertIgnored(t *testing.T) {
	root := buildBST([]int{5, 3, 5, 3})
	want := []int{3, 5}
	if got := root.InOrder(); !reflect.DeepEqual(got, want) {
		t.Errorf("InOrder() = %v, want %v (duplicates should not add nodes)", got, want)
	}
}
