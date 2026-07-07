package main

import (
	"math/rand/v2"
	"strings"
	"testing"
)

func TestCount(t *testing.T) {
	tests := []struct {
		name         string
		in           string
		lines, words int
		bytes        int
	}{
		{"empty", "", 0, 0, 0},
		{"one word no newline", "hello", 0, 1, 5},
		{"one line", "hello world\n", 1, 2, 12},
		{"extra spaces", "  a   b  \n", 1, 2, 10},
		{"two lines", "a\nb\n", 2, 2, 4},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := count(strings.NewReader(tt.in))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if c.Lines != tt.lines || c.Words != tt.words || c.Bytes != tt.bytes {
				t.Errorf("got {l:%d w:%d c:%d}, want {l:%d w:%d c:%d}",
					c.Lines, c.Words, c.Bytes, tt.lines, tt.words, tt.bytes)
			}
		})
	}
}

// randomCountable generates a random mix of letters and the bytes count()
// actually treats specially (space, tab, \n, \r), so generated inputs
// exercise word-boundary logic instead of just being uniform noise.
func randomCountable(rng *rand.Rand, n int) string {
	alphabet := []byte("abcXYZ \t\n\r")
	b := make([]byte, n)
	for i := range b {
		b[i] = alphabet[rng.IntN(len(alphabet))]
	}
	return string(b)
}

// TestCountIsIdempotentAcrossRandomInputs is the property-based test the
// scratch log floated as an alternative to the fixed table above: instead of
// a handful of hand-picked cases, generate N random inputs and check
// invariants that must hold for *any* input, not just the ones someone
// thought to write down.
//
//   - count() is a pure function of its input — notes/http.md's idempotency
//     analysis of /count rests on this claim; here it's actually checked
//     instead of just asserted in prose. Same bytes in, twice, must produce
//     the exact same counts.
//   - Bytes always equals len(input) — a cross-check the table-driven test
//     doesn't exercise directly, since every hand-picked case's byte count
//     was computed by hand rather than derived from the input's actual length.
//
// A fixed seed keeps this deterministic — a property test that fails
// differently on every run would be worse than the table it's supplementing,
// not better, because a red run couldn't be reproduced from the failure
// message alone.
func TestCountIsIdempotentAcrossRandomInputs(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for i := 0; i < 200; i++ {
		in := randomCountable(rng, rng.IntN(200))

		first, err := count(strings.NewReader(in))
		if err != nil {
			t.Fatalf("input %q: unexpected error: %v", in, err)
		}
		if first.Bytes != len(in) {
			t.Fatalf("input %q: Bytes = %d, want len(input) = %d", in, first.Bytes, len(in))
		}

		second, err := count(strings.NewReader(in))
		if err != nil {
			t.Fatalf("input %q: unexpected error on second call: %v", in, err)
		}
		if first != second {
			t.Fatalf("input %q: count() not idempotent — got %+v then %+v", in, first, second)
		}
	}
}
