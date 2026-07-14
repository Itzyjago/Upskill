package algo

import "reflect"
import "testing"

func TestMaxSlidingWindow(t *testing.T) {
	cases := []struct {
		nums []int
		k    int
		want []int
	}{
		{[]int{1, 3, -1, -3, 5, 3, 6, 7}, 3, []int{3, 3, 5, 5, 6, 7}},
		{[]int{9, 8, 7}, 1, []int{9, 8, 7}},
		{[]int{4}, 1, []int{4}},
		{[]int{1, 2, 3}, 3, []int{3}},
	}
	for _, c := range cases {
		if got := MaxSlidingWindow(c.nums, c.k); !reflect.DeepEqual(got, c.want) {
			t.Errorf("MaxSlidingWindow(%v, %d) = %v, want %v", c.nums, c.k, got, c.want)
		}
	}
}

func TestMaxSlidingWindowInvalidK(t *testing.T) {
	if got := MaxSlidingWindow([]int{1, 2, 3}, 0); got != nil {
		t.Errorf("k=0 should return nil, got %v", got)
	}
	if got := MaxSlidingWindow([]int{1, 2, 3}, 4); got != nil {
		t.Errorf("k > len(nums) should return nil, got %v", got)
	}
}
