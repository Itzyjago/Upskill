package algo

// TopoSort returns a topological ordering of g (a valid ordering where every
// edge u->v has u before v), or ok=false if g has a cycle. Uses iterative
// DFS with a three-color (white/gray/black) visited state to detect
// back-edges, which is what a cycle looks like in DFS terms.
func (g *Graph) TopoSort() (order []int, ok bool) {
	const (
		white = 0 // unvisited
		gray  = 1 // on the current DFS stack
		black = 2 // fully processed
	)
	color := make(map[int]int)
	for node := range g.adj {
		color[node] = white
	}

	var stack []int
	var visit func(int) bool
	visit = func(node int) bool {
		color[node] = gray
		for _, next := range g.adj[node] {
			switch color[next] {
			case gray:
				return false // back-edge: cycle
			case white:
				if !visit(next) {
					return false
				}
			}
		}
		color[node] = black
		stack = append(stack, node)
		return true
	}

	for node := range g.adj {
		if color[node] == white {
			if !visit(node) {
				return nil, false
			}
		}
	}

	// stack is in reverse postorder; reverse it for the forward topo order.
	order = make([]int, len(stack))
	for i, n := range stack {
		order[len(stack)-1-i] = n
	}
	return order, true
}
