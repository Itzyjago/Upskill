# Learning Roadmap

Status legend: `🟢 solid` · `🟡 in progress` · `⚪ not started`

## Languages
- 🟢 JavaScript / TypeScript — async, modules, the type system
- 🟡 Python — stdlib fluency, packaging, virtualenvs
- 🟢 Go — concurrency, context cancellation, built a small CLI

## Foundations
- 🟢 Git — branching, rebase, recovering from mistakes
- 🟢 SQL — joins, indexes, query planning
- 🟡 Data structures — revisit trees and hash maps
- 🟡 Shell scripting — bash strict mode, expansion, pipelines
- 🟢 Algorithms — big-O, search/sort, common patterns
- 🟡 Linux — processes, signals, permissions, file descriptors
- 🟡 Regular expressions — groups, lookarounds, greedy vs lazy
- 🟡 Make — task running, phony targets, automatic variables

## Web / APIs
- 🟡 HTTP & REST — methods, status codes, idempotency, caching
- 🟡 Networking — TCP/UDP, DNS, the TLS handshake

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
- 🟡 Testing — the pyramid, table-driven tests, doubles
- 🟡 Security — injection, authz, secrets, OWASP basics
- 🟡 System design — caching, load balancing, queues, scaling

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

### Next up
16. Security: the `/count` body read is unbounded — cap it, and write up why
    that's a resource-exhaustion footgun, not just a style nit.
