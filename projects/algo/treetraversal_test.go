package algo

import (
	"reflect"
	"testing"
)

func TestBSTLevelOrder(t *testing.T) {
	root := buildBST([]int{5, 3, 8, 1, 4, 7, 9})
	want := [][]int{
		{5},
		{3, 8},
		{1, 4, 7, 9},
	}
	if got := root.LevelOrder(); !reflect.DeepEqual(got, want) {
		t.Errorf("LevelOrder() = %v, want %v", got, want)
	}
}

func TestBSTLevelOrderNil(t *testing.T) {
	var root *BSTNode
	if got := root.LevelOrder(); got != nil {
		t.Errorf("LevelOrder() on nil tree = %v, want nil", got)
	}
}

func TestBSTHeight(t *testing.T) {
	var root *BSTNode
	if got := root.Height(); got != 0 {
		t.Errorf("Height() of nil tree = %d, want 0", got)
	}
	root = buildBST([]int{5})
	if got := root.Height(); got != 1 {
		t.Errorf("Height() of single node = %d, want 1", got)
	}
	// sorted-order inserts on a plain BST degrade to a linked list: height == count.
	root = buildBST([]int{1, 2, 3, 4, 5})
	if got := root.Height(); got != 5 {
		t.Errorf("Height() after sorted inserts = %d, want 5 (degraded BST)", got)
	}
}
