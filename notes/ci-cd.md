# CI/CD notes

## CI vs CD
- **CI** (Continuous Integration): every push runs build + tests so broken code
  is caught in minutes, not at release.
- **CD** can mean *Delivery* (always have a deployable artifact, ship on a
  button) or *Deployment* (every green build auto-ships to prod).

## Pipeline shape
```
lint → test → build → (scan) → deploy
```
- Run cheap, fast jobs first (lint) so failures surface quickly.
- Independent jobs run in parallel; use stages/needs to express order.

## Caching = the biggest speedup
- Cache dependency dirs (`node_modules`, `~/.m2`, pip wheels) keyed by the
  lockfile hash — invalidate only when deps change.
- Use layer caching for Docker builds; order Dockerfile so deps install before
  copying source.

## Deploy gates
- Require green tests + review approval before the deploy job runs.
- Protect prod with a manual approval step or environment protection rules.
- Keep deploys reversible: blue/green or canary, plus a one-command rollback.

## Principles
- A red pipeline blocks merge — keep it fast or people route around it.
- Build the artifact once, promote the *same* artifact through environments.
- Secrets live in the CI secret store, never in the repo (see [git.md](git.md)).

## `version: latest` is a live wire (wordcount CI)
`ci.yml`'s lint job pinned `golangci-lint-action@v6` but left the linter
itself on `version: latest` — a newer release started flagging code that was
already sitting in the repo (unused param, a shadowed builtin), so the
pipeline went red on a push that hadn't touched any of the flagged lines.
Came back to actually fix the root cause instead of just noting it:
`golangci-lint` crossed a v1 → v2 major version, and the action's own docs
only claim v6 is compatible with v1.x — `latest` had silently drifted onto a
v2 release the action doesn't officially support. v2 also restructures
`.golangci.yml` (formatters split out of `linters`), so upgrading the action
to a v2-compatible release would mean migrating the config too — a bigger,
separate piece of work, not a one-line fix.
- **Fix applied**: pinned `version: v1.64.8` (the latest v1.x release) —
  verified locally by installing that exact version and running it against
  the repo; it passes clean against the existing `.golangci.yml`.
- **The lesson**: `latest` isn't "the version that works with my config," it's
  "whatever the tool maintainers shipped this week." Pin CI tooling the same
  way you pin dependencies — a reproducible build includes the linter, not
  just the code.
