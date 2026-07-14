package algo

// HasCycle detects a cycle in a linked list using Floyd's tortoise-and-
// hare: a slow pointer advances one node per step, a fast pointer two —
// if there's a cycle they must eventually meet inside it, O(n) time,
// O(1) space, no visited-set needed.
func HasCycle(head *ListNode) bool {
	slow, fast := head, head
	for fast != nil && fast.Next != nil {
		slow = slow.Next
		fast = fast.Next.Next
		if slow == fast {
			return true
		}
	}
	return false
}

// MiddleNode returns the middle node of a linked list (the second of two
// middles for even length), using the same fast/slow pointer trick: when
// fast reaches the end, slow is exactly halfway.
func MiddleNode(head *ListNode) *ListNode {
	slow, fast := head, head
	for fast != nil && fast.Next != nil {
		slow = slow.Next
		fast = fast.Next.Next
	}
	return slow
}
