package algo

import "sort"

// LongestIncreasingSubsequence returns the length of the longest strictly
// increasing subsequence of nums. O(n log n): maintains tails[k] = the
// smallest possible tail value of an increasing subsequence of length k+1,
// and binary-searches it on each element — the O(n^2) DP version compares
// every pair, this only ever does a log-n search per element.
func LongestIncreasingSubsequence(nums []int) int {
	var tails []int
	for _, v := range nums {
		i := sort.SearchInts(tails, v)
		if i == len(tails) {
			tails = append(tails, v)
		} else {
			tails[i] = v
		}
	}
	return len(tails)
}
