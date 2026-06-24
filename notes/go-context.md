# Go context notes

The piece that was still fuzzy: how cancellation actually propagates.

## What context is for
- Carries deadlines, cancellation signals, and request-scoped values *across API
  boundaries* and goroutines.
- Rule of thumb: pass `ctx` as the **first** argument, never store it in a struct.

## Cancellation propagates down a tree
```go
ctx, cancel := context.WithCancel(context.Background())
defer cancel() // always call cancel to release resources, even on the happy path

go worker(ctx)
// ...later...
cancel() // every goroutine watching ctx.Done() unblocks
```
- `cancel()` closes the `Done()` channel; children derived from `ctx` are
  cancelled too. Cancellation flows *down* the tree, never up.

## Deadlines and timeouts
```go
ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
defer cancel()

select {
case res := <-doWork(ctx):
    return res, nil
case <-ctx.Done():
    return nil, ctx.Err() // context.DeadlineExceeded or context.Canceled
}
```
- `WithDeadline` is the absolute-time variant of `WithTimeout`.

## Gotchas
- Always `defer cancel()` — leaking a context leaks the goroutine/timer behind it.
- `ctx.Value` is for request-scoped data (request ID, auth), **not** for passing
  optional function params. Keep keys to unexported custom types to avoid clashes.
- A cancelled context stays cancelled — check `ctx.Err()` to tell which.
