package algo

import "testing"

// knownSolutionCounts are OEIS A000170 — the number of distinct n-queens
// solutions for small n, a well-known independently-verifiable sequence
// to check the backtracking search against, rather than trusting a
// hand-run of the algorithm itself.
var knownSolutionCounts = map[int]int{
	1: 1, 2: 0, 3: 0, 4: 2, 5: 10, 6: 4, 7: 40, 8: 92,
}

func TestSolveNQueensMatchesKnownCounts(t *testing.T) {
	for n, want := range knownSolutionCounts {
		got := len(SolveNQueens(n))
		if got != want {
			t.Errorf("len(SolveNQueens(%d)) = %d, want %d", n, got, want)
		}
	}
}

func TestSolveNQueensSolutionsAreValid(t *testing.T) {
	n := 8
	for _, sol := range SolveNQueens(n) {
		if len(sol) != n {
			t.Fatalf("solution has %d rows, want %d", len(sol), n)
		}
		for r1 := 0; r1 < n; r1++ {
			for r2 := r1 + 1; r2 < n; r2++ {
				c1, c2 := sol[r1], sol[r2]
				if c1 == c2 {
					t.Errorf("solution %v: rows %d,%d share column %d", sol, r1, r2, c1)
				}
				if r1-c1 == r2-c2 || r1+c1 == r2+c2 {
					t.Errorf("solution %v: rows %d,%d share a diagonal", sol, r1, r2)
				}
			}
		}
	}
}
