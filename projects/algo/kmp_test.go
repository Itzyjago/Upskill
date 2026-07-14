package algo

import (
	"reflect"
	"strings"
	"testing"
)

func TestKMPSearch(t *testing.T) {
	cases := []struct {
		text, pattern string
		want          []int
	}{
		{"ababcababcabc", "abc", []int{2, 7, 10}},
		{"aaaa", "aa", []int{0, 1, 2}}, // overlapping matches
		{"hello", "xyz", nil},
		{"abc", "abc", []int{0}},
		{"abc", "", nil},
	}
	for _, c := range cases {
		if got := KMPSearch(c.text, c.pattern); !reflect.DeepEqual(got, c.want) {
			t.Errorf("KMPSearch(%q, %q) = %v, want %v", c.text, c.pattern, got, c.want)
		}
	}
}

func TestKMPSearchAgreesWithNaiveScan(t *testing.T) {
	text := "the quick brown fox jumps over the lazy dog the fox"
	pattern := "the"
	want := naiveIndices(text, pattern)
	if got := KMPSearch(text, pattern); !reflect.DeepEqual(got, want) {
		t.Errorf("KMPSearch(%q, %q) = %v, want %v (naive)", text, pattern, got, want)
	}
}

func naiveIndices(text, pattern string) []int {
	var out []int
	for i := 0; i+len(pattern) <= len(text); i++ {
		if strings.HasPrefix(text[i:], pattern) {
			out = append(out, i)
		}
	}
	return out
}
