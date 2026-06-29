# OTLP — the wire the spans ride out on

[trace-context.md](trace-context.md) gets a `traceparent` across a hop, and
[distributed-tracing.md](distributed-tracing.md) covers the model. But a
propagated header still records **no timings** — there's no waterfall to look at
until a span's start/end actually leaves the process. OTLP is how it leaves.
Wrote this exporting wordcount's spans to Jaeger (`projects/wordcount/otlp.go`).

## What OTLP is
- **OTLP** = OpenTelemetry Protocol: the standard wire format a process uses to
  ship traces/metrics/logs to a collector or backend. Instrument once, point it
  at Jaeger, Tempo, Honeycomb — the *protocol* is the contract, not the vendor.
- Two transports, same payload schema:
  - **gRPC** on `:4317` — the default for SDKs, protobuf over HTTP/2.
  - **HTTP** on `:4318` — protobuf *or* JSON `POST`ed to `/v1/traces`. The JSON
    flavor is plain `application/json`, so you can hand-roll it with `net/http`
    and `encoding/json` — no SDK, same spirit as the hand-rolled `/metrics`.

## The payload shape (JSON)
A trace export is a tree of three nesting levels, outermost first:
```
ResourceSpans      // one per service — carries the Resource (service.name, ...)
  ScopeSpans       // one per instrumentation library/"scope"
    Span[]         // the actual spans
```
```jsonc
{ "resourceSpans": [{
  "resource": { "attributes": [
    { "key": "service.name", "value": { "stringValue": "wordcount" } } ] },
  "scopeSpans": [{
    "scope": { "name": "wordcount/serve" },
    "spans": [{
      "traceId":      "<32 hex>",      // the SAME ids from traceparent
      "spanId":       "<16 hex>",
      "parentSpanId": "<16 hex|omit>", // empty on a root span
      "name":         "POST /count",
      "kind":         2,                // 2 = SERVER
      "startTimeUnixNano": "1719...000",
      "endTimeUnixNano":   "1719...123",
      "attributes": [ /* http.method, http.route, http.status_code */ ],
      "status": { "code": 2 }           // 0 unset, 1 OK, 2 ERROR
    }]
  }]
}] }
```

## Gotchas I hit
- **Times are unix *nanoseconds*, as JSON strings.** They're 64-bit and would
  lose precision as JSON numbers, so OTLP/JSON encodes int64s as strings.
  `strconv.FormatInt(t.UnixNano(), 10)`, not a bare number.
- **IDs stay hex.** OTLP/protobuf wants raw bytes, but OTLP/**JSON** wants the
  lowercase-hex string — which is exactly what `trace.go` already produces. No
  conversion needed for the JSON transport.
- **Export off the hot path.** Don't block the response on a network `POST` to
  the collector. Fire it from a goroutine after the handler returns (a real SDK
  batches; one-span-per-goroutine is the toy version — fine to learn the wire).
- **Best-effort, like propagation.** A collector that's down must never fail a
  request. Swallow the export error (log at debug) and move on.

## What clicked
- The `traceId`/`spanId`/`parentSpanId` in the OTLP payload are *literally* the
  ids `trace.go` already mints for the `traceparent`. The header was always
  carrying the trace's identity across hops; OTLP just adds the **timings** so a
  backend can draw the tree as a waterfall. Propagation and export are two halves
  of the same span — the header moves the id, OTLP reports what happened to it.
