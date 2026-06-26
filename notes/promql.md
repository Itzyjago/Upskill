# PromQL — querying the metrics

Once wordcount exposes `/metrics` and Prometheus scrapes it, PromQL is how you
turn raw counters into answers. The mental shift: you almost never read a metric
directly — you compute a **rate** or **quantile** over a time window.

## The data model
- Everything is a labeled **time series**: `metric{label="v"} value @timestamp`.
- Four expression types: **instant vector** (one sample per series, now),
  **range vector** (a window of samples, written `metric[5m]`), **scalar**, and
  **string**.

## rate() — the workhorse for counters
```promql
rate(http_requests_total[5m])
```
- Per-second average increase over 5m. Handles counter **resets** (restarts)
  automatically — that's why counters can reset to 0 safely.
- `rate` = average over the window (smooth); `irate` = last two samples (spiky,
  for fast-moving graphs). Reach for `rate` by default.
- Feed `rate` a range vector (`[5m]`), never a bare counter.

## Aggregation — collapse labels with by/without
```promql
sum(rate(http_requests_total[5m])) by (path)        # total req/s per route
sum(rate(http_requests_total{status=~"5.."}[5m]))   # error req/s, all routes
```
- `by (labels)` keeps only those labels; `without (labels)` drops them.
- `=~` is a regex match — `status=~"5.."` is every 5xx.

## Quantiles from a histogram
```promql
histogram_quantile(0.95,
  sum(rate(http_request_duration_seconds_bucket[5m])) by (le, path))
```
- p95 latency **per route**. You *must* aggregate `by (le)` (and keep any label
  you want to split on) — the function reads the cumulative `le` buckets.
- This is the server-side-quantile payoff from [prometheus.md](prometheus.md):
  the app only ships buckets; Prometheus does the percentile math, and it
  aggregates correctly across instances (summaries can't).

## An error-rate ratio (the kind of thing an SLO watches)
```promql
sum(rate(http_requests_total{status=~"5.."}[5m]))
  /
sum(rate(http_requests_total[5m]))
```
- Errors as a fraction of all traffic. Division matches series by label, so
  aggregate both sides to bare scalars (or matching labels) first.

## What clicked
- `rate()[window]` is the unit of thought, not the raw counter. And
  `histogram_quantile` is why histograms beat summaries: percentiles stay
  meaningful after you `sum` across pods. Maps straight onto the golden signals
  in [golden-signals.md](golden-signals.md).
