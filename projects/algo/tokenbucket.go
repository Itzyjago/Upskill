package algo

import "time"

// TokenBucket is a rate limiter: capacity tokens refill at refillRate
// tokens/sec, each Allow() call consumes one. now is injectable so tests
// can control time directly instead of sleeping and racing a real clock —
// the same boundary-fake approach wordcount's httptest doubles use.
type TokenBucket struct {
	capacity   float64
	tokens     float64
	refillRate float64 // tokens per second
	last       time.Time
	now        func() time.Time
}

func NewTokenBucket(capacity, refillRate float64, now func() time.Time) *TokenBucket {
	return &TokenBucket{
		capacity:   capacity,
		tokens:     capacity,
		refillRate: refillRate,
		last:       now(),
		now:        now,
	}
}

// Allow reports whether a request may proceed right now, consuming one
// token if so. Refills lazily on each call rather than on a background
// ticker — no goroutine, no ticker to leak or shut down.
func (b *TokenBucket) Allow() bool {
	current := b.now()
	elapsed := current.Sub(b.last).Seconds()
	b.last = current

	b.tokens += elapsed * b.refillRate
	if b.tokens > b.capacity {
		b.tokens = b.capacity
	}
	if b.tokens < 1 {
		return false
	}
	b.tokens--
	return true
}
