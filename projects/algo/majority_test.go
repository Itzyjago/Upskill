package algo

import "testing"

func TestMajorityElement(t *testing.T) {
	cases := []struct {
		in   []int
		want int
	}{
		{[]int{3, 2, 3}, 3},
		{[]int{2, 2, 1, 1, 1, 2, 2}, 2},
		{[]int{1}, 1},
		{[]int{5, 5, 5, 1, 2}, 5},
	}
	for _, c := range cases {
		got, ok := MajorityElement(c.in)
		if !ok {
			t.Fatalf("MajorityElement(%v) ok = false, want true", c.in)
		}
		if got != c.want {
			t.Errorf("MajorityElement(%v) = %d, want %d", c.in, got, c.want)
		}
	}
}

func TestMajorityElementEmpty(t *testing.T) {
	if _, ok := MajorityElement(nil); ok {
		t.Error("MajorityElement(nil) ok = true, want false")
	}
}

// verifyMajority independently counts occurrences, to cross-check the
// voting algorithm's candidate rather than trusting it blindly.
func verifyMajority(nums []int, candidate int) bool {
	count := 0
	for _, n := range nums {
		if n == candidate {
			count++
		}
	}
	return count > len(nums)/2
}

func TestMajorityElementCandidateVerifiedByCounting(t *testing.T) {
	nums := []int{7, 7, 7, 7, 7, 1, 2, 3, 4}
	got, ok := MajorityElement(nums)
	if !ok || !verifyMajority(nums, got) {
		t.Errorf("MajorityElement(%v) = %d, %v — fails independent count verification", nums, got, ok)
	}
}
