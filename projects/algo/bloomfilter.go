package algo

import "hash/fnv"

// BloomFilter is a probabilistic set: Contains never false-negatives (if it
// says "not present", it's really not present) but can false-positive (it
// may say "present" for something never added). Trades exactness for O(k)
// space-efficient membership tests, k = number of hash functions.
type BloomFilter struct {
	bits []bool
	k    int // number of hash functions (derived positions per Add/Contains)
}

func NewBloomFilter(size, k int) *BloomFilter {
	return &BloomFilter{bits: make([]bool, size), k: k}
}

func (b *BloomFilter) Add(s string) {
	for _, pos := range b.positions(s) {
		b.bits[pos] = true
	}
}

func (b *BloomFilter) Contains(s string) bool {
	for _, pos := range b.positions(s) {
		if !b.bits[pos] {
			return false
		}
	}
	return true
}

// positions derives k bit positions from two independent hashes of s
// (double hashing: h1 + i*h2), the standard way to simulate k hash
// functions without computing k actual hash functions.
func (b *BloomFilter) positions(s string) []int {
	h1 := fnv.New64a()
	h1.Write([]byte(s))
	sum1 := h1.Sum64()

	h2 := fnv.New32a()
	h2.Write([]byte(s))
	sum2 := uint64(h2.Sum32())

	out := make([]int, b.k)
	for i := 0; i < b.k; i++ {
		out[i] = int((sum1 + uint64(i)*sum2) % uint64(len(b.bits)))
	}
	return out
}
