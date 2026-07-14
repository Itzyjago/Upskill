package algo

import "testing"

func TestSegmentTreeQuery(t *testing.T) {
	nums := []int{1, 3, 5, 7, 9, 11}
	st := NewSegmentTree(nums)

	naiveSum := func(lo, hi int) int {
		sum := 0
		for i := lo; i <= hi; i++ {
			sum += nums[i]
		}
		return sum
	}
	for lo := 0; lo < len(nums); lo++ {
		for hi := lo; hi < len(nums); hi++ {
			if got, want := st.Query(lo, hi), naiveSum(lo, hi); got != want {
				t.Errorf("Query(%d, %d) = %d, want %d", lo, hi, got, want)
			}
		}
	}
}

func TestSegmentTreeUpdateReflectsInQuery(t *testing.T) {
	nums := []int{1, 3, 5}
	st := NewSegmentTree(nums)
	if got := st.Query(0, 2); got != 9 {
		t.Fatalf("Query(0, 2) = %d, want 9", got)
	}
	st.Update(1, 100) // nums conceptually become [1, 100, 5]
	if got := st.Query(0, 2); got != 106 {
		t.Errorf("Query(0, 2) after update = %d, want 106", got)
	}
	if got := st.Query(0, 0); got != 1 {
		t.Errorf("Query(0, 0) after unrelated update = %d, want 1", got)
	}
}

func TestSegmentTreeSingleElement(t *testing.T) {
	st := NewSegmentTree([]int{42})
	if got := st.Query(0, 0); got != 42 {
		t.Errorf("Query(0, 0) = %d, want 42", got)
	}
}
