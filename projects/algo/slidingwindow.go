package algo

// MaxSlidingWindow returns the max of every contiguous window of size k in
// nums. O(n) time using a monotonic deque of indices (decreasing values) —
// each index enters and leaves the deque at most once, vs. the naive
// O(n*k) of scanning every window from scratch.
func MaxSlidingWindow(nums []int, k int) []int {
	if k <= 0 || k > len(nums) {
		return nil
	}
	var deque []int // holds indices, values strictly decreasing front-to-back
	var out []int
	for i, v := range nums {
		for len(deque) > 0 && nums[deque[len(deque)-1]] <= v {
			deque = deque[:len(deque)-1]
		}
		deque = append(deque, i)
		if deque[0] <= i-k {
			deque = deque[1:]
		}
		if i >= k-1 {
			out = append(out, nums[deque[0]])
		}
	}
	return out
}
