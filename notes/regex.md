# Regular expressions notes

## Anchors & boundaries
- `^` start, `$` end (of string, or line in multiline mode).
- `\b` word boundary — `\bcat\b` matches "cat" but not "category".

## Character classes & quantifiers
- `\d \w \s` (digit / word / whitespace), uppercase negates: `\D \W \S`.
- `[a-z]`, `[^...]` negated set.
- `*` 0+, `+` 1+, `?` 0-1, `{2,4}` range. Add `?` to make lazy: `.*?`.

## Greedy vs lazy (the classic trap)
- `<.*>` on `<a><b>` grabs the whole thing (greedy).
- `<.*?>` stops at the first `>` → matches `<a>` then `<b>`.

## Groups
- `( )` capturing — back-reference with `\1`, or by name `(?<year>\d{4})`.
- `(?: )` non-capturing — group without the capture overhead.

## Lookarounds (match without consuming)
- Lookahead: `\d+(?= USD)` — digits followed by " USD".
- Negative lookahead: `foo(?!bar)`.
- Lookbehind: `(?<=\$)\d+` — digits preceded by `$`.

## Gotchas
- Escape regex metachars in literals: `.` `*` `+` `?` `(` `)` `[` `{` `\` `|` `^` `$`.
- Flavors differ — PCRE/JS/Go(RE2) aren't identical. RE2 (Go) drops
  backreferences/lookarounds for guaranteed linear time.
- Test against real data; prefer a parser over regex for nested structures (HTML).
