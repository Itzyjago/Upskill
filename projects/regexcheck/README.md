# regexcheck

Makes `notes/regex.md`'s Go-vs-lookaround claim runnable instead of asserted.

- `parsePrice` uses capture groups (`^(\$)?(\d+\.\d+)( each)?$`) instead of
  lookaround to pull the same fields `(?<=\$)\d+\.\d+(?= each)` would in a
  backtracking engine — the "what to do instead in Go" from the notes, as
  real code with table tests.
- `TestRE2RejectsLookahead` / `TestRE2RejectsLookbehind` pin the actual
  compile error text RE2 returns for lookaround syntax, so the notes' claim
  is a regression test, not just prose someone has to trust.

```
go run . '$42.50 each'
"$42.50 each" -> dollar=true amount=42.50 each=true
```
