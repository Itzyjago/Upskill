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
- 🟢 Observability — all three pillars in: `/metrics`, structured slog, and
  hand-rolled W3C trace context (`traceparent` + `trace_id` on the log line)
- 🟡 Kubernetes — pods, deployments, services, probes (ran the manifest on a
  local kind cluster)
- 🟢 Health checks — liveness vs readiness vs startup probes (watched readiness
  gate a rollout on kind)
- 🟢 Prometheus — metric types, exposition format, PromQL, *and* a real server:
  a compose stack scrapes `/metrics`, Grafana graphs the golden signals, alert
  rules load

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
9. Export real spans, not just the header — push OTLP to Jaeger/Tempo and see an
   actual waterfall (the half `trace.go` doesn't do: it propagates context but
   records no span timings yet).
10. Add Alertmanager to the stack and route one alert somewhere real (a webhook),
    so a firing rule actually notifies instead of just lighting up the UI.
11. Recording rules — precompute the golden-signal expressions so dashboards and
    alerts read cheap, pre-aggregated series instead of re-running `rate()` math.
