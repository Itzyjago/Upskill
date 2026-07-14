package algo

// LevelOrder returns a binary tree's nodes level by level (BFS), reusing
// the generic Queue built earlier instead of a bespoke slice-based one.
func (n *BSTNode) LevelOrder() [][]int {
	if n == nil {
		return nil
	}
	var levels [][]int
	var q Queue[*BSTNode]
	q.Enqueue(n)
	for q.Len() > 0 {
		size := q.Len()
		var level []int
		for i := 0; i < size; i++ {
			node, _ := q.Dequeue()
			level = append(level, node.Val)
			if node.Left != nil {
				q.Enqueue(node.Left)
			}
			if node.Right != nil {
				q.Enqueue(node.Right)
			}
		}
		levels = append(levels, level)
	}
	return levels
}

// Height returns the tree's height (a single node has height 1, nil has
// height 0) — the same quantity AVLNode tracks incrementally, computed
// here directly for a plain BSTNode instead.
func (n *BSTNode) Height() int {
	if n == nil {
		return 0
	}
	l, r := n.Left.Height(), n.Right.Height()
	if l > r {
		return l + 1
	}
	return r + 1
}
