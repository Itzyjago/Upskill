package algo

import "sort"

// GroupAnagrams groups words that are anagrams of each other. O(n * k log k)
// time (n words, k = max word length) — sorting each word's letters gives a
// canonical key so anagrams collide in the same map bucket, no pairwise
// comparison needed.
func GroupAnagrams(words []string) [][]string {
	groups := make(map[string][]string)
	var keys []string // preserve first-seen order for deterministic output
	for _, w := range words {
		key := sortLetters(w)
		if _, ok := groups[key]; !ok {
			keys = append(keys, key)
		}
		groups[key] = append(groups[key], w)
	}
	out := make([][]string, len(keys))
	for i, k := range keys {
		out[i] = groups[k]
	}
	return out
}

func sortLetters(s string) string {
	b := []byte(s)
	sort.Slice(b, func(i, j int) bool { return b[i] < b[j] })
	return string(b)
}
