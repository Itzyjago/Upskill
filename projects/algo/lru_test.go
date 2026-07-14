package algo

import "testing"

func TestLRUCacheEvictsLeastRecentlyUsed(t *testing.T) {
	c := NewLRUCache(2)
	c.Put(1, 100)
	c.Put(2, 200)
	if v, ok := c.Get(1); !ok || v != 100 {
		t.Fatalf("Get(1) = %d, %v, want 100, true", v, ok)
	}
	// 1 was just touched, so 2 is now least-recently-used and gets evicted.
	c.Put(3, 300)
	if _, ok := c.Get(2); ok {
		t.Fatal("Get(2) should have been evicted")
	}
	if v, ok := c.Get(1); !ok || v != 100 {
		t.Fatalf("Get(1) = %d, %v, want 100, true", v, ok)
	}
	if v, ok := c.Get(3); !ok || v != 300 {
		t.Fatalf("Get(3) = %d, %v, want 300, true", v, ok)
	}
}

func TestLRUCachePutOverwritesAndRefreshes(t *testing.T) {
	c := NewLRUCache(2)
	c.Put(1, 1)
	c.Put(2, 2)
	c.Put(1, 111) // overwrite + refresh recency for key 1
	c.Put(3, 3)   // should evict key 2, not key 1
	if _, ok := c.Get(2); ok {
		t.Fatal("Get(2) should have been evicted")
	}
	if v, ok := c.Get(1); !ok || v != 111 {
		t.Fatalf("Get(1) = %d, %v, want 111, true", v, ok)
	}
}

func TestLRUCacheMissOnUnknownKey(t *testing.T) {
	c := NewLRUCache(2)
	if _, ok := c.Get(42); ok {
		t.Fatal("Get on empty cache should miss")
	}
}
