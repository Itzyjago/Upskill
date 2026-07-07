package main

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"sync"
	"time"
)

// idempotencyTTL is how long a cached response stays eligible for replay —
// long enough to cover a realistic client retry window, short enough that the
// store doesn't grow unbounded between the opportunistic sweeps in store().
const idempotencyTTL = 5 * time.Minute

// errIdempotencyKeyReused means a client reused an Idempotency-Key with a
// different request body. The store can't tell whether that's a genuine new
// request or a client bug, so it refuses instead of guessing which body to
// trust — same reasoning as the Stripe-style API design this mirrors.
var errIdempotencyKeyReused = errors.New("idempotency key reused with a different request body")

type idempotencyEntry struct {
	bodyHash string
	status   int
	body     []byte
	expires  time.Time
}

// idempotencyStore caches a (status, body) pair per Idempotency-Key so a
// retried /count replays the first response instead of counting twice. This
// is the design notes/http.md's #21 entry flagged as needed before
// upstreamClient grows retry logic (roadmap #21) — applied here at the
// client-facing hop first; wiring it into the edge->upstream forward is the
// next step once upstreamClient actually retries.
type idempotencyStore struct {
	mu      sync.Mutex
	ttl     time.Duration
	entries map[string]idempotencyEntry
}

func newIdempotencyStore() *idempotencyStore {
	return newIdempotencyStoreTTL(idempotencyTTL)
}

// newIdempotencyStoreTTL takes an explicit TTL so tests can exercise
// expiry without sleeping for the real 5-minute window.
func newIdempotencyStoreTTL(ttl time.Duration) *idempotencyStore {
	return &idempotencyStore{ttl: ttl, entries: make(map[string]idempotencyEntry)}
}

// hashBody fingerprints a request body so a key reused with a different body
// is detected instead of silently replaying the wrong response.
func hashBody(body []byte) string {
	sum := sha256.Sum256(body)
	return hex.EncodeToString(sum[:])
}

// lookup returns a cached (status, body) for key if present and unexpired.
// ok is false on a miss (new key, or expired); err is set only when key
// exists but bodyHash doesn't match — a reused key, different body.
func (s *idempotencyStore) lookup(key, bodyHash string) (status int, body []byte, ok bool, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	e, found := s.entries[key]
	if !found || time.Now().After(e.expires) {
		return 0, nil, false, nil
	}
	if e.bodyHash != bodyHash {
		return 0, nil, false, errIdempotencyKeyReused
	}
	return e.status, e.body, true, nil
}

// store records a response for key. Sweeps expired entries first — cheap
// enough at this scale that a separate background goroutine would just be
// more moving parts for no real benefit.
func (s *idempotencyStore) store(key, bodyHash string, status int, body []byte) {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	for k, e := range s.entries {
		if now.After(e.expires) {
			delete(s.entries, k)
		}
	}
	s.entries[key] = idempotencyEntry{
		bodyHash: bodyHash,
		status:   status,
		body:     append([]byte(nil), body...),
		expires:  now.Add(s.ttl),
	}
}
