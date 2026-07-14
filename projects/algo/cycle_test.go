package algo

import "testing"

func TestHasCycleNoCycle(t *testing.T) {
	head := NewList([]int{1, 2, 3, 4})
	if HasCycle(head) {
		t.Error("HasCycle on an acyclic list should be false")
	}
}

func TestHasCycleWithCycle(t *testing.T) {
	head := NewList([]int{1, 2, 3, 4})
	// find node with val 3 and point it back to node with val 2,
	// forming a cycle: 1 -> 2 -> 3 -> 2 -> 3 -> ...
	var n2, n3 *ListNode
	for n := head; n != nil; n = n.Next {
		if n.Val == 2 {
			n2 = n
		}
		if n.Val == 3 {
			n3 = n
		}
	}
	n3.Next = n2
	if !HasCycle(head) {
		t.Error("HasCycle should detect the cycle")
	}
}

func TestMiddleNode(t *testing.T) {
	if got := MiddleNode(NewList([]int{1, 2, 3, 4, 5})).Val; got != 3 {
		t.Errorf("MiddleNode odd length = %d, want 3", got)
	}
	if got := MiddleNode(NewList([]int{1, 2, 3, 4})).Val; got != 3 {
		t.Errorf("MiddleNode even length = %d, want 3 (second of two middles)", got)
	}
	if got := MiddleNode(NewList([]int{1})).Val; got != 1 {
		t.Errorf("MiddleNode single element = %d, want 1", got)
	}
}
