# Go notes

## Goroutines and channels
- A goroutine is a cheap, runtime-scheduled "thread": `go doWork()`.
- Channels pass values *and* synchronize: a send blocks until a receive is ready
  (for unbuffered channels).
```go
ch := make(chan int)
go func() { ch <- 42 }()  // blocks until main receives
fmt.Println(<-ch)         // 42
```
- "Don't communicate by sharing memory; share memory by communicating."

## select
```go
select {
case v := <-ch:   fmt.Println(v)
case <-time.After(time.Second): fmt.Println("timeout")
}
```
- `select` waits on multiple channel ops; great for timeouts and cancellation.
- Pair with `context.Context` to cancel work cleanly.

## Error handling
- Errors are values, returned explicitly — no exceptions.
- Wrap with context: `fmt.Errorf("load config: %w", err)`, then `errors.Is` /
  `errors.As` to inspect the chain.

## Gotchas
- A `nil` slice/map is readable but writing to a nil map panics — `make` it first.
- Loop variable capture (pre-1.22): goroutines closing over `i` all saw the last
  value. 1.22+ gives each iteration its own copy.
- `defer` runs LIFO at function return — handy for `Close()` and unlocking.

## Building a CLI
- Stdlib `flag` covers basics; `cobra` for subcommands.
- `go build` produces a single static binary — easy to ship.
