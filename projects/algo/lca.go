package algo

// LowestCommonAncestor finds the lowest common ancestor of p and q in a
// BST. O(h) time (h = tree height) — a BST's ordering means the split
// point where p and q diverge to different subtrees *is* the LCA, no need
// to search both subtrees like a plain binary tree would require.
func (n *BSTNode) LowestCommonAncestor(p, q int) *BSTNode {
	cur := n
	for cur != nil {
		switch {
		case p < cur.Val && q < cur.Val:
			cur = cur.Left
		case p > cur.Val && q > cur.Val:
			cur = cur.Right
		default:
			return cur // p and q are on different sides (or one equals cur.Val)
		}
	}
	return nil
}
