package algo

import (
	"reflect"
	"testing"
)

func TestRotateMatrix90(t *testing.T) {
	m := [][]int{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
	}
	want := [][]int{
		{7, 4, 1},
		{8, 5, 2},
		{9, 6, 3},
	}
	RotateMatrix90(m)
	if !reflect.DeepEqual(m, want) {
		t.Errorf("RotateMatrix90() = %v, want %v", m, want)
	}
}

func TestRotateMatrix90FourTimesIsIdentity(t *testing.T) {
	original := [][]int{
		{1, 2},
		{3, 4},
	}
	m := [][]int{{1, 2}, {3, 4}}
	for i := 0; i < 4; i++ {
		RotateMatrix90(m)
	}
	if !reflect.DeepEqual(m, original) {
		t.Errorf("four 90-degree rotations should return to the original: got %v, want %v", m, original)
	}
}

func TestRotateMatrix90SingleElement(t *testing.T) {
	m := [][]int{{1}}
	RotateMatrix90(m)
	if !reflect.DeepEqual(m, [][]int{{1}}) {
		t.Errorf("1x1 matrix should be unchanged, got %v", m)
	}
}
