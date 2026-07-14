package algo

import "testing"

func TestCoinChange(t *testing.T) {
	cases := []struct {
		coins  []int
		amount int
		want   int
	}{
		{[]int{1, 2, 5}, 11, 3}, // 5+5+1
		{[]int{2}, 3, -1},
		{[]int{1}, 0, 0},
		{[]int{1}, 2, 2},
		{[]int{1, 3, 4}, 6, 2}, // 3+3, not 4+1+1 or 1*6
	}
	for _, c := range cases {
		if got := CoinChange(c.coins, c.amount); got != c.want {
			t.Errorf("CoinChange(%v, %d) = %d, want %d", c.coins, c.amount, got, c.want)
		}
	}
}
