package algo

import "testing"

func TestQueue(t *testing.T) {
	var q Queue[int]
	if _, ok := q.Dequeue(); ok {
		t.Fatal("Dequeue on empty queue should return ok=false")
	}
	for i := 1; i <= 5; i++ {
		q.Enqueue(i)
	}
	if got := q.Len(); got != 5 {
		t.Fatalf("Len() = %d, want 5", got)
	}
	for _, want := range []int{1, 2, 3, 4, 5} {
		v, ok := q.Dequeue()
		if !ok || v != want {
			t.Fatalf("Dequeue() = %d, %v, want %d, true", v, ok, want)
		}
	}
	if _, ok := q.Dequeue(); ok {
		t.Fatal("Dequeue after draining should return ok=false")
	}
}

func TestQueueCompaction(t *testing.T) {
	var q Queue[int]
	for i := 0; i < 100; i++ {
		q.Enqueue(i)
		if _, ok := q.Dequeue(); !ok {
			t.Fatalf("Dequeue() should succeed at i=%d", i)
		}
	}
	if got := q.Len(); got != 0 {
		t.Fatalf("Len() = %d, want 0", got)
	}
}
