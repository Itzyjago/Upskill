package algo

// Graph is an unweighted, directed adjacency list over integer node ids.
type Graph struct {
	adj map[int][]int
}

func NewGraph() *Graph {
	return &Graph{adj: make(map[int][]int)}
}

func (g *Graph) AddEdge(from, to int) {
	g.adj[from] = append(g.adj[from], to)
	if _, ok := g.adj[to]; !ok {
		g.adj[to] = nil // register the node even if it has no outgoing edges
	}
}

// BFS returns nodes reachable from start, in breadth-first visitation order.
func (g *Graph) BFS(start int) []int {
	visited := map[int]bool{start: true}
	queue := []int{start}
	var order []int
	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]
		order = append(order, node)
		for _, next := range g.adj[node] {
			if !visited[next] {
				visited[next] = true
				queue = append(queue, next)
			}
		}
	}
	return order
}

// DFS returns nodes reachable from start, in depth-first (preorder) visitation order.
func (g *Graph) DFS(start int) []int {
	visited := map[int]bool{}
	var order []int
	var visit func(int)
	visit = func(node int) {
		if visited[node] {
			return
		}
		visited[node] = true
		order = append(order, node)
		for _, next := range g.adj[node] {
			visit(next)
		}
	}
	visit(start)
	return order
}
