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

## Coverage & discipline
- Coverage shows what's *unexecuted*, not what's *correct*. 100% ≠ bug-free.
- Write the failing test first (red → green → refactor) — see TDD.
- Fix flaky tests immediately; a quarantined flake erodes trust in the suite.
