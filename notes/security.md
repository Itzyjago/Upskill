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

## Mindset
- Validate input at the boundary; treat all external data as hostile.
- Fail closed (deny on error), log security events, don't leak details in errors.
