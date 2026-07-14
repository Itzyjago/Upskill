package algo

// CountingSort sorts nums (ascending), assuming values fall in [0, max]
// inclusive. O(n+max) time, not comparison-based — it beats the
// n*log(n) lower bound that applies to comparison sorts (QuickSort,
// MergeSort) by using the values themselves as bucket indices instead of
// comparing pairs.
func CountingSort(nums []int, max int) []int {
	counts := make([]int, max+1)
	for _, v := range nums {
		counts[v]++
	}
	out := make([]int, 0, len(nums))
	for v, c := range counts {
		for i := 0; i < c; i++ {
			out = append(out, v)
		}
	}
	return out
}
