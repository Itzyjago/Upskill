package algo

import "testing"

func TestFenwickTreePrefixAndRangeSum(t *testing.T) {
	vals := []int{3, 2, -1, 6, 5, 4, -3, 3, 7, 2}
	f := NewFenwickTree(len(vals))
	for i, v := range vals {
		f.Add(i, v)
	}

	naivePrefix := func(i int) int {
		sum := 0
		for j := 0; j <= i; j++ {
			sum += vals[j]
		}
		return sum
	}
	for i := range vals {
		if got, want := f.PrefixSum(i), naivePrefix(i); got != want {
			t.Errorf("PrefixSum(%d) = %d, want %d", i, got, want)
		}
	}
	if got, want := f.RangeSum(2, 5), naivePrefix(5)-naivePrefix(1); got != want {
		t.Errorf("RangeSum(2, 5) = %d, want %d", got, want)
	}
}

func TestFenwickTreeUpdateReflectsInPrefixSum(t *testing.T) {
	f := NewFenwickTree(5)
	for i := 0; i < 5; i++ {
		f.Add(i, 1)
	}
	if got := f.PrefixSum(4); got != 5 {
		t.Fatalf("PrefixSum(4) = %d, want 5", got)
	}
	f.Add(2, 10) // vals become [1,1,11,1,1]
	if got := f.PrefixSum(4); got != 15 {
		t.Errorf("PrefixSum(4) after update = %d, want 15", got)
	}
	if got := f.RangeSum(0, 1); got != 2 {
		t.Errorf("RangeSum(0, 1) after update = %d, want 2 (unaffected range)", got)
	}
}
