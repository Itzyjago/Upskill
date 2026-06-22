# Observability notes

## The three pillars
- **Logs** — discrete events with detail. "What happened here?"
- **Metrics** — numbers over time, cheap to store/aggregate. "How much / how
  often / how fast?"
- **Traces** — one request's path across services. "Where did the time go?"

## Logs
- Structured (JSON) beats free text — you can filter and aggregate on fields.
- Attach a `request_id` / `trace_id` to every line so you can stitch a request
  together across services.
- Levels mean something: `ERROR` = needs attention, `INFO` = normal milestones,
  `DEBUG` = off in prod by default.

## Metrics
- Four types: counter (only goes up), gauge (up/down), histogram, summary.
- Watch the **RED** method for services: Rate, Errors, Duration.
- For resources, **USE**: Utilization, Saturation, Errors.
- Alert on symptoms users feel (latency, error rate), not every CPU blip.

## Traces
- A trace is a tree of spans; each span = one operation with start/end + tags.
- Context propagation (W3C `traceparent` header) is what links spans across
  service boundaries.
- OpenTelemetry is the vendor-neutral standard for emitting all three signals.

## Mindset
- Observability ≠ monitoring. Monitoring answers known questions; observability
  lets you ask *new* questions about an unfamiliar failure after the fact.
