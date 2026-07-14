package algo

// Subsets returns the power set of nums (all 2^n subsets, including the
// empty one), via backtracking: at each index, branch into "exclude" and
// "include" — the same choice/undo pattern as Permutations, just deciding
// membership instead of ordering.
func Subsets(nums []int) [][]int {
	var out [][]int
	var cur []int
	var backtrack func(idx int)
	backtrack = func(idx int) {
		if idx == len(nums) {
			out = append(out, append([]int(nil), cur...))
			return
		}
		backtrack(idx + 1) // exclude nums[idx]
		cur = append(cur, nums[idx])
		backtrack(idx + 1) // include nums[idx]
		cur = cur[:len(cur)-1]
	}
	backtrack(0)
	return out
}
