package algo

import "container/heap"

// WeightedGraph is a directed graph with non-negative edge weights.
type WeightedGraph struct {
	adj map[int][]weightedEdge
}

type weightedEdge struct {
	to     int
	weight int
}

func NewWeightedGraph() *WeightedGraph {
	return &WeightedGraph{adj: make(map[int][]weightedEdge)}
}

// AddEdge adds a directed edge. weight must be >= 0 — Dijkstra's algorithm
// assumes no negative weights; use Bellman-Ford instead if that's needed.
func (g *WeightedGraph) AddEdge(from, to, weight int) {
	g.adj[from] = append(g.adj[from], weightedEdge{to: to, weight: weight})
	if _, ok := g.adj[to]; !ok {
		g.adj[to] = nil
	}
}

type pqItem struct {
	node int
	dist int
}

// dijkstraPQ implements container/heap.Interface as a min-heap on dist.
type dijkstraPQ []pqItem

func (pq dijkstraPQ) Len() int            { return len(pq) }
func (pq dijkstraPQ) Less(i, j int) bool  { return pq[i].dist < pq[j].dist }
func (pq dijkstraPQ) Swap(i, j int)       { pq[i], pq[j] = pq[j], pq[i] }
func (pq *dijkstraPQ) Push(x interface{}) { *pq = append(*pq, x.(pqItem)) }
func (pq *dijkstraPQ) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[:n-1]
	return item
}

// ShortestPaths returns the shortest distance from start to every reachable
// node. O((V+E) log V) using a binary heap as the priority queue.
func (g *WeightedGraph) ShortestPaths(start int) map[int]int {
	dist := map[int]int{start: 0}
	pq := &dijkstraPQ{{node: start, dist: 0}}
	heap.Init(pq)

	for pq.Len() > 0 {
		cur := heap.Pop(pq).(pqItem)
		if d, ok := dist[cur.node]; ok && cur.dist > d {
			continue // stale entry, a shorter path to this node already won
		}
		for _, e := range g.adj[cur.node] {
			next := cur.dist + e.weight
			if d, ok := dist[e.to]; !ok || next < d {
				dist[e.to] = next
				heap.Push(pq, pqItem{node: e.to, dist: next})
			}
		}
	}
	return dist
}
