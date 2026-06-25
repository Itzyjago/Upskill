# Progress log

A running journal — newest first. One short entry per session.

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
