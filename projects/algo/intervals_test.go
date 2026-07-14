package algo

import (
	"reflect"
	"testing"
)

func TestMergeIntervals(t *testing.T) {
	cases := []struct {
		in   []Interval
		want []Interval
	}{
		{
			[]Interval{{1, 3}, {2, 6}, {8, 10}, {15, 18}},
			[]Interval{{1, 6}, {8, 10}, {15, 18}},
		},
		{
			[]Interval{{1, 4}, {4, 5}}, // touching, not overlapping
			[]Interval{{1, 5}},
		},
		{
			[]Interval{{1, 2}, {3, 4}}, // disjoint
			[]Interval{{1, 2}, {3, 4}},
		},
		{
			[]Interval{{5, 8}, {1, 3}}, // out of order input
			[]Interval{{1, 3}, {5, 8}},
		},
	}
	for _, c := range cases {
		if got := MergeIntervals(c.in); !reflect.DeepEqual(got, c.want) {
			t.Errorf("MergeIntervals(%v) = %v, want %v", c.in, got, c.want)
		}
	}
}

func TestMergeIntervalsEmpty(t *testing.T) {
	if got := MergeIntervals(nil); got != nil {
		t.Errorf("MergeIntervals(nil) = %v, want nil", got)
	}
}
