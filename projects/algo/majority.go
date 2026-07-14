package algo

// MajorityElement returns the element appearing more than len(nums)/2
// times, using the Boyer-Moore voting algorithm: O(n) time, O(1) space —
// no hash map of counts needed. Works because a true majority element
// can't be fully "cancelled out" by every other element paired against it.
// ok is false if nums is empty; if no true majority exists the function
// still returns some value (the candidate), so callers who aren't sure one
// exists should verify it themselves by counting.
func MajorityElement(nums []int) (v int, ok bool) {
	if len(nums) == 0 {
		return 0, false
	}
	candidate := nums[0]
	count := 0
	for _, n := range nums {
		if count == 0 {
			candidate = n
		}
		if n == candidate {
			count++
		} else {
			count--
		}
	}
	return candidate, true
}
