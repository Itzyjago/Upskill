package algo

// WordBreak reports whether s can be segmented into a space-separated
// sequence of words from dict. O(len(s)^2) DP: dp[i] means "s[:i] can be
// fully segmented", built from smaller solved prefixes instead of
// re-trying every split from scratch (which is what makes the naive
// recursive version exponential).
func WordBreak(s string, dict []string) bool {
	words := make(map[string]bool, len(dict))
	for _, w := range dict {
		words[w] = true
	}
	dp := make([]bool, len(s)+1)
	dp[0] = true
	for i := 1; i <= len(s); i++ {
		for j := 0; j < i; j++ {
			if dp[j] && words[s[j:i]] {
				dp[i] = true
				break
			}
		}
	}
	return dp[len(s)]
}
