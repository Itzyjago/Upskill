package algo

import "testing"

func TestDequePushPopBothEnds(t *testing.T) {
	var d Deque[int]
	d.PushBack(2)
	d.PushBack(3)
	d.PushFront(1)
	d.PushFront(0)
	// order front-to-back: 0, 1, 2, 3
	if got := d.Len(); got != 4 {
		t.Fatalf("Len() = %d, want 4", got)
	}
	if v, ok := d.PopFront(); !ok || v != 0 {
		t.Fatalf("PopFront() = %d, %v, want 0, true", v, ok)
	}
	if v, ok := d.PopBack(); !ok || v != 3 {
		t.Fatalf("PopBack() = %d, %v, want 3, true", v, ok)
	}
	if v, ok := d.PopFront(); !ok || v != 1 {
		t.Fatalf("PopFront() = %d, %v, want 1, true", v, ok)
	}
	if v, ok := d.PopBack(); !ok || v != 2 {
		t.Fatalf("PopBack() = %d, %v, want 2, true", v, ok)
	}
	if d.Len() != 0 {
		t.Fatalf("Len() = %d, want 0 after draining", d.Len())
	}
}

func TestDequeEmptyPops(t *testing.T) {
	var d Deque[string]
	if _, ok := d.PopFront(); ok {
		t.Error("PopFront on empty deque should return ok=false")
	}
	if _, ok := d.PopBack(); ok {
		t.Error("PopBack on empty deque should return ok=false")
	}
}

func TestDequeSingleElementBothPopsSeeIt(t *testing.T) {
	var d Deque[int]
	d.PushFront(9)
	v, ok := d.PopBack()
	if !ok || v != 9 {
		t.Fatalf("PopBack() on single-element deque = %d, %v, want 9, true", v, ok)
	}
	if d.Len() != 0 {
		t.Fatalf("Len() = %d, want 0", d.Len())
	}
}
