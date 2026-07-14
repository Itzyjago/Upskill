package algo

import (
	"math"
	"testing"
)

func TestIntSqrt(t *testing.T) {
	cases := []struct{ n, want int }{
		{0, 0}, {1, 1}, {4, 2}, {8, 2}, {9, 3}, {15, 3}, {16, 4}, {2147395600, 46340},
	}
	for _, c := range cases {
		if got := IntSqrt(c.n); got != c.want {
			t.Errorf("IntSqrt(%d) = %d, want %d", c.n, got, c.want)
		}
	}
}

func TestIntSqrtAgreesWithMathSqrtForSmallN(t *testing.T) {
	for n := 0; n < 1000; n++ {
		want := int(math.Sqrt(float64(n)))
		if got := IntSqrt(n); got != want {
			t.Errorf("IntSqrt(%d) = %d, want %d (math.Sqrt reference)", n, got, want)
		}
	}
}

func TestIntSqrtNegativePanics(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error("IntSqrt(-1) should panic")
		}
	}()
	IntSqrt(-1)
}
