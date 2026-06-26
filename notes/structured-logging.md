# Structured logging (and Go's slog)

Wrote this after switching wordcount's request log to `log/slog`. The shift:
stop logging *sentences*, start logging **key/value events** a machine can
filter and aggregate. `"request done in 5ms"` is for humans; `{"msg":"request",
"path":"/count","dur_ms":5}` is for queries.

## Why structured beats printf
- **Queryable**: `status>=500 AND path="/count"` instead of regexing free text.
- **Stable**: adding a field doesn't break a parser the way reordering a
  sentence does.
- **Aggregatable**: log pipelines (Loki, ELK) index the fields directly.

## Go's slog (stdlib since 1.21)
```go
slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
slog.Info("request", "method", r.Method, "path", path, "status", code, "dur_ms", ms)
```
- A **Handler** decides the output format: `TextHandler` (dev) vs `JSONHandler`
  (prod). Swapping is a one-line change.
- **Levels**: Debug/Info/Warn/Error, with a configurable threshold via
  `HandlerOptions{Level: ...}`.
- `slog.With("request_id", id)` returns a logger that stamps that field on every
  line — the idiom for a **correlation ID** threaded through a request.
- Alternating key, value pairs are easy to typo; `slog.String("k", v)` /
  `slog.Int(...)` attrs are the type-safe form for hot paths.

## Logs to stdout, in containers
- A 12-factor / containerized service logs to **stdout/stderr** and lets the
  platform (Docker, k8s) collect the stream. Don't manage log files in-process.
- One JSON object **per line** (newline-delimited) — that's what collectors
  expect; pretty-printing across lines breaks them.

## Logs vs metrics vs traces — the three pillars
- **Metrics** — cheap aggregate numbers over time. "How many 500s/sec?" (See
  [prometheus.md](prometheus.md).) Constant cost regardless of traffic.
- **Logs** — discrete events with detail. "*Why* did *this* request 500?" Cost
  scales with volume.
- **Traces** — one request's path across services, with timing per hop. "*Where*
  did the 800ms go?" (See [distributed-tracing.md](distributed-tracing.md).)
- They compose: a metric alerts, a trace localizes, a log explains. Share a
  `request_id`/`trace_id` so you can jump between them.

## What clicked
- Same dimensions, two destinations: wordcount records `{method, path, status,
  dur}` into both the metric *and* the log line. The metric answers "how often /
  how slow" in aggregate; the log keeps the individual evidence for when a metric
  spikes and you need the specific offender.
