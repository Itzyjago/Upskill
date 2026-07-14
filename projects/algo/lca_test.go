package algo

import "testing"

func TestBSTLowestCommonAncestor(t *testing.T) {
	//          6
	//       /     \
	//      2       8
	//     / \     / \
	//    0   4   7   9
	//       / \
	//      3   5
	root := buildBST([]int{6, 2, 8, 0, 4, 7, 9, 3, 5})

	cases := []struct {
		p, q, want int
	}{
		{2, 8, 6}, // split at root
		{2, 4, 2}, // q is in p's subtree
		{3, 5, 4}, // both under 4
		{0, 5, 2}, // split at 2
		{7, 9, 8}, // split at 8
	}
	for _, c := range cases {
		got := root.LowestCommonAncestor(c.p, c.q)
		if got == nil || got.Val != c.want {
			gotVal := -1
			if got != nil {
				gotVal = got.Val
			}
			t.Errorf("LowestCommonAncestor(%d, %d) = %d, want %d", c.p, c.q, gotVal, c.want)
		}
	}
}
