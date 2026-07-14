package algo

import "testing"

func TestHashMapPutGetDelete(t *testing.T) {
	m := NewHashMap()
	m.Put("a", 1)
	m.Put("b", 2)
	if v, ok := m.Get("a"); !ok || v != 1 {
		t.Fatalf("Get(a) = %d, %v, want 1, true", v, ok)
	}
	m.Put("a", 100) // overwrite
	if v, ok := m.Get("a"); !ok || v != 100 {
		t.Fatalf("Get(a) after overwrite = %d, %v, want 100, true", v, ok)
	}
	if !m.Delete("b") {
		t.Fatal("Delete(b) = false, want true")
	}
	if _, ok := m.Get("b"); ok {
		t.Fatal("Get(b) after delete should miss")
	}
	if m.Delete("nonexistent") {
		t.Fatal("Delete on missing key should return false")
	}
}

func TestHashMapGrowsAndKeepsAllEntries(t *testing.T) {
	m := NewHashMap()
	const n = 500 // forces several grow() calls past the default 8 buckets
	for i := 0; i < n; i++ {
		m.Put(string(rune(i))+"-key", i)
	}
	if got := m.Len(); got != n {
		t.Fatalf("Len() = %d, want %d", got, n)
	}
	for i := 0; i < n; i++ {
		key := string(rune(i)) + "-key"
		v, ok := m.Get(key)
		if !ok || v != i {
			t.Fatalf("Get(%q) = %d, %v, want %d, true", key, v, ok, i)
		}
	}
}

func TestHashMapLenTracksPutsAndDeletes(t *testing.T) {
	m := NewHashMap()
	m.Put("x", 1)
	m.Put("y", 2)
	m.Put("x", 999) // overwrite shouldn't increase Len
	if got := m.Len(); got != 2 {
		t.Fatalf("Len() = %d, want 2", got)
	}
	m.Delete("x")
	if got := m.Len(); got != 1 {
		t.Fatalf("Len() after delete = %d, want 1", got)
	}
}
