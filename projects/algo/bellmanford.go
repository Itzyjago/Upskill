package algo

import "math"

// BellmanFord computes shortest distances from start over directed edges
// (reuses Edge from mst.go — From/To/Weight, here read as directed rather
// than undirected). Tolerates negative weights, unlike Dijkstra. O(V*E) —
// relax every edge V-1 times, then one more pass to detect a
// negative-weight cycle reachable from start (if any edge can still relax,
// a cycle keeps improving forever).
func BellmanFord(n int, edges []Edge, start int) (dist []int, negativeCycle bool) {
	dist = make([]int, n)
	for i := range dist {
		dist[i] = math.MaxInt64
	}
	dist[start] = 0

	for i := 0; i < n-1; i++ {
		for _, e := range edges {
			if dist[e.From] == math.MaxInt64 {
				continue
			}
			if next := dist[e.From] + e.Weight; next < dist[e.To] {
				dist[e.To] = next
			}
		}
	}

	for _, e := range edges {
		if dist[e.From] == math.MaxInt64 {
			continue
		}
		if dist[e.From]+e.Weight < dist[e.To] {
			return dist, true
		}
	}
	return dist, false
}
