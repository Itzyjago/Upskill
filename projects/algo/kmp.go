package algo

// KMPSearch finds all starting indices of pattern in text using the
// Knuth-Morris-Pratt algorithm. O(n+m) time, vs. naive search's O(n*m) —
// the failure function lets a mismatch skip ahead without rescanning text
// already known to match.
func KMPSearch(text, pattern string) []int {
	if len(pattern) == 0 {
		return nil
	}
	lps := buildLPS(pattern)
	var matches []int
	i, j := 0, 0 // i = text index, j = pattern index
	for i < len(text) {
		if text[i] == pattern[j] {
			i++
			j++
			if j == len(pattern) {
				matches = append(matches, i-j)
				j = lps[j-1]
			}
		} else if j > 0 {
			j = lps[j-1]
		} else {
			i++
		}
	}
	return matches
}

// buildLPS computes the "longest proper prefix that's also a suffix" array
// for pattern — the backbone of KMP's skip logic.
func buildLPS(pattern string) []int {
	lps := make([]int, len(pattern))
	length := 0
	i := 1
	for i < len(pattern) {
		if pattern[i] == pattern[length] {
			length++
			lps[i] = length
			i++
		} else if length > 0 {
			length = lps[length-1]
		} else {
			lps[i] = 0
			i++
		}
	}
	return lps
}
