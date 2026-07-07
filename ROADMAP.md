# Learning Roadmap

Status legend: `🟢 solid` · `🟡 in progress` · `⚪ not started`

## Languages
- 🟢 JavaScript / TypeScript — async, modules, the type system
- 🟢 Python — stdlib fluency (applied in `scripts/parse_access_log.py`),
  packaging, virtualenvs
- 🟢 Go — concurrency, context cancellation, built a small CLI

## Foundations
- 🟢 Git — branching, rebase, recovering from mistakes
- 🟢 SQL — joins, indexes, query planning
- 🟢 Data structures — trees and hash maps, plus a real one traced end to end
  (`metrics.go`'s `map[labelKey]int64` and why `sortedKeys()` exists — map
  iteration order is randomized, not just unspecified)
- 🟡 Shell scripting — bash strict mode, expansion, pipelines
- 🟢 Algorithms — big-O, search/sort, common patterns
- 🟡 Linux — processes, signals, permissions, file descriptors
- 🟢 Regular expressions — groups, lookarounds, greedy vs lazy, and the Go
  trap (RE2 rejects lookaround outright) turned into permanent tests:
  `scripts/regex_lookaround.py` (PCRE side) + `projects/regexcheck` (RE2
  rejection + the capture-group workaround)
- 🟢 Make — task running, phony targets, automatic variables, and a real
  three-deep phony chain from `Makefile` (`image -> kind-load -> kind-deploy`)

## Web / APIs
- 🟢 HTTP & REST — methods, status codes, idempotency (a real
  `Idempotency-Key` store wired into `/count`, `idempotency.go`), and why
  `ETag`/`304` caching doesn't fit a `POST` in the first place
- 🟢 Networking — TCP/UDP, DNS, the TLS handshake, and `projects/netcheck`
  verifying DNS + the TCP open against a real host instead of trusting prose

## Platform / DevOps
- 🟢 Docker — images vs containers, multi-stage builds (containerized the CLI)
- 🟢 CI/CD — built a real lint → test → build pipeline + tag-based release
- 🟢 Observability — all three pillars *exported*: `/metrics`, structured slog,
  and real OTLP spans shipped to Jaeger (hand-rolled OTLP/HTTP, no SDK) — an
  actual waterfall, not just a propagated `traceparent`
- 🟢 Kubernetes — pods, deployments, services, probes (ran the manifest on a
  local kind cluster), and a HorizontalPodAutoscaler (needs metrics-server,
  not installed on kind by default)
- 🟢 Health checks — liveness vs readiness vs startup probes (watched readiness
  gate a rollout on kind)
- 🟢 Prometheus — metric types, exposition format, PromQL, recording rules, *and*
  a real server: a compose stack scrapes `/metrics`, Grafana graphs the golden
  signals, recording rules precompute them, alert rules fire
- 🟢 Alerting — Prometheus rules + the `for:` window, and Alertmanager routing a
  page to a real webhook receiver (grouping, inhibition, the routing tree);
  also provisioned the same alert as a Grafana-native rule to compare the two
  engines' shapes
- 🟢 Distributed tracing — W3C trace context propagation + OTLP/HTTP export to
  Jaeger; spans carry timings, the log line cross-links via `trace_id`; a real
  two-service trace (client span forwards to a second instance, server span on
  the far side parents to it); an OTel Collector tail-samples (errors + slow
  always kept, 10% baseline) before spans reach Jaeger

## Cross-cutting
- 🟢 Testing — the pyramid, table-driven tests, and real doubles: an
  `httptest`-backed fake for both outbound calls this repo makes
  (`upstreamClient` and, now, the OTLP collector via `middleware_test.go`'s
  `fakeCollector`) instead of nil-ing dependencies out
- 🟢 Security — injection, authz, secrets, OWASP basics, resource exhaustion;
  applied for real three times now (the original `/count` body cap, plus a
  follow-up audit that found and fixed the same unbounded-read bug in
  `webhook.go` and `client.go`'s upstream response, each covered by a test)
- 🟡 System design — caching, load balancing, queues, scaling; worked
  wordcount's own scaling story through concretely (`notes/system-design.md`)

## Next up
1. ✅ Containerize one small project end to end → `projects/wordcount` (Dockerfile).
2. ✅ Build a small Go CLI to make the concurrency notes stick → `wordcount`.
3. ✅ Stand up a real CI pipeline (lint → test → build) → `.github/workflows/ci.yml`.
4. ✅ Deploy the wordcount container + readiness probe → applied `deploy/k8s.yaml`
   on a local kind cluster.
5. ✅ Run it on a local cluster and watch the readiness probe gate traffic during
   a rollout → done on kind (`make kind-deploy`, `kubectl rollout status`).
6. ✅ Add observability to the serve mode → `/metrics`, structured slog logging,
   and trace context are all in.
7. ✅ Wire tracing into the serve mode → hand-rolled W3C `traceparent`
   (`trace.go`): continue/extract a trace, mint a child span, stamp
   `trace_id`/`span_id` on the log line. SDK-free, like the metrics exporter.
8. ✅ Stand up Prometheus + Grafana locally → `deploy/observability/` compose
   stack scrapes `/metrics`, a provisioned Grafana dashboard graphs the golden
   signals, and alert rules load (`make obs-up`).
9. ✅ Export real spans, not just the header → hand-rolled OTLP/HTTP JSON exporter
   (`otlp.go`) ships a timed span per request to Jaeger; `trace.go` now records
   start/end + parent id, and the compose stack runs Jaeger (UI on :16686).
10. ✅ Add Alertmanager and route an alert somewhere real → Alertmanager in the
    stack, `alertmanager.yml` routes `severity=page` to a webhook, and a `wc
    -webhook` sink logs the routed alert (reuses the binary).
11. ✅ Recording rules — `rules/recording.yml` precomputes the golden signals
    (`job:...` series); the alerts and the error-ratio panel read them so the
    dashboard and the alert can't drift apart.

12. ✅ Export real OTLP spans for the *outbound* hop too — `client.go` forwards
    `/count` to a second wordcount instance wrapped in a client span; the
    compose stack runs both and Jaeger shows a real two-service waterfall.

13. ✅ Add tail sampling at the collector — an OTel Collector now sits between
    both wordcount instances and Jaeger; `tail_sampling` keeps every errored or
    >500ms trace and a flat 10% baseline of the rest (`otel-collector.yaml`).

14. ✅ Grafana alerting vs. Prometheus/Alertmanager — provisioned a
    Grafana-native mirror of `HighErrorRate` (alert rule, contact point,
    notification policy) next to the Prometheus/Alertmanager original; same
    shape, different engine — `notes/grafana.md`.

15. ✅ Kubernetes: a HorizontalPodAutoscaler for the deployment — 2-5 replicas,
    70% CPU target; `make kind-metrics-server` covers kind's missing add-on.

16. ✅ Security: the `/count` body read was unbounded — `http.MaxBytesReader`
    caps it at 10MB, `statusForBodyErr` turns a tripped cap into a `413`
    instead of a generic `400` (`notes/security.md`).

17. ✅ Testing: the OTLP export goroutine (`middleware.go`) was untestable —
    `otlpExporter` was just nil'd out in every test, an escape hatch, not a
    double. `middleware_test.go`'s `fakeCollector` (an `httptest.Server` that
    signals a channel on receipt) now asserts export actually happens, no
    production code changed (`notes/testing.md`).
18. ✅ System design: wordcount past one edge + one upstream is a **load
    balancer** problem, not a queue one — `/count` is synchronous
    request/response, so decoupling with a queue means either blocking
    anyway or changing the API contract; multiple upstream replicas behind
    one address doesn't (`notes/system-design.md`).

### Next up
19. ⛔ **Blocked on tooling, not skipped.** Actually apply the load-balancer
    verdict from #18 — `deploy/k8s.yaml`'s upstream is still one Service
    backing one Deployment; scale replicas and confirm the Service's built-in
    load balancing (kube-proxy) spreads `/count` traffic across pods, not
    just liveness/readiness gating a rollout. This machine has no
    Docker/kind/kubectl (checked 2026-07-08: `go`/`python`/`node` resolve,
    those three don't), so there's no cluster here to actually watch
    kube-proxy spread traffic — marking this ✅ without that would be
    exactly the "assume it's fine" failure mode the scratch log exists to
    catch. Stays open until it runs on a machine that has the tooling.
20. ✅ Data structures: was 🟡 "revisit trees and hash maps" with nothing to
    show for it — turned out this repo already had a real hash map doing
    work (`metrics.go`'s `map[labelKey]int64`), including a live example of
    "unordered by design" (`sortedKeys()` exists because map iteration order
    is randomized, not just unspecified) — `notes/data-structures.md`.
21. ✅ HTTP: idempotency notes were abstract — worked `/count` through for
    real. Its *response* is idempotent (pure function of the body); its
    *metrics* aren't (every call increments counters, retry or not); and
    `upstreamClient` has no retry logic yet, which is exactly why adding one
    later needs an `Idempotency-Key` design up front, not bolted on after a
    double-counted metric shows up (`notes/http.md`).

### Next up
22. ✅ Built the `Idempotency-Key` design #21 said would be needed —
    `idempotency.go` is a real in-memory store wired into `/count`
    (`countHandlerFunc`): same key + same body replays the cached response,
    same key + different body is a `409`. Applied at the client-facing hop
    only; wiring it into `forwardCountHandler` waits for `upstreamClient` to
    actually grow retry logic, per #21's own rule about not bolting this on
    speculatively (`notes/http.md`).
23. ✅ Networking: DNS/TCP were prose with nothing behind them —
    `projects/netcheck` resolves a host and times a raw TCP connect, tested
    with an injected resolver/dialer (no real socket in CI), and run for
    real against `example.com` for actual numbers (`notes/networking.md`).
24. ✅ Regex: the lookaround examples and the RE2-rejects-them claim were
    verified once by hand and written down — now permanent tests instead:
    `scripts/regex_lookaround.py` (PCRE side, Python's `re`) and
    `projects/regexcheck` (RE2's actual compile-error text, plus the
    capture-group workaround as real code) (`notes/regex.md`).
