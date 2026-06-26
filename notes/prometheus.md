# Prometheus & the metrics model

Wrote these while adding a `/metrics` endpoint to wordcount. Prometheus *pulls*
(scrapes) a plaintext endpoint on an interval — the app just exposes numbers, it
never pushes. That inversion is the whole mental model.

## The four metric types
- **Counter** — monotonically increasing; only goes up (or resets to 0 on
  restart). Requests served, errors, bytes. You almost never read a counter
  raw — you `rate()` it. Name them `_total`.
- **Gauge** — a value that goes up *and* down. In-flight requests, queue depth,
  temperature, memory in use.
- **Histogram** — samples observations into cumulative **buckets** plus a
  `_sum` and `_count`. Lets you compute quantiles *server-side* with
  `histogram_quantile`. Latency and sizes live here.
- **Summary** — like a histogram but quantiles are computed **client-side** and
  can't be aggregated across instances. Prefer histograms unless you have a
  reason.

## The text exposition format
Just lines of `name{label="v",...} value`, with `# HELP` and `# TYPE` headers:
```
# HELP http_requests_total Total HTTP requests handled.
# TYPE http_requests_total counter
http_requests_total{method="POST",path="/count",status="200"} 42

# TYPE http_request_duration_seconds histogram
http_request_duration_seconds_bucket{path="/count",le="0.01"} 40
http_request_duration_seconds_bucket{path="/count",le="+Inf"} 42
http_request_duration_seconds_sum{path="/count"} 0.137
http_request_duration_seconds_count{path="/count"} 42
```
- Histogram buckets are **cumulative**: `le="0.01"` counts every observation
  ≤ 10ms, so it includes the ≤1ms ones too. The `+Inf` bucket equals `_count`.
- Content type is `text/plain; version=0.0.4`. That's it — I hand-rolled the
  whole thing with `fmt` instead of pulling in `client_golang`, to actually
  understand it.

## Labels = dimensions (use with care)
- Each distinct label-set is its own time series. **High-cardinality labels**
  (user IDs, request paths with IDs in them) explode series count and melt the
  TSDB. Keep labels bounded: method, route *template*, status class.

## Scraping
- A `prometheus.yml` `scrape_config` lists targets (or service-discovers them).
- In Kubernetes you annotate pods/services with `prometheus.io/scrape: "true"`
  and a `prometheus.io/port`, and Prometheus' k8s SD finds them.

## What clicked
- The app stays dumb: it exposes a snapshot, the server does the time-series
  math. A counter that "resets on restart" is fine *because* `rate()` is built
  to ignore the negative jump.

See [golden-signals.md](golden-signals.md) for *what* to measure and
[promql.md](promql.md) for querying it.
