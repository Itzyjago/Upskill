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
- `POST /count` — counts the request body, returns a JSON tally. The body is
  capped at 10MB (`http.MaxBytesReader`) — an unbounded read is a
  resource-exhaustion DoS, not just a style nit (roadmap #16,
  `notes/security.md`). Over the cap → `413`.
- `GET /metrics` — Prometheus text exposition (see below).
- Every response carries a `traceparent` header; set `OTEL_EXPORTER_OTLP_ENDPOINT`
  to ship the span to a collector (see Tracing below).
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
annotates the pod for Prometheus scraping. It also defines a
HorizontalPodAutoscaler (2-5 replicas, 70% CPU) — roadmap #15. On kind that
needs metrics-server, which isn't installed by default:
```sh
make kind-up kind-metrics-server kind-deploy
kubectl get hpa wordcount   # TARGETS shows a real percentage once metrics-server is up,
                             # not <unknown>
```

## Observability stack (Prometheus + Grafana)
`deploy/observability/` is a docker-compose stack that actually scrapes the
`/metrics` above and graphs the golden signals — roadmap #8.
```sh
make obs-up          # build the app + start Prometheus and Grafana
# drive some traffic so there's something to graph:
for i in $(seq 200); do curl -s --data-binary "hello world" localhost:8080/count >/dev/null; done
open http://localhost:3000     # Grafana — the "Golden Signals" dashboard auto-loads
open http://localhost:9090     # Prometheus — Status > Targets shows wordcount UP
make obs-down        # tear it all down
```
- Prometheus config + alert rules: `deploy/observability/prometheus.yml`,
  `deploy/observability/rules/`.
- Grafana is provisioned (datasource + dashboard) from
  `deploy/observability/grafana/` — no manual clicking, it's all in git.
- Grafana Alerting tab also shows `HighErrorRate (Grafana-native)` — a second,
  Grafana-evaluated copy of the same Prometheus alert, provisioned from
  `grafana/provisioning/alerting/`, routed to the same webhook sink. It's there
  to compare Grafana's alert rule / contact point / notification policy shape
  against Prometheus rule / receiver / route (roadmap #14,
  `notes/grafana.md`) — a real setup should pick one, not run both.

## Tracing (OTLP → Jaeger)
The stack also runs Jaeger, and the app ships a span per request to it via a
hand-rolled OTLP/HTTP JSON exporter (`otlp.go`, no OpenTelemetry SDK) — roadmap
#9. The trace ids are the same ones from the `traceparent` header; OTLP just adds
the timings so Jaeger can draw a waterfall.
```sh
make obs-up          # now also starts Jaeger (OTLP ingest + UI)
for i in $(seq 50); do curl -s --data-binary "hello world" localhost:8080/count >/dev/null; done
open http://localhost:16686    # Jaeger — pick service "wordcount", find traces
```
- Export is opt-in: set `OTEL_EXPORTER_OTLP_ENDPOINT` (the compose file points it
  at `http://jaeger:4318`). Unset → spans are still timed and logged, just not shipped.
- It's best-effort and off the hot path: a collector that's down never fails or
  slows a request (`notes/otlp.md`).

### Two-service trace (roadmap #12)
The compose stack runs a second wordcount instance (`wordcount-upstream`). Set
`WORDCOUNT_UPSTREAM_URL` on an instance and its `/count` **forwards** the
request there instead of counting locally — `client.go` wraps the outbound
call in a client span and injects its own id into the traceparent it sends, so
the upstream's server span parents to *that*, not to the original caller. One
request, two services, one trace:
```sh
curl -s --data-binary "hello world" localhost:8080/count   # hits the edge instance
open http://localhost:16686    # Jaeger — pick "wordcount", the trace has two
                                # spans: wordcount (SERVER) -> wordcount-upstream (SERVER),
                                # joined by wordcount's CLIENT span in between
```
- `OTEL_SERVICE_NAME` (also a standard OTel env var) is set per instance in
  compose, so the two show up as distinct services in Jaeger rather than the
  same service name twice.

### Tail sampling (roadmap #13)
Both instances now export to an OTel **Collector** (`otel-collector.yaml`)
instead of Jaeger directly. The collector holds each trace open for
`decision_wait` (10s) so every span — edge *and* upstream — has landed before
it decides, then applies (in order): keep it if any span errored, keep it if
it's slower than 500ms, otherwise keep a flat 10% baseline. See
`notes/otlp.md` "Sampling: head vs. tail" for why the decision has to happen
*after* the trace, not per-request in the app.
```sh
for i in $(seq 100); do curl -s --data-binary "hello world" localhost:8080/count >/dev/null; done
open http://localhost:16686    # ~10 traces show up, not 100 — the baseline policy
                                # dropped the rest, all clean 200s in a few ms
docker compose -f deploy/observability/docker-compose.yml stop wordcount-upstream
curl -s --data-binary "hello world" localhost:8080/count   # 502 from forwardCountHandler
                                # this trace is a guaranteed keeper: the errors policy
                                # doesn't roll the dice like the baseline one does
```

## Test
```sh
go test ./...
```
Every outbound call this codebase makes gets a real fake, not a nil
escape hatch: `client_test.go` stands up an `httptest.Server` for the
upstream hop, and `middleware_test.go`'s `fakeCollector` does the same for
the OTLP export goroutine — pushing each received payload onto a channel so
a test can wait on a fire-and-forget `go func()` deterministically instead
of racing it with a sleep (`notes/testing.md`).

## What I practiced
- `flag` for CLI parsing and the "no flags → do everything" default.
- Reading from stdin *or* file args; streaming with `bufio` instead of slurping.
- Explicit error returns, non-zero exit on failure, errors to stderr.
- Table-driven tests (see `main_test.go`).
- Test doubles for both outbound calls (`client_test.go`, `middleware_test.go`)
  instead of nil-ing the dependency out and only testing the rest of the path.

Then: containerized it (see the `Dockerfile`), added a `-serve` HTTP mode with a
`/healthz` probe, and a `deploy/` k8s manifest — building toward a real deploy.
