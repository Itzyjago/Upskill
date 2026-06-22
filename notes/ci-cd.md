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
