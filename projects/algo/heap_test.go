package algo

import "testing"

func TestMinHeapPopsInSortedOrder(t *testing.T) {
	var h MinHeap
	in := []int{5, 3, 8, 1, 9, 2, 7}
	for _, v := range in {
		h.Push(v)
	}
	want := []int{1, 2, 3, 5, 7, 8, 9}
	for _, w := range want {
		v, ok := h.Pop()
		if !ok || v != w {
			t.Fatalf("Pop() = %d, %v, want %d, true", v, ok, w)
		}
	}
	if _, ok := h.Pop(); ok {
		t.Fatal("Pop on empty heap should return ok=false")
	}
}

func TestMinHeapLen(t *testing.T) {
	var h MinHeap
	h.Push(1)
	h.Push(2)
	if got := h.Len(); got != 2 {
		t.Fatalf("Len() = %d, want 2", got)
	}
	h.Pop()
	if got := h.Len(); got != 1 {
		t.Fatalf("Len() after Pop = %d, want 1", got)
	}
}
