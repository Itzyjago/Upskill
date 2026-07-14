package algo

import (
	"reflect"
	"testing"
)

func TestMergeKLists(t *testing.T) {
	lists := []*ListNode{
		NewList([]int{1, 4, 5}),
		NewList([]int{1, 3, 4}),
		NewList([]int{2, 6}),
	}
	want := []int{1, 1, 2, 3, 4, 4, 5, 6}
	got := MergeKLists(lists).ToSlice()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("MergeKLists(...) = %v, want %v", got, want)
	}
}

func TestMergeKListsWithEmptyLists(t *testing.T) {
	lists := []*ListNode{nil, NewList([]int{1}), nil}
	got := MergeKLists(lists).ToSlice()
	want := []int{1}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("MergeKLists with nil lists = %v, want %v", got, want)
	}
}

func TestMergeKListsAllEmpty(t *testing.T) {
	if got := MergeKLists([]*ListNode{nil, nil}); got != nil {
		t.Errorf("MergeKLists(all nil) = %v, want nil", got)
	}
	if got := MergeKLists(nil); got != nil {
		t.Errorf("MergeKLists(nil) = %v, want nil", got)
	}
}
