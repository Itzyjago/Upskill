package algo

// FibMemo computes the nth Fibonacci number (F(0)=0, F(1)=1) via top-down
// memoized recursion. O(n) time and space, vs. naive recursion's O(2^n) —
// the memo turns each subproblem into work done once instead of
// exponentially many times.
func FibMemo(n int) int {
	memo := make(map[int]int)
	var fib func(int) int
	fib = func(n int) int {
		if n <= 1 {
			return n
		}
		if v, ok := memo[n]; ok {
			return v
		}
		v := fib(n-1) + fib(n-2)
		memo[n] = v
		return v
	}
	return fib(n)
}

// FibTabulation computes the nth Fibonacci number bottom-up, iteratively.
// O(n) time, O(1) space — no call stack, no map, just two rolling values.
func FibTabulation(n int) int {
	if n <= 1 {
		return n
	}
	prev, cur := 0, 1
	for i := 2; i <= n; i++ {
		prev, cur = cur, prev+cur
	}
	return cur
}
