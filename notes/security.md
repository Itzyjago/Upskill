# Security notes

Defensive basics for app developers — the stuff that actually bites.

## Injection
- Never build queries/commands by string concat. Use **parameterized queries**
  (prepared statements). Same idea for shell: avoid `sh -c "$userInput"`.
- Output-encode for the sink: HTML-escape for the DOM to stop XSS; a templating
  engine with auto-escaping is your friend.

## AuthN vs AuthZ
- Authentication = who you are; authorization = what you may do. Check *both* on
  every request — don't trust a hidden field or client-side role.
- Prefer short-lived tokens; validate signature, issuer, audience, and expiry.

## Passwords & secrets
- Hash passwords with a slow, salted KDF: **bcrypt / scrypt / Argon2**. Never
  MD5/SHA-1, never plaintext.
- Secrets go in env/secret managers, never in git. Rotate on exposure. Add a
  pre-commit secret scanner.

## Transport & data
- TLS everywhere; HSTS to prevent downgrade. Validate certs (don't disable!).
- Principle of least privilege for DB users, cloud IAM, API tokens.

## Common OWASP-ish pitfalls
- Broken access control (IDOR): `/api/users/123` — verify the caller owns 123.
- SSRF: validate/allow-list outbound URLs the server fetches on user input.
- Dependency risk: pin versions, run `npm audit` / `govulncheck`, update often.

## Resource exhaustion (the DoS nobody thinks is "security")
- An **unbounded read** — `io.ReadAll(r.Body)` with no limit — isn't just
  sloppy, it's a footgun: one request with a multi-GB body (or a body that
  never ends) can OOM the process before any application logic runs. It's the
  memory-exhaustion sibling of an injection bug: untrusted input controls a
  resource the server should own.
- `http.MaxBytesReader(w, r.Body, limit)` (Go's stdlib answer) wraps the body
  so a read past `limit` fails fast with an error instead of buffering
  forever. Cheap insurance — one line at the boundary, applies whether the
  handler is `count`-ing 12 bytes or a `io.ReadAll` of arbitrary size.
- This is a **boundary** control, same category as parameterized queries and
  output encoding above: the size cap belongs where untrusted bytes first
  enter the process (the handler), not scattered through downstream code that
  assumes a well-behaved caller.
- Applies to slice/map growth from user input too, not just body reads — same
  shape of bug: "the size of this resource is attacker-controlled and
  unbounded" is the actual vulnerability, HTTP bodies are just the most common
  door it walks in through.

## Reuse audit: is the fix actually everywhere?
`countHandler`/`forwardCountHandler` got the `MaxBytesReader` fix (above);
went back and actually `grep`ed the whole codebase for every other place a
body gets read, instead of assuming the pattern had propagated. It hadn't —
two more unbounded reads turned up, and they're the two directions the
`/count` fix doesn't cover:
- **`webhook.go`'s alert sink** — `alertWebhookHandler` decoded Alertmanager's
  POST body with no cap at all. Same bug, same fix, one detail worth keeping:
  it was originally wired as `json.NewDecoder(MaxBytesReader(...)).Decode()`
  directly, and testing that turned up that `errors.As(err,
  *http.MaxBytesError)` **doesn't reliably survive going through
  `json.Decoder`** — whether the cap trip reaches the caller as a
  `*MaxBytesError` or gets reshaped into a `SyntaxError` first depends on
  exactly where in the JSON the cutoff lands (verified both ways). Switched to
  `io.ReadAll(MaxBytesReader(...))` first, `json.Unmarshal` the resulting
  bytes second — the same two-step split `countHandler` already used, now
  understood as load-bearing rather than incidental.
- **`client.go`'s upstream response** — `upstreamClient.count` decoded the
  *upstream's* response with no cap either. Easy to miss because it's not
  "untrusted user input" in the usual sense — it's the mirror case: a
  compromised or just-buggy upstream can exhaust the *caller* the same way a
  hostile client can exhaust a server. `http.MaxBytesReader` doesn't apply
  here (it needs a `ResponseWriter` to signal the trip to, which only exists
  server-side); `io.LimitReader(resp.Body, cap)` is the client-side
  equivalent — no distinct error type, it just truncates, which surfaces
  naturally as a JSON decode error either way.
- **The actual lesson**: "apply the same fix everywhere" isn't a checklist
  item, it's a search — the two misses here weren't in the code that was
  already being looked at when the original bug was fixed, they were in
  adjacent files nobody re-audited. Both are now covered by tests
  (`TestAlertWebhookHandlerRejectsOversizedBody`,
  `TestUpstreamClientCapsOversizedResponse`) instead of just inspection.

## Mindset
- Validate input at the boundary; treat all external data as hostile.
- Fail closed (deny on error), log security events, don't leak details in errors.
