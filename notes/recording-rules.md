# Recording rules — precompute the expensive PromQL

Alerting rules ([alerting.md](alerting.md)) turn a PromQL expression into a
firing alert. **Recording rules** are the other half of `rule_files`: they
evaluate an expression on a timer and save the result as a *new time series*, so
dashboards and alerts read a cheap pre-aggregated metric instead of re-running
the math every refresh. Wrote this adding them to the wordcount stack
(`deploy/observability/rules/recording.yml`).

## The problem they solve
- The golden-signal queries in [promql.md](promql.md) aren't free:
  `histogram_quantile(0.95, sum(rate(...[5m])) by (le))` re-runs `rate()` over
  every bucket series, every time a panel refreshes or an alert evaluates.
- A dashboard with 6 panels open on 5 people's screens + the alert evaluating it
  too = the same heavy expression run dozens of times a minute. Wasteful and
  slow on a real series count.
- A recording rule runs it **once per `evaluation_interval`** and stores the
  answer. Everything downstream reads one flat series.

## The naming convention
Prometheus' own convention, and it's load-bearing — follow it:
```
level:metric:operations
```
- **level** — the aggregation level / labels kept (e.g. `job`, `instance`).
- **metric** — the underlying metric name.
- **operations** — what was applied, newest first (`rate5m`, `sum`, `ratio`).
```yaml
groups:
  - name: wordcount-golden-signals-recording
    interval: 5s
    rules:
      - record: job:http_requests:rate5m
        expr: sum(rate(http_requests_total[5m]))
      - record: job:http_requests:error_ratio5m
        expr: |
          sum(rate(http_requests_total{status=~"5.."}[5m]))
            / sum(rate(http_requests_total[5m]))
      - record: job:http_request_duration_seconds:p95_5m
        expr: |
          histogram_quantile(0.95,
            sum(rate(http_request_duration_seconds_bucket[5m])) by (le))
```
- The colon-delimited name is *reserved for recording rules* — raw metrics use
  `_`, never `:`. Seeing a `:` in a series name tells you instantly it's derived.

## When to use them (and not)
- **Use** for expressions that are expensive *and* queried often: dashboard
  panels, alert expressions, anything aggregated across many series.
- **Don't** pre-record everything — each rule is a new series stored forever and
  costs evaluation time. Record the hot, reused queries, not one-offs.
- **Alerts can read recorded series too**, so the alert and the dashboard share
  one definition of "p95 latency" instead of drifting apart — a real footgun
  when the two copies of the PromQL diverge.

## What clicked
- It's the classic compute trade: **do the work once, ahead of time, and read it
  cheap**, vs. recomputing on every read. Same instinct as a materialized view in
  SQL ([sql-indexing.md](sql-indexing.md)) — denormalize the expensive aggregate
  so the read path is flat. The `:` convention is the tell that you're looking at
  a precomputed answer, not a raw measurement.
