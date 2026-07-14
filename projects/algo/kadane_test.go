package algo

import (
	"math"
	"testing"
)

// bruteForceMaxSubarraySum is the O(n^2) reference, used only to cross-check
// Kadane's result rather than trusting hand-picked expected values alone.
func bruteForceMaxSubarraySum(nums []int) int {
	best := math.MinInt64
	for i := range nums {
		sum := 0
		for j := i; j < len(nums); j++ {
			sum += nums[j]
			if sum > best {
				best = sum
			}
		}
	}
	return best
}

func TestMaxSubarraySumAgreesWithBruteForce(t *testing.T) {
	cases := [][]int{
		{-2, 1, -3, 4, -1, 2, 1, -5, 4},
		{1},
		{-1},
		{-2, -1, -3},
		{5, 4, -1, 7, 8},
		{0, 0, 0},
	}
	for _, nums := range cases {
		got := MaxSubarraySum(nums)
		want := bruteForceMaxSubarraySum(nums)
		if got != want {
			t.Errorf("MaxSubarraySum(%v) = %d, want %d (brute force)", nums, got, want)
		}
	}
}

func TestMaxSubarraySumEmpty(t *testing.T) {
	if got := MaxSubarraySum(nil); got != 0 {
		t.Errorf("MaxSubarraySum(nil) = %d, want 0", got)
	}
}
