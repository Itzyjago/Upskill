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

## `set -e` is narrower than "strict mode" sounds
Verified each of these against a real shell rather than trusting the mental
model — `-e`'s actual rule is "abort if a command's exit status isn't being
*checked* by something," and "checked" has more exceptions than it looks
like at first:
- **Inside an `if`/`while`/`until` condition, `-e` never fires** —
  `if false; then ...; fi` runs the rest of the script normally. The whole
  point of testing a command's exit status is to react to failure yourself;
  `-e` gets out of the way.
- **A command on the left of `&&`/`||` is exempt too** — `false && echo x`
  doesn't abort the script (it just skips `echo x`, ordinary `&&`
  short-circuiting, unrelated to `-e`). This is the one that bites in CI
  scripts: `some_check && echo "ok"` looks safe, but if `some_check` is
  supposed to be a hard gate, this pattern quietly turns it into a soft one.
- **A failing command inside a function is exempt exactly when the function
  call itself is in one of the contexts above** — `f() { false; echo
  "reached"; }; f || true` prints `"reached"`, because `f`'s own exit status
  is being tested by `|| true`, so nothing inside `f` triggers `-e` either.
  Call the same function plainly — `f() { false; echo "reached"; }; f` — and
  `-e` fires the instant `false` runs, aborting before the `echo` and before
  anything after the call.
- **The practical rule**: `-e` protects "a command failed and nobody's
  handling it," not "a command failed." If failure is meant to be fatal,
  don't wrap the check in `if`/`&&`/`||` at all — let the bare command fail
  and let `-e` do its job; the moment you *do* test it, you've told `-e` you
  own the failure path, on purpose or not.
- All four claims above are now `scripts/verify_set_e_exemptions.sh` —
  assertions, not just prose to trust. `scripts/check-scratch-log.sh` is the
  same discipline applied to something actually useful (warns when
  `scratch-log.md` needs clearing).

## A real strict-mode bug this session's scripts hit
Writing those two scripts on Windows surfaced something `set -euo pipefail`
doesn't protect against at all: line endings. This repo's `.editorconfig`
declares `end_of_line = lf`, but nothing enforced it at *checkout* — with
`core.autocrlf=true`, a fresh clone on Windows would silently rewrite the
committed LF script to CRLF, which breaks `#!/usr/bin/env bash` outright
(`env` looks for a program literally named `bash\r`). Fixed with
`.gitattributes` (`* text=auto eol=lf`, plus an explicit `*.sh` rule) —
found by checking what `git show HEAD:script.sh | xxd` actually contained
instead of assuming the working copy matched the commit.
