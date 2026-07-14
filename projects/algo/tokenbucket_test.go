package algo

import (
	"testing"
	"time"
)

// fakeClock lets the test control time directly, no real sleeping.
type fakeClock struct {
	t time.Time
}

func (c *fakeClock) now() time.Time { return c.t }
func (c *fakeClock) advance(d time.Duration) {
	c.t = c.t.Add(d)
}

func TestTokenBucketAllowsBurstUpToCapacity(t *testing.T) {
	clock := &fakeClock{t: time.Unix(0, 0)}
	b := NewTokenBucket(3, 1, clock.now)

	for i := 0; i < 3; i++ {
		if !b.Allow() {
			t.Fatalf("Allow() #%d = false, want true (within initial capacity)", i)
		}
	}
	if b.Allow() {
		t.Fatal("Allow() after draining capacity should be false")
	}
}

func TestTokenBucketRefillsOverTime(t *testing.T) {
	clock := &fakeClock{t: time.Unix(0, 0)}
	b := NewTokenBucket(2, 1, clock.now) // 1 token/sec refill

	b.Allow()
	b.Allow()
	if b.Allow() {
		t.Fatal("bucket should be empty immediately after draining")
	}

	clock.advance(1500 * time.Millisecond) // 1.5 tokens refilled
	if !b.Allow() {
		t.Fatal("Allow() should succeed after enough time passed to refill a token")
	}
	if b.Allow() {
		t.Fatal("only ~1.5 tokens refilled, a second Allow() should fail")
	}
}

func TestTokenBucketRefillCapsAtCapacity(t *testing.T) {
	clock := &fakeClock{t: time.Unix(0, 0)}
	b := NewTokenBucket(2, 100, clock.now) // fast refill

	b.Allow()
	b.Allow()
	clock.advance(10 * time.Second) // would overfill without the capacity cap

	allowed := 0
	for i := 0; i < 5; i++ {
		if b.Allow() {
			allowed++
		}
	}
	if allowed != 2 {
		t.Errorf("allowed = %d after long idle period, want 2 (capped at capacity)", allowed)
	}
}
