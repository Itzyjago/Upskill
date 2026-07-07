# Testing notes

## The pyramid
- Many fast **unit** tests, fewer **integration** tests, fewest **e2e** tests.
- Inverted (mostly e2e) = slow, flaky suites. Push logic down so it's unit-testable.

## What makes a good test
- Tests behavior, not implementation — refactors shouldn't break it.
- Arrange / Act / Assert. One logical assertion per test.
- Deterministic: no real clock, network, or randomness. Inject those.

## Table-driven tests (Go)
```go
tests := []struct {
    name string; in int; want int
}{
    {"zero", 0, 0},
    {"positive", 3, 9},
}
for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
        if got := square(tt.in); got != tt.want {
            t.Errorf("square(%d) = %d, want %d", tt.in, got, tt.want)
        }
    })
}
```

## Test doubles
- **Stub** returns canned data. **Mock** also asserts it was called as expected.
- **Fake** is a working lightweight impl (in-memory DB). Prefer fakes over deep
  mock chains — mocks couple tests to call structure.

## The fake-collector-shaped hole (wordcount) — closed
Two different answers to "how do I test code that calls an external
collector?" sit side by side in `projects/wordcount`, and only one of them is
a real double:
- `otlpExporter` (the OTLP-to-Jaeger client) is nil-able everywhere it's used
  — `middleware.go`, `client.go` — and every test just passes `nil` to skip
  export entirely (`newMux(m, nil, nil)`). That's not a double, it's an
  **escape hatch**: it proves the *rest* of the request path works, but it
  asserts nothing about export itself — a payload-shape bug in `otlp.go`
  would sail through every handler test untouched. (`otlp_test.go` does cover
  the payload directly, just via a separate, narrower set of tests — the gap
  is that nothing exercises "middleware calls export and the right thing
  happens.")
- `upstreamClient` (the *other* outbound call, `client.go`) gets tested with
  an actual **fake**: `client_test.go` stands up a real `httptest.NewServer`
  and asserts on what it received (the injected `traceparent`) and returns
  canned responses (a 500, a slow body). That's Testing 101 "prefer fakes over
  deep mocks" — a working lightweight HTTP server beats asserting call
  structure, and it costs nothing extra since `net/http/httptest` is already
  in every Go toolchain.
- The difference: `upstreamClient` talks HTTP request/response, trivial to
  fake with `httptest`. `otlpExporter`'s effect is fire-and-forget from a
  goroutine (`middleware.go`'s `go func() { ... tr.export(...) }()`) — testing
  *that* actually happened means either sleeping and racing the goroutine, or
  giving `metrics`/`otlpExporter` a way to synchronize, which nothing here
  does yet. Nil-as-escape-hatch was the path of least resistance, not a
  deliberate design choice — worth coming back to.
- **Resolution**: turned out the synchronization didn't need to live in
  production code at all — `middleware_test.go`'s `fakeCollector` wraps
  `httptest.NewServer` and pushes each decoded payload onto a buffered
  channel; the test just receives off that channel with a bounded
  `time.After` instead of sleeping. The goroutine in `middleware.go` is
  untouched — the fake observes it from the outside, the same way
  `client_test.go` observes `upstreamClient` from the outside. Escape hatch
  → fake didn't require adding a synchronization *feature*, just testing at
  the boundary (the HTTP call) instead of the internals (the goroutine).

## Coverage & discipline
- Coverage shows what's *unexecuted*, not what's *correct*. 100% ≠ bug-free.
- Write the failing test first (red → green → refactor) — see TDD.
- Fix flaky tests immediately; a quarantined flake erodes trust in the suite.
