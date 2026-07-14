package algo

import "testing"

func TestFibMemoAndTabulationAgree(t *testing.T) {
	want := []int{0, 1, 1, 2, 3, 5, 8, 13, 21, 34, 55}
	for n, w := range want {
		if got := FibMemo(n); got != w {
			t.Errorf("FibMemo(%d) = %d, want %d", n, got, w)
		}
		if got := FibTabulation(n); got != w {
			t.Errorf("FibTabulation(%d) = %d, want %d", n, got, w)
		}
	}
}

func TestFibLargerN(t *testing.T) {
	// F(30) = 832040, a value large enough that naive exponential
	// recursion would be noticeably slow, memoization/tabulation isn't.
	if got := FibMemo(30); got != 832040 {
		t.Errorf("FibMemo(30) = %d, want 832040", got)
	}
	if got := FibTabulation(30); got != 832040 {
		t.Errorf("FibTabulation(30) = %d, want 832040", got)
	}
}
