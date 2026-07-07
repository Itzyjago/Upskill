package main

import (
	"testing"
	"time"
)

func TestIdempotencyStoreMissThenHit(t *testing.T) {
	s := newIdempotencyStore()
	hash := hashBody([]byte("hello"))

	if _, _, ok, err := s.lookup("key-1", hash); ok || err != nil {
		t.Fatalf("lookup on empty store: ok=%v err=%v, want ok=false err=nil", ok, err)
	}

	s.store("key-1", hash, 200, []byte(`{"lines":1}`))

	status, body, ok, err := s.lookup("key-1", hash)
	if !ok || err != nil {
		t.Fatalf("lookup after store: ok=%v err=%v, want ok=true err=nil", ok, err)
	}
	if status != 200 || string(body) != `{"lines":1}` {
		t.Errorf("lookup = (%d, %q), want (200, %q)", status, body, `{"lines":1}`)
	}
}

func TestIdempotencyStoreKeyReusedDifferentBody(t *testing.T) {
	s := newIdempotencyStore()
	s.store("key-1", hashBody([]byte("hello")), 200, []byte(`{"lines":1}`))

	_, _, ok, err := s.lookup("key-1", hashBody([]byte("goodbye")))
	if ok {
		t.Fatal("lookup with mismatched body hash returned ok=true, want a conflict")
	}
	if err != errIdempotencyKeyReused {
		t.Errorf("err = %v, want errIdempotencyKeyReused", err)
	}
}

func TestIdempotencyStoreExpiry(t *testing.T) {
	s := newIdempotencyStoreTTL(10 * time.Millisecond)
	hash := hashBody([]byte("hello"))
	s.store("key-1", hash, 200, []byte("cached"))

	time.Sleep(20 * time.Millisecond)

	if _, _, ok, err := s.lookup("key-1", hash); ok || err != nil {
		t.Fatalf("lookup after expiry: ok=%v err=%v, want a clean miss", ok, err)
	}
}

func TestIdempotencyStoreSweepsExpiredOnStore(t *testing.T) {
	s := newIdempotencyStoreTTL(10 * time.Millisecond)
	s.store("stale", hashBody([]byte("a")), 200, []byte("a"))
	time.Sleep(20 * time.Millisecond)

	s.store("fresh", hashBody([]byte("b")), 200, []byte("b"))

	s.mu.Lock()
	_, staleStillPresent := s.entries["stale"]
	n := len(s.entries)
	s.mu.Unlock()
	if staleStillPresent {
		t.Error("expired entry survived a later store() call — sweep didn't run")
	}
	if n != 1 {
		t.Errorf("entries after sweep = %d, want 1 (just the fresh one)", n)
	}
}
