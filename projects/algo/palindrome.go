package algo

// LongestPalindromicSubstring returns the longest substring of s that
// reads the same forwards and backwards. O(n^2) time via expand-around-
// center: every index (and every gap between two indices, for even-length
// palindromes) is tried as a center, expanding outward while it still
// matches.
func LongestPalindromicSubstring(s string) string {
	if len(s) == 0 {
		return ""
	}
	start, end := 0, 0
	expand := func(l, r int) (int, int) {
		for l >= 0 && r < len(s) && s[l] == s[r] {
			l--
			r++
		}
		return l + 1, r - 1
	}
	for i := 0; i < len(s); i++ {
		l1, r1 := expand(i, i)     // odd-length, center at i
		l2, r2 := expand(i, i+1)   // even-length, center between i and i+1
		if r1-l1 > end-start {
			start, end = l1, r1
		}
		if r2-l2 > end-start {
			start, end = l2, r2
		}
	}
	return s[start : end+1]
}
