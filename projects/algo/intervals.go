package algo

import "sort"

// Interval is an inclusive [Start, End] range.
type Interval struct {
	Start, End int
}

// MergeIntervals sorts by start and merges any overlapping or touching
// intervals. O(n log n), dominated by the sort.
func MergeIntervals(intervals []Interval) []Interval {
	if len(intervals) == 0 {
		return nil
	}
	sorted := make([]Interval, len(intervals))
	copy(sorted, intervals)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].Start < sorted[j].Start })

	merged := []Interval{sorted[0]}
	for _, cur := range sorted[1:] {
		last := &merged[len(merged)-1]
		if cur.Start <= last.End {
			if cur.End > last.End {
				last.End = cur.End
			}
			continue
		}
		merged = append(merged, cur)
	}
	return merged
}
