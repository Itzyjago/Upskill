package algo

import (
	"sort"
	"testing"
)

func TestPermutationsCountAndDistinctness(t *testing.T) {
	nums := []int{1, 2, 3, 4}
	perms := Permutations(nums)
	if got, want := len(perms), 24; got != want { // 4! = 24
		t.Fatalf("len(Permutations(%v)) = %d, want %d", nums, got, want)
	}
	seen := make(map[string]bool)
	for _, p := range perms {
		key := ""
		for _, v := range p {
			key += string(rune('0' + v))
		}
		if seen[key] {
			t.Errorf("duplicate permutation: %v", p)
		}
		seen[key] = true

		sorted := append([]int(nil), p...)
		sort.Ints(sorted)
		for i, v := range sorted {
			if v != nums[i] {
				t.Errorf("permutation %v is not a rearrangement of %v", p, nums)
				break
			}
		}
	}
}

func TestPermutationsSingleElement(t *testing.T) {
	perms := Permutations([]int{7})
	if len(perms) != 1 || perms[0][0] != 7 {
		t.Errorf("Permutations([7]) = %v, want [[7]]", perms)
	}
}
