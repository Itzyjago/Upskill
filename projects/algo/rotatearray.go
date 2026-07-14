package algo

// RotateArray rotates nums right by k positions, in place. O(n) time,
// O(1) extra space via the reverse trick: reverse the whole array, then
// reverse each of the two halves — three linear passes instead of
// allocating a second array to shift into.
func RotateArray(nums []int, k int) {
	n := len(nums)
	if n == 0 {
		return
	}
	k %= n
	if k < 0 {
		k += n
	}
	reverse(nums, 0, n-1)
	reverse(nums, 0, k-1)
	reverse(nums, k, n-1)
}

func reverse(nums []int, l, r int) {
	for l < r {
		nums[l], nums[r] = nums[r], nums[l]
		l++
		r--
	}
}
