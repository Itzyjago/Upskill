# HTTP & REST notes

## Methods and what they promise
- `GET` — read, no side effects. Safe and idempotent.
- `POST` — create / trigger. Not idempotent (calling twice may make two things).
- `PUT` — replace a resource wholesale. Idempotent.
- `PATCH` — partial update.
- `DELETE` — remove. Idempotent (deleting twice ends in the same state).

## Status codes worth knowing cold
- `2xx` success — `200 OK`, `201 Created`, `204 No Content`.
- `3xx` redirect — `301` permanent, `302/307` temporary, `304 Not Modified`.
- `4xx` you messed up — `400` bad request, `401` unauthenticated,
  `403` authenticated-but-forbidden, `404` not found, `409` conflict,
  `422` validation failed, `429` rate limited.
- `5xx` server messed up — `500` generic, `502/503/504` upstream/unavailable/timeout.

## Idempotency in practice
- Network retries are inevitable; design writes so a retry is safe.
- For non-idempotent `POST`, accept an `Idempotency-Key` header and dedupe.

## Caching
- `Cache-Control: max-age=...` for freshness; `ETag` + `If-None-Match` for
  cheap revalidation (server replies `304`, no body).

## REST design
- URLs are nouns (`/users/42/orders`), HTTP methods are the verbs.
- Use query params for filtering/paging, not new endpoints.
- Return the right status code — don't `200 OK` with `{"error": ...}` inside.
- Version the API (`/v1/...`) so you can evolve without breaking clients.
