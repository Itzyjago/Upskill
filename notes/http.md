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

## Idempotency, worked through on a real endpoint (wordcount's `/count`)
The roadmap flagged this as still abstract — `/count` is a `POST` with no
`Idempotency-Key`, so is it idempotent or not? Worked it through instead of
guessing:
- **The response is idempotent.** `countHandler`/`count()` is a pure
  function of the request body — same bytes in, same `{lines, words,
  bytes}` out, no database write, no state that changes what a *second*
  identical call returns. Retrying a timed-out `/count` is safe in the sense
  that matters most: it can't corrupt data, because there's no data to
  corrupt.
- **But it isn't idempotent in every sense** — every call, including a
  retry of an already-succeeded request, increments
  `http_requests_total`/the latency histogram (`metrics.go`). "Same result"
  and "no side effects" are different properties; `/count` only has the
  first. A retry is safe for the *caller*, but it's a phantom extra request
  in the *metrics* — traffic and error-rate golden signals, not the actual
  data.
- **Where this stops being harmless**: `upstreamClient.count` (`client.go`)
  currently has **no retry logic at all** — one `client.Do`, no wrapping —
  so today this is a non-issue in practice. But it's the exact place a retry
  would get added first (network calls are where retries earn their keep),
  and a naive "retry on any network error" there is *not* safe to add
  blindly: the edge doesn't know whether the upstream's request actually
  landed before the connection dropped. Retry, and the upstream may process
  the same text twice — its own `http_requests_total` and duration
  histogram now show two requests for one logical count, quietly inflating
  the very numbers `rules/alerts.yml`'s `HighErrorRate`/`HighLatencyP95`
  watch. A body-content dedupe key wouldn't fix this properly either (two
  different legitimate requests can share identical text); a real fix would
  need a client-generated `Idempotency-Key` header the upstream can dedupe
  server-side by key, not by content.
- **Verdict**: `/count` doesn't need an `Idempotency-Key` today because it
  has no retry path yet and no state to corrupt. It would need one the
  moment `upstreamClient` grows retries — and the design for that key has to
  exist *before* the retry logic does, not be bolted on after a duplicate
  shows up in a Grafana panel.

## REST design
- URLs are nouns (`/users/42/orders`), HTTP methods are the verbs.
- Use query params for filtering/paging, not new endpoints.
- Return the right status code — don't `200 OK` with `{"error": ...}` inside.
- Version the API (`/v1/...`) so you can evolve without breaking clients.
