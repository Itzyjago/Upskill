package algo

import "testing"

func TestLFUCacheEvictsLeastFrequentlyUsed(t *testing.T) {
	c := NewLFUCache(2)
	c.Put(1, 1)
	c.Put(2, 2)
	c.Get(1)    // freq(1)=2, freq(2)=1
	c.Put(3, 3) // capacity full, freq(2)=1 is the minimum -> evict key 2

	if _, ok := c.Get(2); ok {
		t.Fatal("key 2 should have been evicted (lowest frequency)")
	}
	if v, ok := c.Get(1); !ok || v != 1 {
		t.Fatalf("Get(1) = %d, %v, want 1, true", v, ok)
	}
	if v, ok := c.Get(3); !ok || v != 3 {
		t.Fatalf("Get(3) = %d, %v, want 3, true", v, ok)
	}
}

func TestLFUCacheTiebreaksOnLeastRecentlyUsed(t *testing.T) {
	c := NewLFUCache(2)
	c.Put(1, 1)
	c.Put(2, 2) // both at freq 1; 1 is the older of the two
	c.Put(3, 3) // tie at freq 1 -> evict the least-recently-touched, key 1

	if _, ok := c.Get(1); ok {
		t.Fatal("key 1 should have been evicted (tie broken by recency)")
	}
	if _, ok := c.Get(2); !ok {
		t.Fatal("key 2 should still be present")
	}
}

func TestLFUCachePutOverwriteBumpsFrequency(t *testing.T) {
	c := NewLFUCache(1)
	c.Put(1, 100)
	c.Put(1, 200) // overwrite, same key, should not evict itself
	if v, ok := c.Get(1); !ok || v != 200 {
		t.Fatalf("Get(1) = %d, %v, want 200, true", v, ok)
	}
}

func TestLFUCacheZeroCapacity(t *testing.T) {
	c := NewLFUCache(0)
	c.Put(1, 1)
	if _, ok := c.Get(1); ok {
		t.Fatal("zero-capacity cache should never retain anything")
	}
}
