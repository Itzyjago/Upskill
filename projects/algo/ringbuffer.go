package algo

// RingBuffer is a fixed-capacity FIFO backed by a circular array — no
// reslicing or shifting on Dequeue, indices just wrap via modulo. Useful
// where Queue's occasional backing-array reallocation isn't acceptable
// (e.g. a fixed memory budget).
type RingBuffer struct {
	data        []int
	head, count int // head = index of oldest element
}

func NewRingBuffer(capacity int) *RingBuffer {
	return &RingBuffer{data: make([]int, capacity)}
}

func (r *RingBuffer) Len() int      { return r.count }
func (r *RingBuffer) Cap() int      { return len(r.data) }
func (r *RingBuffer) IsFull() bool  { return r.count == len(r.data) }
func (r *RingBuffer) IsEmpty() bool { return r.count == 0 }

// Enqueue adds v. Returns false if the buffer is already full (capacity is
// fixed, not grown).
func (r *RingBuffer) Enqueue(v int) bool {
	if r.IsFull() {
		return false
	}
	tail := (r.head + r.count) % len(r.data)
	r.data[tail] = v
	r.count++
	return true
}

func (r *RingBuffer) Dequeue() (v int, ok bool) {
	if r.IsEmpty() {
		return 0, false
	}
	v = r.data[r.head]
	r.head = (r.head + 1) % len(r.data)
	r.count--
	return v, true
}
