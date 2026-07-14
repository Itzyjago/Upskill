package algo

// Queue is a generic FIFO built on a slice. Dequeue is O(1) amortized: it
// slices off the front instead of shifting, and only reslices from index 0
// once the backing array's free space gets large relative to its length.
type Queue[T any] struct {
	items []T
	head  int
}

func (q *Queue[T]) Enqueue(v T) {
	q.items = append(q.items, v)
}

func (q *Queue[T]) Dequeue() (v T, ok bool) {
	if q.head >= len(q.items) {
		return v, false
	}
	v = q.items[q.head]
	q.items[q.head] = *new(T) // drop the reference so it can be GC'd
	q.head++
	if q.head > 8 && q.head*2 > len(q.items) {
		q.items = append([]T(nil), q.items[q.head:]...)
		q.head = 0
	}
	return v, true
}

func (q *Queue[T]) Len() int {
	return len(q.items) - q.head
}
