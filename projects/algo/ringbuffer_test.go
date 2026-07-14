package algo

import "testing"

func TestRingBufferFIFOOrder(t *testing.T) {
	rb := NewRingBuffer(3)
	for _, v := range []int{1, 2, 3} {
		if !rb.Enqueue(v) {
			t.Fatalf("Enqueue(%d) = false, want true (buffer not yet full)", v)
		}
	}
	if rb.Enqueue(4) {
		t.Fatal("Enqueue on full buffer should return false")
	}
	for _, want := range []int{1, 2, 3} {
		v, ok := rb.Dequeue()
		if !ok || v != want {
			t.Fatalf("Dequeue() = %d, %v, want %d, true", v, ok, want)
		}
	}
	if _, ok := rb.Dequeue(); ok {
		t.Fatal("Dequeue on empty buffer should return ok=false")
	}
}

func TestRingBufferWrapsAroundCorrectly(t *testing.T) {
	rb := NewRingBuffer(3)
	rb.Enqueue(1)
	rb.Enqueue(2)
	rb.Dequeue() // drop 1, head moves forward
	rb.Enqueue(3)
	rb.Enqueue(4) // wraps the tail index around to slot 0
	if !rb.IsFull() {
		t.Fatal("buffer should be full after 3 live elements at capacity 3")
	}
	for _, want := range []int{2, 3, 4} {
		v, ok := rb.Dequeue()
		if !ok || v != want {
			t.Fatalf("Dequeue() = %d, %v, want %d, true", v, ok, want)
		}
	}
}
