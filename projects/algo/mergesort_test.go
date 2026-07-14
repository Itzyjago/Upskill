package algo

import (
	"reflect"
	"testing"
)

func TestMergeSort(t *testing.T) {
	cases := []struct {
		in   []int
		want []int
	}{
		{[]int{5, 3, 1, 4, 2}, []int{1, 2, 3, 4, 5}},
		{[]int{1}, []int{1}},
		{[]int{}, []int{}},
		{[]int{2, 2, 1}, []int{1, 2, 2}},
	}
	for _, c := range cases {
		if got := MergeSort(c.in); !reflect.DeepEqual(got, c.want) {
			t.Errorf("MergeSort(%v) = %v, want %v", c.in, got, c.want)
		}
	}
}

func TestMergeSortDoesNotMutateInput(t *testing.T) {
	in := []int{3, 1, 2}
	_ = MergeSort(in)
	want := []int{3, 1, 2}
	if !reflect.DeepEqual(in, want) {
		t.Errorf("MergeSort mutated input: got %v, want %v", in, want)
	}
}
