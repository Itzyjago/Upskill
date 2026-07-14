package algo

import (
	"reflect"
	"testing"
)

func TestGCD(t *testing.T) {
	cases := []struct{ a, b, want int }{
		{48, 18, 6},
		{17, 5, 1},
		{0, 5, 5},
		{-12, 8, 4},
	}
	for _, c := range cases {
		if got := GCD(c.a, c.b); got != c.want {
			t.Errorf("GCD(%d, %d) = %d, want %d", c.a, c.b, got, c.want)
		}
	}
}

func TestLCM(t *testing.T) {
	cases := []struct{ a, b, want int }{
		{4, 6, 12},
		{3, 5, 15},
		{0, 7, 0},
	}
	for _, c := range cases {
		if got := LCM(c.a, c.b); got != c.want {
			t.Errorf("LCM(%d, %d) = %d, want %d", c.a, c.b, got, c.want)
		}
	}
}

func TestSieveOfEratosthenes(t *testing.T) {
	want := []int{2, 3, 5, 7, 11, 13, 17, 19}
	if got := SieveOfEratosthenes(20); !reflect.DeepEqual(got, want) {
		t.Errorf("SieveOfEratosthenes(20) = %v, want %v", got, want)
	}
	if got := SieveOfEratosthenes(1); got != nil {
		t.Errorf("SieveOfEratosthenes(1) = %v, want nil", got)
	}
}

func TestSieveAgreesWithTrialDivision(t *testing.T) {
	isPrimeTrial := func(n int) bool {
		if n < 2 {
			return false
		}
		for i := 2; i*i <= n; i++ {
			if n%i == 0 {
				return false
			}
		}
		return true
	}
	primes := SieveOfEratosthenes(100)
	primeSet := make(map[int]bool, len(primes))
	for _, p := range primes {
		primeSet[p] = true
	}
	for n := 0; n <= 100; n++ {
		if primeSet[n] != isPrimeTrial(n) {
			t.Errorf("sieve disagrees with trial division at n=%d: sieve=%v, trial=%v", n, primeSet[n], isPrimeTrial(n))
		}
	}
}

func TestModPow(t *testing.T) {
	cases := []struct{ base, exp, m, want int }{
		{2, 10, 1000, 24}, // 2^10 = 1024, mod 1000 = 24
		{3, 0, 7, 1},      // anything^0 = 1
		{5, 3, 1, 0},      // mod 1 is always 0
		{7, 2, 13, 10},    // 49 mod 13 = 10
	}
	for _, c := range cases {
		if got := ModPow(c.base, c.exp, c.m); got != c.want {
			t.Errorf("ModPow(%d, %d, %d) = %d, want %d", c.base, c.exp, c.m, got, c.want)
		}
	}
}
