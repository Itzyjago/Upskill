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

## Mindset
- Validate input at the boundary; treat all external data as hostile.
- Fail closed (deny on error), log security events, don't leak details in errors.
