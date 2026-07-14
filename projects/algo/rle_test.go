package algo

import "testing"

func TestRunLengthEncode(t *testing.T) {
	cases := []struct{ in, want string }{
		{"aaabbbccd", "a3b3c2d"},
		{"abc", "abc"},
		{"", ""},
		{"aaaaaaaaaaaa", "a12"},
	}
	for _, c := range cases {
		if got := RunLengthEncode(c.in); got != c.want {
			t.Errorf("RunLengthEncode(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestRunLengthEncodeDecodeRoundTrip(t *testing.T) {
	inputs := []string{"aaabbbccd", "abc", "", "aaaaaaaaaaaa", "wwwwaaadexeeee"}
	for _, in := range inputs {
		encoded := RunLengthEncode(in)
		decoded := RunLengthDecode(encoded)
		if decoded != in {
			t.Errorf("round trip failed: %q -> %q -> %q", in, encoded, decoded)
		}
	}
}
