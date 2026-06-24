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
- 🟡 CI/CD — pipelines, caching, deploy gates
- 🟡 Observability — logs, metrics, traces
- 🟡 Kubernetes — pods, deployments, services, reconciliation

## Cross-cutting
- 🟡 Testing — the pyramid, table-driven tests, doubles
- 🟡 Security — injection, authz, secrets, OWASP basics
- 🟡 System design — caching, load balancing, queues, scaling

## Next up
1. ✅ Containerize one small project end to end → `projects/wordcount` (Dockerfile).
2. ✅ Build a small Go CLI to make the concurrency notes stick → `wordcount`.
3. Stand up a real CI pipeline (lint → test → build) on the wordcount project.
4. Deploy the wordcount container somewhere and add a readiness probe.
