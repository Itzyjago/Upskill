package algo

// ListNode is a singly linked list node.
type ListNode struct {
	Val  int
	Next *ListNode
}

// NewList builds a linked list from a slice, in order.
func NewList(vals []int) *ListNode {
	dummy := &ListNode{}
	cur := dummy
	for _, v := range vals {
		cur.Next = &ListNode{Val: v}
		cur = cur.Next
	}
	return dummy.Next
}

// ToSlice flattens a linked list back into a slice, for easy comparison.
func (n *ListNode) ToSlice() []int {
	var out []int
	for cur := n; cur != nil; cur = cur.Next {
		out = append(out, cur.Val)
	}
	return out
}

// ReverseList reverses a singly linked list in place, iteratively, and
// returns the new head. O(n) time, O(1) space.
func ReverseList(head *ListNode) *ListNode {
	var prev *ListNode
	cur := head
	for cur != nil {
		next := cur.Next
		cur.Next = prev
		prev = cur
		cur = next
	}
	return prev
}
