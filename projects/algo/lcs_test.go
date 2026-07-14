package algo

import "testing"

func TestLongestCommonSubsequence(t *testing.T) {
	cases := []struct {
		a, b string
		want int
	}{
		{"abcde", "ace", 3},
		{"abc", "abc", 3},
		{"abc", "def", 0},
		{"", "abc", 0},
		{"aggtab", "gxtxayb", 4},
	}
	for _, c := range cases {
		if got := LongestCommonSubsequence(c.a, c.b); got != c.want {
			t.Errorf("LongestCommonSubsequence(%q, %q) = %d, want %d", c.a, c.b, got, c.want)
		}
	}
}

func TestEditDistance(t *testing.T) {
	cases := []struct {
		a, b string
		want int
	}{
		{"horse", "ros", 3},
		{"intention", "execution", 5},
		{"", "abc", 3},
		{"same", "same", 0},
	}
	for _, c := range cases {
		if got := EditDistance(c.a, c.b); got != c.want {
			t.Errorf("EditDistance(%q, %q) = %d, want %d", c.a, c.b, got, c.want)
		}
	}
}

func TestEditDistanceIsSymmetric(t *testing.T) {
	a, b := "kitten", "sitting"
	if EditDistance(a, b) != EditDistance(b, a) {
		t.Errorf("EditDistance should be symmetric: got %d and %d", EditDistance(a, b), EditDistance(b, a))
	}
}
