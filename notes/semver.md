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

## Checked against the real workflow
Went back to actually verify these notes against `release.yml` and `git
tag -l` instead of trusting what I'd written down weeks ago. Two things
turned up:
- **No tag has ever been pushed.** `git tag -l` on this repo is empty — the
  release workflow (`.github/workflows/release.yml`) has never actually run.
  Everything above about `docker/metadata-action` deriving `1.2.3`/`1.2`/sha
  tags is correct by reading the config (`type=semver,pattern={{version}}`,
  `type=semver,pattern={{major}}.{{minor}}`, `type=sha` — matches), but it's
  unverified in practice. Worth remembering the difference between "the
  config says this will happen" and "I watched it happen."
- **The workflow's own comment contradicted these notes.** It said `git tag
  v0.1.0` — a **lightweight** tag (just a ref, no tagger/message/date) — while
  this file has said "annotated tag is the source of truth" since it was
  written. `docker/metadata-action` doesn't actually care which kind you
  push, so this wasn't a functional bug, just an inconsistency between the
  advice and the copy-pasteable command sitting right next to it. Fixed the
  comment to `git tag -a v0.1.0 -m "v0.1.0"` so the two agree.
