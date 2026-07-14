package algo

// BSTNode is a node in an (unbalanced) binary search tree.
type BSTNode struct {
	Val         int
	Left, Right *BSTNode
}

// Insert adds v to the tree rooted at n, returning the (possibly new) root.
// Duplicate values are ignored.
func (n *BSTNode) Insert(v int) *BSTNode {
	if n == nil {
		return &BSTNode{Val: v}
	}
	switch {
	case v < n.Val:
		n.Left = n.Left.Insert(v)
	case v > n.Val:
		n.Right = n.Right.Insert(v)
	}
	return n
}

// Contains reports whether v is in the tree. O(h) time, h = tree height —
// O(log n) balanced, O(n) worst case on sorted-order inserts.
func (n *BSTNode) Contains(v int) bool {
	if n == nil {
		return false
	}
	switch {
	case v == n.Val:
		return true
	case v < n.Val:
		return n.Left.Contains(v)
	default:
		return n.Right.Contains(v)
	}
}

// InOrder returns the tree's values in sorted order — the defining property
// of a BST's in-order traversal.
func (n *BSTNode) InOrder() []int {
	if n == nil {
		return nil
	}
	out := n.Left.InOrder()
	out = append(out, n.Val)
	out = append(out, n.Right.InOrder()...)
	return out
}
