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
- 🟡 Observability — metrics + structured logs shipped (`/metrics`, slog);
  tracing still to wire
- 🟡 Kubernetes — pods, deployments, services, probes (ran the manifest on a
  local kind cluster)
- 🟢 Health checks — liveness vs readiness vs startup probes (watched readiness
  gate a rollout on kind)
- 🟡 Prometheus — metric types, exposition format, PromQL (hand-rolled an
  exporter; haven't stood up a server to scrape it yet)

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
6. 🟡 Add observability to the serve mode → `/metrics` (Prometheus format) and
   structured slog request logging are in; **traces still to wire**.
7. Wire OpenTelemetry tracing into the serve mode (`otelhttp` + a `trace_id` on
   the log line) — the remaining half of #6.
8. Stand up a real Prometheus + Grafana locally to actually scrape `/metrics`
   and graph the golden signals with the PromQL from `notes/promql.md`.
