package algo

// Deque is a generic doubly linked list supporting O(1) push/pop at both
// ends — the operation LRUCache's node type also needed, pulled out here
// as its own reusable structure instead of copy-pasted.
type Deque[T any] struct {
	head, tail *dequeNode[T]
	length     int
}

type dequeNode[T any] struct {
	val        T
	prev, next *dequeNode[T]
}

func (d *Deque[T]) Len() int { return d.length }

func (d *Deque[T]) PushFront(v T) {
	n := &dequeNode[T]{val: v, next: d.head}
	if d.head != nil {
		d.head.prev = n
	}
	d.head = n
	if d.tail == nil {
		d.tail = n
	}
	d.length++
}

func (d *Deque[T]) PushBack(v T) {
	n := &dequeNode[T]{val: v, prev: d.tail}
	if d.tail != nil {
		d.tail.next = n
	}
	d.tail = n
	if d.head == nil {
		d.head = n
	}
	d.length++
}

func (d *Deque[T]) PopFront() (v T, ok bool) {
	if d.head == nil {
		return v, false
	}
	n := d.head
	d.head = n.next
	if d.head != nil {
		d.head.prev = nil
	} else {
		d.tail = nil
	}
	d.length--
	return n.val, true
}

func (d *Deque[T]) PopBack() (v T, ok bool) {
	if d.tail == nil {
		return v, false
	}
	n := d.tail
	d.tail = n.prev
	if d.tail != nil {
		d.tail.next = nil
	} else {
		d.head = nil
	}
	d.length--
	return n.val, true
}
