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
The "without consuming" part is the whole point: a lookaround asserts
something about the surrounding text but isn't part of the match itself, so
it doesn't show up in the result and doesn't get consumed as you scan
forward. Verified each of these against real input, not just written from
memory:
- **Lookahead** `\d+(?= USD)` — digits followed by literal " USD". Against
  `"150 USD"` this matches `150`, not `150 USD`; against `"150 EUR"` it
  doesn't match at all. The " USD" never appears in the captured result —
  that's the "without consuming" part made concrete.
- **Negative lookahead** `foo(?!bar)` — "foo" *not* immediately followed by
  "bar". Matches the "foo" in `"foobaz"`; matches nothing in `"foobar"`. Easy
  to misread as "match anything except foobar" — it's narrower than that: it
  only blocks *this specific* continuation, not "reject the whole string."
- **Lookbehind** `(?<=\$)\d+` — digits preceded by a literal `$`. Against
  `"$150"` matches `150` (the `$` itself isn't in the result); against
  `"USD150"` it doesn't match.
- They compose: `(?<=\$)\d+\.\d+(?= each)` against `"Total: $42.50 each"`
  pulls out exactly `42.50` — price format, but only when both a `$` prefix
  and an " each" suffix are present, neither of which end up in the match.

### The trap this repo would actually hit: Go doesn't have any of the above
All three examples above are JS-flavor (`node`'s regex engine, PCRE-like).
Go's `regexp` package is **RE2**, and RE2 rejects lookaround syntax outright
— `regexp.Compile(`\d+(?= USD)`)` fails to compile with `invalid or
unsupported Perl syntax`, and the lookbehind form fails as an "invalid named
capture" (RE2 parses `(?<=...)` as an attempt at a *named group* with an
illegal name, a different error for the same missing feature). This isn't a
performance opinion, it's a hard guarantee RE2 makes: no
backtracking-dependent construct is allowed, because that's precisely what
lets a regex engine go exponential on adversarial input (see the ReDoS note
below) — RE2 trades expressiveness for a linear-time guarantee, full stop.
- **What to do instead in Go**: capture groups plus post-processing.
  `(\$)?(\d+\.\d+)( each)?` captures the surrounding markers instead of
  asserting on them, then ordinary code decides whether the optional groups
  were present. It's more verbose than one clever lookaround, but it's the
  actual answer for wordcount or any other Go code in this repo — reaching
  for a lookaround here is reaching for syntax that doesn't exist.

## Gotchas
- Escape regex metachars in literals: `.` `*` `+` `?` `(` `)` `[` `{` `\` `|` `^` `$`.
- Flavors differ — PCRE/JS/Go(RE2) aren't identical. RE2 (Go) drops
  backreferences/lookarounds for guaranteed linear time (see above).
- **ReDoS**: nested/overlapping quantifiers (`(a+)+b` against a long run of
  `a`s with no trailing `b`) can make a backtracking engine's match time
  explode — the exact failure mode RE2's no-backtracking guarantee exists to
  rule out. A reason to actually prefer Go's `regexp` for anything parsing
  untrusted input, not just a consolation prize for missing features.
- Test against real data; prefer a parser over regex for nested structures (HTML).
