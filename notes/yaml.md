# YAML notes

Picked these up the hard way writing CI and k8s manifests.

## It's whitespace-significant
- Indentation is **spaces only** — tabs are a syntax error. Two spaces per
  level is the norm.
- Structure comes from indentation, not braces. A misindented key silently
  becomes a child of the wrong parent.

## Scalars, lists, maps
```yaml
name: upskill          # string
count: 3               # int
enabled: true          # bool
items:                 # list
  - one
  - two
config:                # map
  key: value
```

## Multiline strings
- `|` — **block**, keeps newlines (good for scripts).
- `>` — **folded**, joins lines with spaces (good for long prose).
- Add `-` to strip the trailing newline: `|-`, `>-`.

```yaml
script: |
  set -euo pipefail
  go test ./...
```

## Anchors & aliases (DRY)
```yaml
defaults: &defaults
  retries: 3
  timeout: 30
job_a:
  <<: *defaults          # merge the anchor in
  name: a
```

## The Norway problem (gotchas)
- Unquoted `no`, `yes`, `on`, `off` parse as booleans — Norway's code `NO`
  becomes `false`. **Quote** strings that look like bools/numbers.
- `1.20` parses as the float `1.2`; a Go version pin needs quotes: `"1.20"`.
- Leading zeros (`08`) used to parse as octal — quote version/zip-like values.
- `null`, `~`, and an empty value are all null.

## Rule of thumb
When a string could be misread as another type, quote it. Costs nothing,
prevents the weird bugs.
