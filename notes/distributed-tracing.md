# Distributed tracing & OpenTelemetry

The third pillar (see [structured-logging.md](structured-logging.md) for the
other two). Metrics tell you *something* is slow; a trace tells you **where** —
which hop, in which service, ate the latency. The roadmap's "wire traces" goal;
these are the notes before the code.

## The model: traces and spans
- A **trace** is one request's whole journey, identified by a **trace ID**.
- A **span** is one unit of work in it (an HTTP handler, a DB query), with a
  start/end time, a name, and key/value **attributes**.
- Spans form a **tree** via parent span IDs: the root span is the inbound
  request; children are the calls it makes. The gaps and overlaps in the
  waterfall are where the time goes.

## Context propagation — the crux
- For a child service's spans to join the *same* trace, the trace ID + parent
  span ID must travel **with the request**, usually as HTTP headers.
- **W3C Trace Context** standardized this as the `traceparent` header:
  `version-traceid-spanid-flags`. Everyone speaks it now, so cross-vendor traces
  stitch together.
- Same idea as Go's `context.Context` carrying a deadline down the call tree
  (see [go-context.md](go-context.md)) — except it crosses the network. The span
  context rides *inside* the Go context locally, then gets injected into headers
  on the way out and extracted on the way in.

## OpenTelemetry (OTel) — the vendor-neutral standard
- One set of APIs/SDKs for traces, metrics, and logs; export to Jaeger, Tempo,
  Honeycomb, etc. without rewriting instrumentation.
- **Instrument once, export anywhere.** The **Collector** is a separate process
  that receives, batches, and ships telemetry — so apps don't hard-code a
  backend.
- **Sampling** keeps the cost sane: trace a *fraction* of requests (head
  sampling at the start, or tail sampling that keeps the slow/errored ones).
  100% tracing of high traffic is usually too much data.

## How it'd land in wordcount
- Wrap the handler in OTel's `otelhttp` middleware → a span per request,
  `traceparent` auto-propagated.
- Stamp the `trace_id` into the slog line so a log and its trace cross-link.
- Even single-service, a span tree shows handler vs. body-read vs. encode time.

## Server spans vs. client spans
- Everything above (and `trace.go`) covers the **server** span: the inbound
  request I *receive*. A single-service trace never needs anything else.
- The moment wordcount *calls out* to another service, the outbound hop needs
  its own span, with `kind: CLIENT` instead of `SERVER` (OTLP's span `kind`
  field, `notes/otlp.md`) — same trace-id, a fresh child span-id, one level
  deeper in the tree.
- Propagation on the client side is the mirror image of the server side:
  - **Server**: *extract* the inbound `traceparent`, mint a child, that child
    becomes the server span.
  - **Client**: mint a child of *my* current span, *inject* that child's
    `traceparent` into the outbound request's headers, before sending it.
  - The downstream service's server span then extracts what I just injected —
    my client span's id becomes its parent's id. That's the whole stitch: two
    services, two spans, one trace, connected by one HTTP header.
- Concretely: the client span's `parentID` is the *inbound* server span's id
  (my own request), and the id I inject downstream is the *client* span's id
  — not the server span's. Get that swapped and the downstream service parents
  itself to the wrong span, and the waterfall draws a sibling instead of a
  child.

### The actual shape (wordcount's real two-service trace)
Everything above is the rule; this is what it draws once both wordcount
instances are running. Three spans, same trace-id, each nested inside the one
that started it:
```
trace-id 4bf92f3577b34da6a3ce929d0e0e4736
edge   │ SERVER "POST /count"            0ms ═══════════════════════════ 42ms
       │            │
       │            └ CLIENT "POST /count (upstream)"
       │                 3ms ═══════════════════════════════════ 38ms
       │                        │
upstrm │                        └ SERVER "POST /count"
       │                             5ms ═══════════════════ 35ms
```
- Three spans, but only **two distinct names** — the edge's inbound span and
  the upstream's inbound span are both literally `"POST /count"`
  (`middleware.go` builds the name from `r.Method+" "+path`, and both
  services register `/count` under the same route). In a real waterfall UI
  this reads fine because indentation (parent/child depth) disambiguates
  them, but grepping Jaeger for a span *by name* returns both — the depth is
  what tells you which service it's from, not the name.
  `"POST /count (upstream)"` (`client.go`, a literal string, not derived from
  a route) is the only span name in the trace that's actually unique.
- The nesting *is* the parent-chain from the section above, drawn: the edge's
  CLIENT span's `parentID` is the edge SERVER span's id (same-process, no
  header involved — `withSpan`/`spanFrom` on the Go `context.Context`); the
  upstream SERVER span's `parentID` is the edge CLIENT span's id (crosses the
  network via the injected `traceparent` header). One hop is in-process
  context, the other is the wire — same tree, two different propagation
  mechanisms stitching it together.
- The CLIENT span starts a few ms after its parent (encode + dispatch time)
  and ends a few ms before it (the parent still has response-encoding left to
  do) — a child span's interval nested *inside* its parent's is what a
  correct waterfall always looks like; a child span poking outside its
  parent's bounds is a clock-skew or bug tell.

## What clicked
- Propagation is the whole game. Without passing trace context across the call,
  you get disconnected per-service spans, not one trace. The headers are how a
  request stays one story across a dozen services.
