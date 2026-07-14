package algo

// QuickSort sorts nums in place (ascending) using Lomuto partitioning.
// Average O(n log n), worst case O(n^2) on already-sorted input with this
// naive pivot choice (last element) — the tradeoff vs. MergeSort's
// guaranteed O(n log n) at the cost of extra space.
func QuickSort(nums []int) {
	quickSort(nums, 0, len(nums)-1)
}

func quickSort(nums []int, lo, hi int) {
	if lo >= hi {
		return
	}
	p := partition(nums, lo, hi)
	quickSort(nums, lo, p-1)
	quickSort(nums, p+1, hi)
}

func partition(nums []int, lo, hi int) int {
	pivot := nums[hi]
	i := lo
	for j := lo; j < hi; j++ {
		if nums[j] < pivot {
			nums[i], nums[j] = nums[j], nums[i]
			i++
		}
	}
	nums[i], nums[hi] = nums[hi], nums[i]
	return i
}
