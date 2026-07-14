// Package algo collects small, tested implementations of classic data
// structures and algorithms, built to make the roadmap's algorithms notes
// stick rather than staying textbook prose.
package algo

// BinarySearch returns the index of target in a sorted (ascending) slice,
// or -1 if not present. O(log n) time, O(1) space.
func BinarySearch(sorted []int, target int) int {
	lo, hi := 0, len(sorted)-1
	for lo <= hi {
		mid := lo + (hi-lo)/2
		switch {
		case sorted[mid] == target:
			return mid
		case sorted[mid] < target:
			lo = mid + 1
		default:
			hi = mid - 1
		}
	}
	return -1
}
