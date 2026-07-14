package algo

import "testing"

func TestTwoSum(t *testing.T) {
	nums := []int{2, 7, 11, 15}
	i, j, ok := TwoSum(nums, 9)
	if !ok || nums[i]+nums[j] != 9 {
		t.Fatalf("TwoSum(%v, 9) = %d, %d, %v; nums[i]+nums[j] should be 9", nums, i, j, ok)
	}
	if i != 0 || j != 1 {
		t.Errorf("TwoSum(%v, 9) = %d, %d, want 0, 1", nums, i, j)
	}
}

func TestTwoSumNoSolution(t *testing.T) {
	if _, _, ok := TwoSum([]int{1, 2, 3}, 100); ok {
		t.Error("TwoSum should return ok=false when no pair sums to target")
	}
}

func TestTwoSumUsesEarliestPair(t *testing.T) {
	// Two candidate pairs sum to target; TwoSum should return the one found
	// first scanning left to right (3+3 at indices 1,2 before 3+3 at 1,3).
	nums := []int{0, 3, 3, 3}
	i, j, ok := TwoSum(nums, 6)
	if !ok {
		t.Fatal("expected a solution")
	}
	if i != 1 || j != 2 {
		t.Errorf("TwoSum(%v, 6) = %d, %d, want 1, 2", nums, i, j)
	}
}
