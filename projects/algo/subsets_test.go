package algo

import "testing"

func TestSubsetsCountAndDistinctness(t *testing.T) {
	nums := []int{1, 2, 3}
	got := Subsets(nums)
	if want := 8; len(got) != want { // 2^3
		t.Fatalf("len(Subsets(%v)) = %d, want %d", nums, len(got), want)
	}
	seen := make(map[string]bool)
	for _, s := range got {
		key := ""
		for _, v := range s {
			key += string(rune('0' + v))
		}
		if seen[key] {
			t.Errorf("duplicate subset: %v", s)
		}
		seen[key] = true
	}
	if !seen[""] {
		t.Error("power set should include the empty subset")
	}
	if !seen["123"] {
		t.Error("power set should include the full set")
	}
}

func TestSubsetsEmptyInput(t *testing.T) {
	got := Subsets(nil)
	if len(got) != 1 || len(got[0]) != 0 {
		t.Errorf("Subsets(nil) = %v, want a single empty subset", got)
	}
}
