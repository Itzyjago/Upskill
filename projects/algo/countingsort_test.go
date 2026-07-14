package algo

import (
	"reflect"
	"testing"
)

func TestCountingSort(t *testing.T) {
	cases := []struct {
		in   []int
		max  int
		want []int
	}{
		{[]int{4, 2, 2, 8, 3, 3, 1}, 8, []int{1, 2, 2, 3, 3, 4, 8}},
		{[]int{0, 0, 0}, 0, []int{0, 0, 0}},
		{[]int{}, 5, []int{}},
	}
	for _, c := range cases {
		got := CountingSort(c.in, c.max)
		if len(got) == 0 && len(c.want) == 0 {
			continue // reflect.DeepEqual(nil, []int{}) is false, treat both as "empty"
		}
		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("CountingSort(%v, %d) = %v, want %v", c.in, c.max, got, c.want)
		}
	}
}
