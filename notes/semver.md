# Semantic versioning notes

Wrote these after setting up the tag-triggered release workflow — the image
tags come straight from the git tag.

## MAJOR.MINOR.PATCH
- **MAJOR** — incompatible/breaking API changes.
- **MINOR** — new functionality, backward-compatible.
- **PATCH** — backward-compatible bug fixes.
- Bump the leftmost field that applies and **reset** the ones to its right
  (1.4.2 → a breaking change → 2.0.0).

## Pre-release & build metadata
- `1.0.0-alpha.1`, `1.0.0-rc.2` — a pre-release **sorts lower** than the
  release: `1.0.0-rc.1 < 1.0.0`.
- `+build.5` is build metadata and is **ignored** for ordering.

## Zero-based (0.y.z)
- `0.x` means "anything may change." Many tools treat a `0.MINOR` bump as
  potentially breaking — don't assume 0.x is stable.

## Version ranges (consumers)
- `^1.2.3` — compatible: `>=1.2.3 <2.0.0` (no major bump).
- `~1.2.3` — patch-level: `>=1.2.3 <1.3.0`.
- Caret with `0.x` is special-cased narrow: `^0.2.3` → `>=0.2.3 <0.3.0`.

## Git tags drive releases
- Annotated tag `v1.2.3` is the source of truth; the release workflow's
  `docker/metadata-action` turns it into image tags `1.2.3`, `1.2`, and the sha.
- Keep a `v` prefix on tags (`v1.2.3`) — Go modules and most tooling expect it.

## Practical
- Tag from a green `main`; never re-point a published tag (consumers cache it).
- Pair tags with a short CHANGELOG entry so a version means something to a human.
