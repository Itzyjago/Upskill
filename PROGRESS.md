# Progress log

A running journal — newest first. One short entry per session.

## 2026-07-02
- Cleared #12, the goal from last time: a real two-service trace. `client.go`
  adds `upstreamClient` — wraps an outbound `/count` call in a **client** span
  (`kind` finally travels through `otlp.go` instead of every span being
  hardcoded SERVER), mints a child of whatever's in `ctx`, and injects *that
  child's* id — not its own parent's — into the outbound `traceparent`. Wired
  in via `WORDCOUNT_UPSTREAM_URL` (env, same pattern as the OTLP endpoint) so
  compose can run an "edge" instance forwarding to an "upstream" one; each
  gets its own `OTEL_SERVICE_NAME` so Jaeger's dropdown shows two real
  services, not one instance twice.
- Also fixed a real bug sitting in the repo: `trace_test.go` was calling
  `newMux` with one argument after a prior commit added the exporter param —
  it's been failing to compile. Caught it while reading the code before
  starting #12; two-line fix.
- What clicked: server spans *extract* a parent and mint a child; client spans
  mint a child and *inject* it. Same operation, opposite direction — the id
  that goes out on the wire is always "my new span," never "the span I got."
  Mixing that up is the one way to make Jaeger draw a sibling instead of a
  child.
- Goal for next time: tail sampling at the collector (#13).

## 2026-06-29
- Extended the observability arc — cleared the three open "next up" goals (#9–11),
  all in the project's hand-rolled, no-SDK spirit.
- **#9 OTLP → Jaeger.** Closed the gap `trace.go` left: it propagated a
  `traceparent` but recorded no timings. Added a `span` type (start/end, name,
  parent id, status) and a hand-rolled OTLP/HTTP **JSON** exporter (`otlp.go`)
  that POSTs `ResourceSpans` to `/v1/traces`. The middleware times each request
  and ships the span off the hot path (best-effort, like propagation). Jaeger
  now runs in the compose stack — a real waterfall on :16686. Tests cover the
  payload (id mapping, root-span parent omitted, ns-as-string).
- **#10 Alertmanager.** Added it to the stack and pointed Prometheus at it;
  `alertmanager.yml` routes `severity=page` to a webhook. Wrote the receiver as a
  `wc -webhook` mode (`webhook.go`) so a fired alert lands in a log line instead
  of just lighting up the UI — reusing the same binary, no new tool.
- **#11 Recording rules.** `rules/recording.yml` precomputes the golden signals
  as `job:...` series; the alerts and the error-ratio panel now read those, so
  the dashboard and the alert can't drift apart.
- What clicked: OTLP export and `traceparent` propagation are two halves of one
  span — the header moves the *id* across a hop, OTLP reports the *timings* to a
  backend. And the noise controls live in two layers: `for:` (Prometheus, "is it
  real?") vs. grouping/`repeat_interval` (Alertmanager, "how often does a human
  hear it?"). Recording rules are just a materialized view for metrics.
- Caveat: built without a local Go/Docker toolchain on this box — code is written
  to compile and the configs to load, but `go test ./...` + `make obs-up` still
  need a run on a machine that has them. That's goal one for next time.
- Goal for next time: a two-service trace (#12) — wordcount calling itself, the
  span tree stitching across the hop in Jaeger.

## 2026-06-27 (cont.)
- Finished the observability arc — all three pillars now real. Stood up a
  **Prometheus + Grafana** stack in `deploy/observability/` (compose): Prometheus
  scrapes wordcount's `/metrics`, a **provisioned** Grafana dashboard graphs the
  golden signals straight from the PromQL in `notes/promql.md`, and alert rules
  load under the Alerts tab. `make obs-up` and it's all there — roadmap #8 done.
- Wired **tracing** (#7) the same way I did metrics: hand-rolled, no OTel SDK.
  `trace.go` parses/validates a W3C `traceparent`, keeps the trace-id and mints a
  fresh child span, and the middleware stamps `trace_id`/`span_id` onto the slog
  line so a log and its trace cross-link. Table tests cover the malformed-header
  cases (bad version, all-zero ids, uppercase hex → start a fresh trace).
- Notes: Grafana, the Prometheus scrape side (pull model, `up`, relabeling),
  alerting (rules vs Alertmanager, the `for:` window), and W3C Trace Context.
- What clicked: provisioning is the same lesson as everything else here — if a
  dashboard or data source isn't a file in git, it's a snowflake that dies on
  restart. And propagation really is the whole game for tracing: *keep the
  trace-id, new span-id per hop* is the one rule that turns scattered spans into
  one connected request.
- Goal for next time: export real spans to Jaeger/Tempo for an actual waterfall
  (#9), and add Alertmanager so a firing rule notifies somewhere real (#10).

## 2026-06-27
- Closed out the deploy goal: brought up a **kind** cluster, loaded the locally
  built image (the step everyone forgets), and applied `deploy/k8s.yaml`. Watched
  `kubectl rollout status` block on the **readiness** probe and saw a bad-probe
  pod stall the rollout instead of cutting over — roadmap #4 and #5 done.
- Gave wordcount **observability**: a `/metrics` endpoint in Prometheus text
  format (counter + latency histogram + in-flight gauge), hand-rolled with `fmt`
  — no `client_golang` — to actually learn the exposition format. Instrumented
  the handlers with middleware (a `statusRecorder` to capture the code) and added
  structured per-request logging via `log/slog` (JSON to stdout).
- Notes: Prometheus, the four golden signals, kind, PromQL, structured logging,
  and distributed tracing. Annotated the pods for Prometheus scraping and added
  `make kind-*` targets for the cluster loop.
- What clicked: histograms ship *cumulative buckets* and let Prometheus compute
  p95 server-side with `histogram_quantile` — and those percentiles still mean
  something after you `sum` across pods, which is exactly why histograms beat
  summaries. Also: metrics, logs, and traces are the same event seen three ways —
  aggregate, detail, and causal path.
- Goal for next time: wire OpenTelemetry traces (#7) and stand up a real
  Prometheus + Grafana to scrape `/metrics` and graph the golden signals (#8).

## 2026-06-25
- Knocked out the "real CI pipeline" goal: `.github/workflows/ci.yml` runs
  lint → test → build, staged with `needs` so a red lint blocks the rest, with
  Go module caching keyed on `go.sum`. Added a Makefile and a golangci-lint
  config so local and CI runs can't drift.
- Gave wordcount a `-serve` mode — an HTTP service with a `/healthz` probe and a
  JSON `/count` endpoint, plus graceful SIGTERM shutdown via
  `signal.NotifyContext`. Wrote a k8s Deployment/Service that wires `/healthz`
  to liveness + readiness probes.
- Notes: GitHub Actions, YAML (the Norway problem bit me — `1.20` parsed as
  1.2), health checks, Go JSON, and semantic versioning.
- What clicked: liveness vs readiness is about the *action* on failure — restart
  vs pull-from-LB. Using a dependency check as liveness turns a blip into a
  restart storm. Also: struct fields must be exported for `encoding/json` to see
  them, which is why `counts` went uppercase.
- Goal for next time: actually apply `deploy/k8s.yaml` on a local kind cluster
  and watch the readiness probe gate traffic through a rollout.

## 2026-06-24
- Hit both "next up" goals: built `projects/wordcount`, a small `wc` clone in Go
  (flag parsing, stdin/file streaming, table-driven tests), then containerized it
  with a multi-stage Dockerfile down to a scratch image.
- Cleared the Go context fuzziness — wrote `notes/go-context.md` on how cancellation
  propagates down the tree and why you always `defer cancel()`.
- Breadth pass: Linux, networking, regex, Make, Kubernetes, testing, security,
  system design, and an algorithms refresher.
- What clicked: a multi-stage build means the runtime image carries *only* the
  static binary — no Go toolchain shipped. `CGO_ENABLED=0` is what makes `scratch`
  work.
- Goal for next time: a real CI pipeline (lint → test → build) on wordcount.

## 2026-06-23
- Big notes push: TypeScript type system, a SQL indexing cheat sheet, Go
  concurrency, CI/CD, observability, data structures, bash, and HTTP/REST.
- What clicked: composite-index column order (leftmost-prefix) and why functions
  on indexed columns kill the index — finally makes sense from EXPLAIN output.
- Still fuzzy: Go's context cancellation patterns; want to build a real CLI.
- Updated the roadmap — TS and SQL are now solid; added Web/APIs and shell.
- Goal for next time: containerize a small project end to end, then a Go CLI.

## 2026-06-22
- Set up the Upskill repo: README, roadmap, and topic note structure.
- Wrote first-pass notes for JavaScript, Python, Git, Docker, and SQL.
- Started a curated resources list.
- Goal for next session: turn the SQL notes into a real indexing cheat sheet
  using `EXPLAIN ANALYZE` output from a sample database.

<!-- Template for new entries:
## YYYY-MM-DD
- What I worked on.
- What clicked / what's still fuzzy.
- Goal for next time.
-->
