package algo

// TwoSum returns indices i, j (i < j) such that nums[i]+nums[j] == target,
// or ok=false if no such pair exists. O(n) time via a single-pass hash map
// of value -> index, vs. the naive O(n^2) pair scan.
func TwoSum(nums []int, target int) (i, j int, ok bool) {
	seen := make(map[int]int, len(nums)) // value -> index
	for idx, v := range nums {
		need := target - v
		if prevIdx, found := seen[need]; found {
			return prevIdx, idx, true
		}
		seen[v] = idx
	}
	return 0, 0, false
}
