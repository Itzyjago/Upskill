package algo

// Trie is a prefix tree over lowercase a-z strings, used for fast
// prefix/word lookups (e.g. autocomplete) — O(len(word)) per operation
// regardless of how many words are stored.
type Trie struct {
	children [26]*Trie
	isWord   bool
}

func NewTrie() *Trie {
	return &Trie{}
}

func (t *Trie) Insert(word string) {
	node := t
	for _, ch := range word {
		i := ch - 'a'
		if node.children[i] == nil {
			node.children[i] = &Trie{}
		}
		node = node.children[i]
	}
	node.isWord = true
}

// Search reports whether word was inserted exactly.
func (t *Trie) Search(word string) bool {
	node := t.walk(word)
	return node != nil && node.isWord
}

// StartsWith reports whether any inserted word has prefix as a prefix.
func (t *Trie) StartsWith(prefix string) bool {
	return t.walk(prefix) != nil
}

func (t *Trie) walk(s string) *Trie {
	node := t
	for _, ch := range s {
		i := ch - 'a'
		if node.children[i] == nil {
			return nil
		}
		node = node.children[i]
	}
	return node
}
