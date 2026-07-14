package algo

// GCD returns the greatest common divisor via the Euclidean algorithm.
// O(log(min(a,b))).
func GCD(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	return a
}

// LCM returns the least common multiple, derived from GCD via
// a*b = GCD(a,b)*LCM(a,b).
func LCM(a, b int) int {
	if a == 0 || b == 0 {
		return 0
	}
	return a / GCD(a, b) * b
}

// SieveOfEratosthenes returns every prime <= n. O(n log log n) — each
// composite gets crossed off once per prime factor, not re-tested by
// trial division from scratch.
func SieveOfEratosthenes(n int) []int {
	if n < 2 {
		return nil
	}
	isComposite := make([]bool, n+1)
	var primes []int
	for i := 2; i <= n; i++ {
		if isComposite[i] {
			continue
		}
		primes = append(primes, i)
		for j := i * i; j <= n; j += i {
			isComposite[j] = true
		}
	}
	return primes
}

// ModPow computes (base^exp) mod m using binary exponentiation. O(log exp)
// time, vs. naive repeated multiplication's O(exp) — and it never
// materializes the full base^exp, which would overflow for even modest
// inputs.
func ModPow(base, exp, m int) int {
	if m == 1 {
		return 0
	}
	result := 1
	base %= m
	if base < 0 {
		base += m
	}
	for exp > 0 {
		if exp&1 == 1 {
			result = result * base % m
		}
		exp >>= 1
		base = base * base % m
	}
	return result
}
