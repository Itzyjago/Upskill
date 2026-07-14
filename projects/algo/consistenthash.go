package algo

import (
	"hash/fnv"
	"sort"
	"strconv"
)

// ConsistentHashRing distributes keys across a set of nodes so that adding
// or removing a node only reshuffles ~1/n of the keys, not all of them
// (unlike `hash(key) % nodeCount`, where changing nodeCount remaps almost
// everything). Each node gets `replicas` virtual points on the ring to
// smooth out an uneven distribution from an unlucky hash.
type ConsistentHashRing struct {
	replicas int
	ring     map[uint32]string
	sorted   []uint32
}

func NewConsistentHashRing(replicas int) *ConsistentHashRing {
	return &ConsistentHashRing{replicas: replicas, ring: make(map[uint32]string)}
}

func hashKey(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

func (r *ConsistentHashRing) AddNode(node string) {
	for i := 0; i < r.replicas; i++ {
		h := hashKey(node + "#" + strconv.Itoa(i))
		r.ring[h] = node
		r.sorted = append(r.sorted, h)
	}
	sort.Slice(r.sorted, func(i, j int) bool { return r.sorted[i] < r.sorted[j] })
}

func (r *ConsistentHashRing) RemoveNode(node string) {
	kept := r.sorted[:0]
	for _, h := range r.sorted {
		if r.ring[h] == node {
			delete(r.ring, h)
			continue
		}
		kept = append(kept, h)
	}
	r.sorted = kept
}

// GetNode returns the node responsible for key: walk clockwise from key's
// hash to the first ring point >= it, wrapping around to the first point
// if key's hash is past every point.
func (r *ConsistentHashRing) GetNode(key string) (string, bool) {
	if len(r.sorted) == 0 {
		return "", false
	}
	h := hashKey(key)
	i := sort.Search(len(r.sorted), func(i int) bool { return r.sorted[i] >= h })
	if i == len(r.sorted) {
		i = 0
	}
	return r.ring[r.sorted[i]], true
}
