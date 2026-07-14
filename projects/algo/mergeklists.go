package algo

import "container/heap"

type listHeapItem struct {
	node *ListNode
}

type listMinHeap []listHeapItem

func (h listMinHeap) Len() int            { return len(h) }
func (h listMinHeap) Less(i, j int) bool  { return h[i].node.Val < h[j].node.Val }
func (h listMinHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *listMinHeap) Push(x interface{}) { *h = append(*h, x.(listHeapItem)) }
func (h *listMinHeap) Pop() interface{} {
	old := *h
	n := len(old)
	item := old[n-1]
	*h = old[:n-1]
	return item
}

// MergeKLists merges k already-sorted linked lists into one sorted list.
// O(N log k) time (N = total nodes) using a min-heap over the k lists'
// current heads, vs. merging pairwise which would cost O(N*k) if done
// naively one list at a time.
func MergeKLists(lists []*ListNode) *ListNode {
	h := &listMinHeap{}
	for _, l := range lists {
		if l != nil {
			heap.Push(h, listHeapItem{node: l})
		}
	}
	dummy := &ListNode{}
	tail := dummy
	for h.Len() > 0 {
		item := heap.Pop(h).(listHeapItem)
		tail.Next = item.node
		tail = tail.Next
		if item.node.Next != nil {
			heap.Push(h, listHeapItem{node: item.node.Next})
		}
	}
	return dummy.Next
}
