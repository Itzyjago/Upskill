package algo

import (
	"reflect"
	"testing"
)

func TestGroupAnagrams(t *testing.T) {
	words := []string{"eat", "tea", "tan", "ate", "nat", "bat"}
	want := [][]string{
		{"eat", "tea", "ate"},
		{"tan", "nat"},
		{"bat"},
	}
	if got := GroupAnagrams(words); !reflect.DeepEqual(got, want) {
		t.Errorf("GroupAnagrams(%v) = %v, want %v", words, got, want)
	}
}

func TestGroupAnagramsEmpty(t *testing.T) {
	if got := GroupAnagrams(nil); len(got) != 0 {
		t.Errorf("GroupAnagrams(nil) = %v, want empty", got)
	}
}

func TestGroupAnagramsEveryWordInExactlyOneGroup(t *testing.T) {
	words := []string{"abc", "cab", "xyz", "bca", "zyx"}
	groups := GroupAnagrams(words)
	seen := make(map[string]bool)
	for _, g := range groups {
		for _, w := range g {
			if seen[w] {
				t.Errorf("word %q appeared in more than one group", w)
			}
			seen[w] = true
		}
	}
	for _, w := range words {
		if !seen[w] {
			t.Errorf("word %q missing from output entirely", w)
		}
	}
}
