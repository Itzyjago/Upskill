package algo

import "testing"

func TestBloomFilterNoFalseNegatives(t *testing.T) {
	bf := NewBloomFilter(1024, 4)
	words := []string{"apple", "banana", "cherry", "date", "elderberry"}
	for _, w := range words {
		bf.Add(w)
	}
	for _, w := range words {
		if !bf.Contains(w) {
			t.Errorf("Contains(%q) = false, want true (added earlier — false negatives aren't allowed)", w)
		}
	}
}

func TestBloomFilterAbsentItemsUsuallyNotContained(t *testing.T) {
	bf := NewBloomFilter(1024, 4)
	bf.Add("apple")
	bf.Add("banana")

	// A large, clearly-unrelated set at this size/k should mostly report
	// absent — false positives are possible in principle but should be rare,
	// not the norm, at this bits-per-item ratio.
	falsePositives := 0
	total := 200
	for i := 0; i < total; i++ {
		if bf.Contains(string(rune('a'+i%26)) + "-not-added-" + string(rune(i))) {
			falsePositives++
		}
	}
	if falsePositives > total/4 {
		t.Errorf("falsePositives = %d/%d, expected well under 25%% at this size/k", falsePositives, total)
	}
}

func TestBloomFilterEmpty(t *testing.T) {
	bf := NewBloomFilter(64, 3)
	if bf.Contains("anything") {
		t.Error("Contains on empty filter should be false")
	}
}
