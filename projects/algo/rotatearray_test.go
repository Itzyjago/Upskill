package algo

import (
	"reflect"
	"testing"
)

func TestRotateArray(t *testing.T) {
	cases := []struct {
		in   []int
		k    int
		want []int
	}{
		{[]int{1, 2, 3, 4, 5, 6, 7}, 3, []int{5, 6, 7, 1, 2, 3, 4}},
		{[]int{1, 2}, 3, []int{2, 1}}, // k > len(nums)
		{[]int{1, 2, 3}, 0, []int{1, 2, 3}},
	}
	for _, c := range cases {
		got := append([]int(nil), c.in...)
		RotateArray(got, c.k)
		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("RotateArray(%v, %d) = %v, want %v", c.in, c.k, got, c.want)
		}
	}
}

func TestRotateArrayEmpty(t *testing.T) {
	var nums []int
	RotateArray(nums, 5) // must not panic on an empty slice
	if len(nums) != 0 {
		t.Errorf("RotateArray on empty slice should stay empty, got %v", nums)
	}
}

func TestRotateArrayFullRotationIsIdentity(t *testing.T) {
	nums := []int{1, 2, 3, 4, 5}
	want := append([]int(nil), nums...)
	RotateArray(nums, len(nums))
	if !reflect.DeepEqual(nums, want) {
		t.Errorf("RotateArray by len(nums) should be identity, got %v, want %v", nums, want)
	}
}
