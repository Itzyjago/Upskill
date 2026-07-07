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

## Follow-up: actually did the v2 migration
The v1 pin above was deliberately the smaller, safer fix — this is the
"bigger, separate piece of work" it deferred, done properly instead of left
as a comment:
- **Config**: `.golangci.yml` gets `version: "2"` at the top, `formatters:`
  becomes its own top-level section (`gofmt`/`goimports` move out of
  `linters.enable`), and `linters.default: none` has to be set explicitly —
  v2 defaults to a "standard" bundle if you don't, which would silently
  enable linters this config never asked for. `run.timeout` and
  `issues.max-issues-per-linter`/`max-same-issues` are unchanged.
- **Action**: bumped `golangci-lint-action` to `v9` (the current release;
  its own docs confirm v2 support starting around v7) and the `version:`
  input to `v2.12.2`, matching what got verified locally.
- **Verified before trusting it, same as the v1 pin**: installed
  `golangci-lint/v2` (note the `/v2` module path — that's how a Go module
  signals a major version bump past v1) locally and ran it against the real
  config. First run surfaced 5 genuinely new findings the v1.64.8 binary
  never caught — 3 `errcheck` (unchecked `Close()` calls, two of them in
  this session's own new test code) and 2 `staticcheck` quickfixes. This is
  the exact "a version bump surfaces more than expected" problem from the
  top of this section, playing out again in miniature, except this time
  caught locally before it ever touched CI.
- **The De Morgan quickfix (QF1001) taught its own lesson**: staticcheck
  flagged `!((c>='0'&&c<='9') || (c>='a'&&c<='f'))` and suggested
  `!A && !B`. Applying it literally left the linter *still* unhappy — each
  negated range comparison inside the AND is itself a `!(range)` shape the
  same rule matches, so chasing it to a true fixed-point would have expanded
  into `(c<'0'||c>'9') && (c<'a'||c>'f')`, which is worse to read than the
  original, not better. The actual fix wasn't a mechanical rewrite at all:
  flip the condition to a positive check with `continue` for the valid case
  and `return false` for everything else — no top-level negation left for
  the rule to match, and it reads better than either negated form
  (`trace.go`'s `isLowerHex`). Not every quickfix is worth chasing to
  completion; sometimes the fix is restructuring past the pattern entirely.
