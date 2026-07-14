package algo

// Knapsack01 solves the 0/1 knapsack problem: given item weights/values and
// a capacity, return the max total value achievable without exceeding
// capacity, taking each item at most once. Classic 2D DP, O(n*capacity)
// time and space.
func Knapsack01(weights, values []int, capacity int) int {
	n := len(weights)
	// dp[i][c] = max value using the first i items with capacity c.
	dp := make([][]int, n+1)
	for i := range dp {
		dp[i] = make([]int, capacity+1)
	}
	for i := 1; i <= n; i++ {
		for c := 0; c <= capacity; c++ {
			dp[i][c] = dp[i-1][c] // option: skip item i-1
			if w := weights[i-1]; w <= c {
				if take := dp[i-1][c-w] + values[i-1]; take > dp[i][c] {
					dp[i][c] = take
				}
			}
		}
	}
	return dp[n][capacity]
}
