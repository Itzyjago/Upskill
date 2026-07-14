package algo

// IntSqrt returns floor(sqrt(n)) for n >= 0, via binary search over the
// candidate range [0, n]. O(log n) time, no floating point involved (and
// so no float-precision edge cases near perfect squares).
func IntSqrt(n int) int {
	if n < 0 {
		panic("IntSqrt: negative input")
	}
	lo, hi := 0, n
	ans := 0
	for lo <= hi {
		mid := lo + (hi-lo)/2
		if mid*mid <= n {
			ans = mid
			lo = mid + 1
		} else {
			hi = mid - 1
		}
	}
	return ans
}
