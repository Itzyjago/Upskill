package algo

import "testing"

func TestLongestIncreasingSubsequence(t *testing.T) {
	cases := []struct {
		in   []int
		want int
	}{
		{[]int{10, 9, 2, 5, 3, 7, 101, 18}, 4}, // 2,3,7,101 or 2,3,7,18
		{[]int{0, 1, 0, 3, 2, 3}, 4},
		{[]int{7, 7, 7, 7}, 1}, // strictly increasing, ties don't extend
		{[]int{}, 0},
		{[]int{5}, 1},
	}
	for _, c := range cases {
		if got := LongestIncreasingSubsequence(c.in); got != c.want {
			t.Errorf("LongestIncreasingSubsequence(%v) = %d, want %d", c.in, got, c.want)
		}
	}
}

// bruteForceLIS is an O(2^n)-ish exhaustive reference for small inputs,
// used to cross-check the O(n log n) implementation instead of trusting
// hand-picked expected values alone.
func bruteForceLIS(nums []int) int {
	best := 0
	var rec func(idx int, prev int, curLen int)
	rec = func(idx, prev, curLen int) {
		if curLen > best {
			best = curLen
		}
		for i := idx; i < len(nums); i++ {
			if prev == -1<<63 || nums[i] > prev {
				rec(i+1, nums[i], curLen+1)
			}
		}
	}
	rec(0, -1<<63, 0)
	return best
}

func TestLongestIncreasingSubsequenceAgreesWithBruteForce(t *testing.T) {
	cases := [][]int{
		{4, 2, 3, 6, 10, 1, 12},
		{3, 1, 4, 1, 5, 9, 2, 6},
		{1, 2, 3, 4, 5},
		{5, 4, 3, 2, 1},
	}
	for _, nums := range cases {
		got := LongestIncreasingSubsequence(nums)
		want := bruteForceLIS(nums)
		if got != want {
			t.Errorf("LongestIncreasingSubsequence(%v) = %d, want %d (brute force)", nums, got, want)
		}
	}
}
