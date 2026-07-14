package algo

// QuickSelect returns the k-th smallest element (0-indexed) of nums without
// fully sorting, using the same Lomuto partition as QuickSort but only
// recursing into the one side that can contain index k. Average O(n) time
// — each partition step throws away the side that can't possibly hold the
// answer, unlike a full sort's O(n log n).
func QuickSelect(nums []int, k int) int {
	nums = append([]int(nil), nums...) // don't mutate the caller's slice
	lo, hi := 0, len(nums)-1
	for {
		if lo == hi {
			return nums[lo]
		}
		p := partition(nums, lo, hi)
		switch {
		case k == p:
			return nums[p]
		case k < p:
			hi = p - 1
		default:
			lo = p + 1
		}
	}
}
