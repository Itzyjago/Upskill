package algo

import "testing"

func TestValidParentheses(t *testing.T) {
	cases := []struct {
		s    string
		want bool
	}{
		{"()", true},
		{"()[]{}", true},
		{"(]", false},
		{"([)]", false},
		{"{[]}", true},
		{"", true},
		{"(", false},
		{")", false},
		{"(a+b)*[c-d]", true}, // non-bracket chars are just ignored
	}
	for _, c := range cases {
		if got := ValidParentheses(c.s); got != c.want {
			t.Errorf("ValidParentheses(%q) = %v, want %v", c.s, got, c.want)
		}
	}
}
