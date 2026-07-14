package algo

// Permutations returns every distinct ordering of nums (assumed to have no
// duplicate values), via backtracking: build one position at a time, swap
// candidates into place, recurse, then swap back to try the next
// candidate. O(n!) orderings, O(n) extra space beyond the output.
func Permutations(nums []int) [][]int {
	var out [][]int
	work := append([]int(nil), nums...)
	var backtrack func(k int)
	backtrack = func(k int) {
		if k == len(work) {
			out = append(out, append([]int(nil), work...))
			return
		}
		for i := k; i < len(work); i++ {
			work[k], work[i] = work[i], work[k]
			backtrack(k + 1)
			work[k], work[i] = work[i], work[k] // undo, so the next i sees the original order
		}
	}
	backtrack(0)
	return out
}
