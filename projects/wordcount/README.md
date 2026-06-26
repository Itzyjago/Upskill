# wordcount

A tiny `wc` clone in Go — my "build a small CLI to make the concurrency/stdlib
notes stick" roadmap goal, step one.

## Build & run
```sh
go build -o bin/wc .
echo "hello world" | ./bin/wc          # 1  2  12
./bin/wc -w notes.md                    # words only
./bin/wc *.md                           # per-file + total
```

## Serve mode
Runs the same counter as a small HTTP service — built to make the
liveness/readiness-probe notes stick.
```sh
./bin/wc -serve :8080
curl -s localhost:8080/healthz                 # ok   (probe target)
curl -s --data-binary "hello world" \
     localhost:8080/count                      # {"lines":0,"words":2,"bytes":11}
```
- `GET /healthz` — cheap readiness/liveness probe, always 200 when serving.
- `POST /count` — counts the request body, returns a JSON tally.
- `GET /metrics` — Prometheus text exposition (see below).
- `SIGTERM` triggers a graceful shutdown that drains in-flight requests.

## Metrics
Every request through `/healthz` and `/count` is instrumented; `/metrics`
exposes it in Prometheus format — hand-rolled with `fmt`, no `client_golang`
dependency, to learn the exposition format.
```sh
curl -s localhost:8080/metrics
# http_requests_total{method="POST",path="/count",status="200"} 1
# http_request_duration_seconds_bucket{method="POST",path="/count",le="0.01"} 1
# http_request_duration_seconds_count{method="POST",path="/count"} 1
# http_requests_in_flight 0
```
- **counter** `http_requests_total` — by method, route, and status.
- **histogram** `http_request_duration_seconds` — latency, by route.
- **gauge** `http_requests_in_flight` — requests being served right now.

See `deploy/` for the Kubernetes manifest that wires `/healthz` to a probe and
annotates the pod for Prometheus scraping.

## Test
```sh
go test ./...
```

## What I practiced
- `flag` for CLI parsing and the "no flags → do everything" default.
- Reading from stdin *or* file args; streaming with `bufio` instead of slurping.
- Explicit error returns, non-zero exit on failure, errors to stderr.
- Table-driven tests (see `main_test.go`).

Then: containerized it (see the `Dockerfile`), added a `-serve` HTTP mode with a
`/healthz` probe, and a `deploy/` k8s manifest — building toward a real deploy.
