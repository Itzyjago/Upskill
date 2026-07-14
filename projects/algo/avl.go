package algo

// AVLNode is a self-balancing BST node: after every insert, rotations keep
// the height difference between left and right subtrees at most 1, which
// is what keeps operations O(log n) instead of degrading to O(n) on
// sorted-order inserts like a plain BSTNode does.
type AVLNode struct {
	Val         int
	Left, Right *AVLNode
	height      int
}

func height(n *AVLNode) int {
	if n == nil {
		return 0
	}
	return n.height
}

func balanceFactor(n *AVLNode) int {
	if n == nil {
		return 0
	}
	return height(n.Left) - height(n.Right)
}

func updateHeight(n *AVLNode) {
	l, r := height(n.Left), height(n.Right)
	if l > r {
		n.height = l + 1
	} else {
		n.height = r + 1
	}
}

func rotateRight(y *AVLNode) *AVLNode {
	x := y.Left
	y.Left = x.Right
	x.Right = y
	updateHeight(y)
	updateHeight(x)
	return x
}

func rotateLeft(x *AVLNode) *AVLNode {
	y := x.Right
	x.Right = y.Left
	y.Left = x
	updateHeight(x)
	updateHeight(y)
	return y
}

// Insert adds v, rebalancing as needed, and returns the new subtree root.
func (n *AVLNode) Insert(v int) *AVLNode {
	if n == nil {
		return &AVLNode{Val: v, height: 1}
	}
	if v < n.Val {
		n.Left = n.Left.Insert(v)
	} else if v > n.Val {
		n.Right = n.Right.Insert(v)
	} else {
		return n // duplicate, ignore
	}
	updateHeight(n)

	bf := balanceFactor(n)
	switch {
	case bf > 1 && v < n.Left.Val: // left-left
		return rotateRight(n)
	case bf > 1: // left-right
		n.Left = rotateLeft(n.Left)
		return rotateRight(n)
	case bf < -1 && v > n.Right.Val: // right-right
		return rotateLeft(n)
	case bf < -1: // right-left
		n.Right = rotateRight(n.Right)
		return rotateLeft(n)
	}
	return n
}

func (n *AVLNode) InOrder() []int {
	if n == nil {
		return nil
	}
	out := n.Left.InOrder()
	out = append(out, n.Val)
	out = append(out, n.Right.InOrder()...)
	return out
}
