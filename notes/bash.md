# Bash / shell notes

## Safer scripts start with this
```bash
set -euo pipefail
```
- `-e` exit on any command failure, `-u` error on unset variables,
  `-o pipefail` make a pipeline fail if *any* stage fails (not just the last).
- Quote your variables: `"$var"` — unquoted expansion word-splits and globs.

## Pipes and redirection
- `cmd > file` overwrites, `>> file` appends.
- `2>&1` sends stderr to wherever stdout currently goes (order matters:
  `> file 2>&1`, not `2>&1 > file`).
- `cmd1 | cmd2` wires stdout of one into stdin of the next.

## Useful constructs
```bash
for f in *.log; do echo "$f"; done
if [[ -f "$path" ]]; then ...; fi      # [[ ]] is the bash test, safer than [ ]
name="${1:-default}"                     # default if $1 unset
files=$(find . -name '*.md')             # command substitution
```
- `${var:-default}`, `${var:?msg}`, `${var%%.*}` (strip suffix) save a lot of
  external calls.

## Exit codes
- `0` = success, non-zero = failure. `$?` holds the last code.
- Chain on success/failure: `a && b` runs b only if a succeeded; `a || b` runs
  b only if a failed.

## Gotchas
- A pipeline runs in a subshell — variables set inside `cmd | while read` don't
  survive. Use process substitution `while read ...; do ...; done < <(cmd)`.
- `cd` in a script can fail silently without `set -e`; check it.
