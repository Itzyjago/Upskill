# W3C Trace Context — the `traceparent` wire format

[distributed-tracing.md](distributed-tracing.md) is the concepts (traces, spans,
propagation). This is the actual header I implemented in wordcount — hand-rolled,
same spirit as the hand-rolled Prometheus exposition: build the wire format to
understand it, no OpenTelemetry SDK.

## The header
```
traceparent: 00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01
             └┬┘ └──────────────┬───────────────┘ └──────┬───────┘ └┬┘
           version          trace-id              parent-id     flags
```
- **version** — `00`. Two hex chars. Reject anything you don't understand.
- **trace-id** — 16 bytes / 32 hex. Identifies the *whole* trace. Constant for
  every span across every service. Must not be all-zero.
- **parent-id** — 8 bytes / 16 hex. The span-id of the *sender's* span — i.e.
  the parent of the span you're about to start. Must not be all-zero.
- **flags** — 8 bits / 2 hex. Bit 0 is `sampled` (`01` = record this trace).

All hex is **lowercase**; case matters for validation.

## Extract → child → inject
The receiving side does three things:
1. **Extract**: parse the inbound `traceparent`. Keep the **trace-id** and the
   sender's span-id (that becomes *my* parent).
2. **Start a child span**: same trace-id, a **fresh** span-id (mine). Carry the
   sampled flag through.
3. **Inject**: when calling *downstream*, write a new `traceparent` with the
   trace-id and *my* span-id — so the next hop's parent is me.

Miss any of this and you get disconnected per-service spans instead of one trace.
The trace-id is the thread; a new span-id per hop is what builds the tree.

## Malformed → start fresh, don't fail
If the header is missing or malformed (wrong field count, bad length, non-hex,
all-zero ids, unknown version), you **start a brand-new root trace** — generate a
random trace-id + span-id. A bad upstream header must never 500 the request;
tracing is best-effort telemetry, not part of the contract.

## Generating IDs
- Random bytes from a CSPRNG (`crypto/rand`), hex-encoded. 16 bytes for the
  trace-id, 8 for the span-id. Collisions are astronomically unlikely at these
  widths, so no coordination needed across services.

## Cross-linking logs and traces
- Stamp the `trace_id` (and `span_id`) onto the structured log line (see
  [structured-logging.md](structured-logging.md)). Now a log and its trace point
  at each other: jump from "this request errored" in the logs to the waterfall
  that shows *where* the time went.

## What clicked
- The whole standard is one header carrying *trace-id + parent span-id + flags*,
  and the only real logic is **keep the trace-id, mint a new span-id** at each
  hop. That single rule is what turns a pile of per-service spans into one
  connected story — and it's why everyone agreeing on `traceparent` is what lets
  cross-vendor traces stitch together at all.
