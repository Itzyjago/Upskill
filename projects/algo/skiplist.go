package algo

import "math/rand"

const skipListMaxLevel = 16

// SkipList is a probabilistic ordered set: each node gets a random "tower
// height", and searches drop down a level whenever going right would
// overshoot. Expected O(log n) search/insert/contains without any
// rebalancing logic (unlike AVLNode) — the randomness does the balancing
// work instead of explicit rotations.
type SkipList struct {
	head  *skipListNode
	level int
	rng   *rand.Rand
}

type skipListNode struct {
	val  int
	next []*skipListNode
}

func NewSkipList(rng *rand.Rand) *SkipList {
	return &SkipList{
		head:  &skipListNode{next: make([]*skipListNode, skipListMaxLevel)},
		level: 1,
		rng:   rng,
	}
}

func (s *SkipList) randomLevel() int {
	level := 1
	for level < skipListMaxLevel && s.rng.Intn(2) == 0 {
		level++
	}
	return level
}

func (s *SkipList) Insert(v int) {
	update := make([]*skipListNode, skipListMaxLevel)
	cur := s.head
	for i := s.level - 1; i >= 0; i-- {
		for cur.next[i] != nil && cur.next[i].val < v {
			cur = cur.next[i]
		}
		update[i] = cur
	}
	if cur.next[0] != nil && cur.next[0].val == v {
		return // duplicate, ignore
	}

	level := s.randomLevel()
	if level > s.level {
		for i := s.level; i < level; i++ {
			update[i] = s.head
		}
		s.level = level
	}
	node := &skipListNode{val: v, next: make([]*skipListNode, level)}
	for i := 0; i < level; i++ {
		node.next[i] = update[i].next[i]
		update[i].next[i] = node
	}
}

func (s *SkipList) Contains(v int) bool {
	cur := s.head
	for i := s.level - 1; i >= 0; i-- {
		for cur.next[i] != nil && cur.next[i].val < v {
			cur = cur.next[i]
		}
	}
	cur = cur.next[0]
	return cur != nil && cur.val == v
}

// ToSlice returns every value in sorted order, walking the bottom level —
// the same list every insert threads through regardless of tower height.
func (s *SkipList) ToSlice() []int {
	var out []int
	for cur := s.head.next[0]; cur != nil; cur = cur.next[0] {
		out = append(out, cur.val)
	}
	return out
}
