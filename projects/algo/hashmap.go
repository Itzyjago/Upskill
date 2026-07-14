package algo

import "hash/fnv"

// HashMap is a from-scratch string->int hash map using separate chaining
// (each bucket is a slice of entries) to resolve collisions — built to
// make Go's built-in `map`'s behavior (which this project already relies
// on everywhere else, e.g. LRUCache/LFUCache) concrete rather than opaque.
// Amortized O(1) Get/Put/Delete given a reasonable load factor; grows by
// doubling and rehashing when the load factor gets too high, same
// approach as Go's runtime map.
type HashMap struct {
	buckets [][]hashMapEntry
	count   int
}

type hashMapEntry struct {
	key string
	val int
}

func NewHashMap() *HashMap {
	return &HashMap{buckets: make([][]hashMapEntry, 8)}
}

func (m *HashMap) bucketIndex(key string) int {
	h := fnv.New32a()
	h.Write([]byte(key))
	return int(h.Sum32()) % len(m.buckets)
}

func (m *HashMap) Put(key string, val int) {
	if float64(m.count+1)/float64(len(m.buckets)) > 0.75 {
		m.grow()
	}
	idx := m.bucketIndex(key)
	for i, e := range m.buckets[idx] {
		if e.key == key {
			m.buckets[idx][i].val = val
			return
		}
	}
	m.buckets[idx] = append(m.buckets[idx], hashMapEntry{key: key, val: val})
	m.count++
}

func (m *HashMap) Get(key string) (int, bool) {
	idx := m.bucketIndex(key)
	for _, e := range m.buckets[idx] {
		if e.key == key {
			return e.val, true
		}
	}
	return 0, false
}

func (m *HashMap) Delete(key string) bool {
	idx := m.bucketIndex(key)
	for i, e := range m.buckets[idx] {
		if e.key == key {
			m.buckets[idx] = append(m.buckets[idx][:i], m.buckets[idx][i+1:]...)
			m.count--
			return true
		}
	}
	return false
}

func (m *HashMap) Len() int { return m.count }

func (m *HashMap) grow() {
	old := m.buckets
	m.buckets = make([][]hashMapEntry, len(old)*2)
	m.count = 0
	for _, bucket := range old {
		for _, e := range bucket {
			m.Put(e.key, e.val)
		}
	}
}
