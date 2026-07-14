package algo

import (
	"reflect"
	"testing"
)

func TestReverseList(t *testing.T) {
	cases := []struct {
		in   []int
		want []int
	}{
		{[]int{1, 2, 3, 4}, []int{4, 3, 2, 1}},
		{[]int{1}, []int{1}},
		{[]int{}, nil},
	}
	for _, c := range cases {
		head := NewList(c.in)
		reversed := ReverseList(head)
		got := reversed.ToSlice()
		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("ReverseList(%v) = %v, want %v", c.in, got, c.want)
		}
	}
}

func TestNewListAndToSliceRoundTrip(t *testing.T) {
	in := []int{7, 8, 9}
	got := NewList(in).ToSlice()
	if !reflect.DeepEqual(got, in) {
		t.Errorf("round trip = %v, want %v", got, in)
	}
}
