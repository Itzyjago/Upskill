package algo

// LongestCommonSubsequence returns the length of the longest subsequence
// common to a and b (not necessarily contiguous, order preserved). Classic
// 2D DP, O(len(a)*len(b)) time and space.
func LongestCommonSubsequence(a, b string) int {
	dp := make([][]int, len(a)+1)
	for i := range dp {
		dp[i] = make([]int, len(b)+1)
	}
	for i := 1; i <= len(a); i++ {
		for j := 1; j <= len(b); j++ {
			if a[i-1] == b[j-1] {
				dp[i][j] = dp[i-1][j-1] + 1
			} else if dp[i-1][j] > dp[i][j-1] {
				dp[i][j] = dp[i-1][j]
			} else {
				dp[i][j] = dp[i][j-1]
			}
		}
	}
	return dp[len(a)][len(b)]
}

// EditDistance returns the Levenshtein distance between a and b — the
// minimum number of single-character insert/delete/replace operations to
// turn a into b. Same DP shape as LCS, different recurrence.
func EditDistance(a, b string) int {
	dp := make([][]int, len(a)+1)
	for i := range dp {
		dp[i] = make([]int, len(b)+1)
		dp[i][0] = i
	}
	for j := 0; j <= len(b); j++ {
		dp[0][j] = j
	}
	for i := 1; i <= len(a); i++ {
		for j := 1; j <= len(b); j++ {
			if a[i-1] == b[j-1] {
				dp[i][j] = dp[i-1][j-1]
				continue
			}
			min := dp[i-1][j] // delete
			if dp[i][j-1] < min {
				min = dp[i][j-1] // insert
			}
			if dp[i-1][j-1] < min {
				min = dp[i-1][j-1] // replace
			}
			dp[i][j] = min + 1
		}
	}
	return dp[len(a)][len(b)]
}
