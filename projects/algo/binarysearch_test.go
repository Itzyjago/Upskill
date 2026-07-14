package algo

import "testing"

func TestBinarySearch(t *testing.T) {
	sorted := []int{1, 3, 5, 7, 9, 11, 13}
	cases := []struct {
		target int
		want   int
	}{
		{1, 0},
		{13, 6},
		{7, 3},
		{2, -1},
		{14, -1},
	}
	for _, c := range cases {
		if got := BinarySearch(sorted, c.target); got != c.want {
			t.Errorf("BinarySearch(%v, %d) = %d, want %d", sorted, c.target, got, c.want)
		}
	}
}

func TestBinarySearchEmpty(t *testing.T) {
	if got := BinarySearch(nil, 5); got != -1 {
		t.Errorf("BinarySearch(nil, 5) = %d, want -1", got)
	}
}
