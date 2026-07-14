package algo

// MaxSubarraySum returns the largest sum of any contiguous subarray of
// nums (Kadane's algorithm). O(n) time, O(1) space — at each index it
// decides whether extending the running subarray beats starting fresh
// there, which is enough to find the global best without ever comparing
// all O(n^2) subarrays.
func MaxSubarraySum(nums []int) int {
	if len(nums) == 0 {
		return 0
	}
	best := nums[0]
	cur := nums[0]
	for _, v := range nums[1:] {
		if cur+v > v {
			cur = cur + v
		} else {
			cur = v
		}
		if cur > best {
			best = cur
		}
	}
	return best
}
