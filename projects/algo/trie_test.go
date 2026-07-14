package algo

import "testing"

func TestTrieSearchAndStartsWith(t *testing.T) {
	tr := NewTrie()
	for _, w := range []string{"cat", "car", "cart", "dog"} {
		tr.Insert(w)
	}
	for _, w := range []string{"cat", "car", "cart", "dog"} {
		if !tr.Search(w) {
			t.Errorf("Search(%q) = false, want true", w)
		}
	}
	for _, w := range []string{"ca", "do", "carts", "cats"} {
		if tr.Search(w) {
			t.Errorf("Search(%q) = true, want false (not an inserted word)", w)
		}
	}
	for _, p := range []string{"c", "ca", "car", "d"} {
		if !tr.StartsWith(p) {
			t.Errorf("StartsWith(%q) = false, want true", p)
		}
	}
	if tr.StartsWith("z") {
		t.Error(`StartsWith("z") = true, want false`)
	}
}

func TestTrieEmpty(t *testing.T) {
	tr := NewTrie()
	if tr.Search("anything") {
		t.Error("Search on empty trie should be false")
	}
	if tr.StartsWith("a") {
		t.Error("StartsWith on empty trie should be false")
	}
}
