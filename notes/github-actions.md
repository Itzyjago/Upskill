# GitHub Actions notes

Wrote these while standing up the wordcount pipeline (see
[ci-cd.md](ci-cd.md) for the vendor-neutral concepts).

## The object model
- **Workflow** — a YAML file in `.github/workflows/`. Triggered by events.
- **Job** — runs on a fresh runner (VM). Jobs are parallel by default.
- **Step** — a shell command (`run:`) or a reusable **action** (`uses:`).
- Each job starts clean — nothing carries between jobs except artifacts and
  declared outputs.

## Ordering: needs
- Jobs run in parallel unless one declares `needs: [other]`.
- That's how you get `lint → test → build`: each stage `needs` the prior one,
  so a failure short-circuits the rest.

## Triggers
```yaml
on:
  push:
    branches: [main]
    paths: ["projects/wordcount/**"]   # only run when this dir changes
  pull_request:
```
- `paths` filters keep a monorepo from running every workflow on every push.

## Caching is the big speedup
- `actions/setup-go@v5` (and setup-node, etc.) cache the module/dep download
  when you point them at the lockfile (`cache-dependency-path: .../go.sum`).
- Cache keys should hash the lockfile — invalidate only when deps change.

## Matrix builds
```yaml
strategy:
  matrix:
    go: ["1.21", "1.22"]
    os: [ubuntu-latest, macos-latest]
runs-on: ${{ matrix.os }}
```
- One job definition fans out across every combination.

## Gotchas I hit
- `working-directory` under `defaults.run` applies to `run:` steps but **not**
  to `uses:` actions — those take their own `working-directory` input.
- Secrets are masked in logs but available as `${{ secrets.NAME }}`; never
  `echo` one expecting it hidden in a pasted-out value.
- `GITHUB_TOKEN` is auto-provided per run with repo scope — no PAT needed for
  most things (see [security.md](security.md) on secrets).
