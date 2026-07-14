package algo

import "testing"

func isPalindrome(s string) bool {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		if s[i] != s[j] {
			return false
		}
	}
	return true
}

func TestLongestPalindromicSubstring(t *testing.T) {
	cases := []struct {
		in        string
		wantLen   int
		wantOneOf []string
	}{
		{"babad", 3, []string{"bab", "aba"}},
		{"cbbd", 2, []string{"bb"}},
		{"a", 1, []string{"a"}},
		{"", 0, []string{""}},
	}
	for _, c := range cases {
		got := LongestPalindromicSubstring(c.in)
		if len(got) != c.wantLen {
			t.Errorf("LongestPalindromicSubstring(%q) = %q (len %d), want len %d", c.in, got, len(got), c.wantLen)
			continue
		}
		if !isPalindrome(got) {
			t.Errorf("LongestPalindromicSubstring(%q) = %q, not actually a palindrome", c.in, got)
		}
		found := false
		for _, w := range c.wantOneOf {
			if got == w {
				found = true
			}
		}
		if !found {
			t.Errorf("LongestPalindromicSubstring(%q) = %q, want one of %v", c.in, got, c.wantOneOf)
		}
	}
}

func TestLongestPalindromicSubstringWholeString(t *testing.T) {
	if got := LongestPalindromicSubstring("racecar"); got != "racecar" {
		t.Errorf("LongestPalindromicSubstring(\"racecar\") = %q, want \"racecar\"", got)
	}
}
