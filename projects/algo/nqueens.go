package algo

// SolveNQueens returns every solution to the n-queens problem: each
// solution is a slice where result[row] = column of the queen in that row.
// Backtracking with O(1) conflict checks per placement, using sets for
// occupied columns and the two diagonal directions (row-col and row+col
// are each constant along a diagonal) instead of rescanning the board.
func SolveNQueens(n int) [][]int {
	var solutions [][]int
	placement := make([]int, n)
	cols := make(map[int]bool)
	diag1 := make(map[int]bool) // row - col
	diag2 := make(map[int]bool) // row + col

	var backtrack func(row int)
	backtrack = func(row int) {
		if row == n {
			solutions = append(solutions, append([]int(nil), placement...))
			return
		}
		for col := 0; col < n; col++ {
			if cols[col] || diag1[row-col] || diag2[row+col] {
				continue
			}
			placement[row] = col
			cols[col], diag1[row-col], diag2[row+col] = true, true, true
			backtrack(row + 1)
			cols[col], diag1[row-col], diag2[row+col] = false, false, false
		}
	}
	backtrack(0)
	return solutions
}
