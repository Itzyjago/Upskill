package algo

// RotateMatrix90 rotates an n x n matrix 90 degrees clockwise, in place.
// O(n^2) time, O(1) extra space: transpose (flip across the main diagonal),
// then reverse each row — two well-known O(n^2) steps instead of
// allocating a second matrix.
func RotateMatrix90(m [][]int) {
	n := len(m)
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			m[i][j], m[j][i] = m[j][i], m[i][j]
		}
	}
	for i := 0; i < n; i++ {
		for l, r := 0, n-1; l < r; l, r = l+1, r-1 {
			m[i][l], m[i][r] = m[i][r], m[i][l]
		}
	}
}
