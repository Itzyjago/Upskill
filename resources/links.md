# Resources

Curated, not exhaustive — links I've actually found useful.

## Docs (go to the source first)
- MDN Web Docs — JavaScript & web platform
- Python official docs & tutorial
- Pro Git book (free online) — the chapters on branching and rebasing
- PostgreSQL docs — especially the section on EXPLAIN
- TypeScript Handbook — narrowing and the utility types
- A Tour of Go — the official interactive intro
- OpenTelemetry docs — the vendor-neutral observability standard
- MDN HTTP reference — methods, status codes, and caching headers
- Kubernetes docs — the Concepts section, not just the reference
- regex101.com — live regex tester with an explanation pane
- GNU Make manual — automatic variables and pattern rules
- OWASP Cheat Sheet Series — practical, per-topic security guidance
- GitHub Actions docs — workflow syntax, contexts, and the events reference
- Kubernetes "Configure Liveness, Readiness and Startup Probes" task page
- Semantic Versioning 2.0.0 (semver.org) — the spec, it's short
- The official YAML 1.2 spec + "Learn YAML in Y minutes" for the quick version
- Prometheus docs — data model, exposition format, and the "instrumentation
  best practices" page on naming and labels
- PromQL: the official "Querying" docs + Robust Perception's "rate() vs irate()"
- Google SRE book — the "Monitoring Distributed Systems" chapter (golden signals)
- kind (Kubernetes in Docker) "Quick Start" — clusters, loading images, configs
- OpenTelemetry docs — concepts, plus W3C Trace Context for `traceparent`
- Go `log/slog` package docs — handlers, levels, and attrs
- Grafana docs — "Provisioning" (data sources + dashboards as code) and the
  dashboard JSON model
- Prometheus docs — "Configuration" (scrape_configs, relabeling) and "Alerting
  rules" + the Alertmanager overview
- W3C Trace Context Recommendation — the `traceparent`/`tracestate` spec itself
  (short, and the field layout is right there)
- Jaeger / Grafana Tempo docs — where OTel traces actually land and get viewed
- "The Illustrated TLS 1.3 Connection" (tls13.xargs.org) — the actual bytes on
  the wire for a real handshake, not just the flow diagram
- Kubernetes "Horizontal Pod Autoscaling" concept page — the `behavior`
  field's default stabilization windows aren't obvious from the walkthrough
- Go's RE2 syntax reference (`golang.org/s/re2syntax`) — the actual list of
  what `regexp` supports, which is shorter than PCRE/JS and worth checking
  before assuming a pattern will compile
- golangci-lint docs — the config reference, especially useful the moment a
  `version: latest` pin crosses a major version and the schema changes

## Practice
- Exercism — guided exercises with mentor feedback
- LeetCode / HackerRank — data structures & algorithms drills
- "Build your own X" lists — reimplement tools to understand them
- The wordcount CLI in this repo — my own "build your own wc"

## Reading
- "The Pragmatic Programmer" — habits over hype
- "Designing Data-Intensive Applications" — how real systems store & move data

## Newsletters / blogs
- Julia Evans (wizardzines) — short, deep, friendly explainers
- Changelog — weekly pulse on the ecosystem
