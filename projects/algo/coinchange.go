package algo

// CoinChange returns the minimum number of coins from coins needed to make
// amount, or -1 if impossible. Bottom-up DP, O(amount * len(coins)) — each
// amount from 1..target is built from a smaller already-solved amount plus
// one coin, so no amount is ever recomputed from scratch.
func CoinChange(coins []int, amount int) int {
	const unreachable = 1 << 30
	dp := make([]int, amount+1)
	for i := 1; i <= amount; i++ {
		dp[i] = unreachable
	}
	for i := 1; i <= amount; i++ {
		for _, c := range coins {
			if c <= i && dp[i-c]+1 < dp[i] {
				dp[i] = dp[i-c] + 1
			}
		}
	}
	if dp[amount] == unreachable {
		return -1
	}
	return dp[amount]
}
