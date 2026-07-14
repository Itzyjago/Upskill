package algo

import (
	"reflect"
	"testing"
)

func TestQuickSort(t *testing.T) {
	cases := []struct {
		in   []int
		want []int
	}{
		{[]int{5, 3, 1, 4, 2}, []int{1, 2, 3, 4, 5}},
		{[]int{1}, []int{1}},
		{[]int{}, []int{}},
		{[]int{9, 8, 7, 6}, []int{6, 7, 8, 9}},
		{[]int{1, 1, 1}, []int{1, 1, 1}},
	}
	for _, c := range cases {
		got := make([]int, len(c.in))
		copy(got, c.in)
		QuickSort(got)
		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("QuickSort(%v) = %v, want %v", c.in, got, c.want)
		}
	}
}
