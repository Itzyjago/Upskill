# System design notes

Building blocks I keep reaching for, and the tradeoffs behind each.

## Caching
- Layers: client → CDN → app/in-memory (Redis) → DB. Cache closest to the read.
- Invalidation is the hard part. Strategies: TTL expiry, write-through,
  write-back, cache-aside (lazy load on miss).
- Watch for **stampede** (many misses at once) → use locks/single-flight or
  jittered TTLs. Watch for stale reads.
- **This isn't just theory in this repo anymore** — `idempotency.go`
  (`notes/http.md`) is a real TTL-expiry cache: keyed by `Idempotency-Key`
  instead of a URL, 5-minute TTL, swept opportunistically on write instead of
  a background goroutine. It's a narrower case than a general read cache
  (single-writer-per-key, not "many readers, one writer"), which is exactly
  why it *doesn't* need single-flight/stampede protection — two concurrent
  requests with the *same* key racing each other is a client bug (reusing a
  key for a genuinely concurrent, not sequential-retry, request), not a
  cache-design problem to solve for.
- Where a CDN/reverse-proxy cache *would* sit in front of wordcount: nowhere
  useful today. `/count` is a `POST` (not cacheable by a shared cache
  without explicit opt-in, `notes/http.md`), and the only `GET` worth
  mentioning, `/metrics`, must never be cached — a stale scrape is worse than
  a slow one, it's actively misleading to whoever's alerting on it.

## Load balancing
- Spread traffic across instances: round-robin, least-connections, hashing.
- Health checks eject bad nodes. Sticky sessions only if you can't make the
  service stateless (prefer stateless + shared session store).

## Async & queues
- Decouple producers from consumers with a queue (SQS, RabbitMQ, Kafka).
- Smooths spikes, enables retries and back-pressure. Design consumers to be
  **idempotent** — messages can arrive more than once (at-least-once delivery).

## Data & scaling
- Vertical (bigger box) is simplest; horizontal (more boxes) scales further but
  needs partitioning/sharding and tolerance for eventual consistency.
- Read-heavy? Add read replicas + cache. Write-heavy? Shard by key, batch writes.
- CAP: under a partition you choose availability or consistency, not both.

## Reliability
- Avoid single points of failure; design for graceful degradation.
- Timeouts + retries **with backoff + jitter**; circuit breakers stop cascades.
- Estimate first: QPS, payload size, storage growth — back-of-envelope before code.

## Case study: scaling wordcount past one edge + one upstream
The compose stack (`deploy/observability/docker-compose.yml`) runs exactly two
wordcount instances: an edge that forwards `/count` to a single upstream via
one hardcoded `WORDCOUNT_UPSTREAM_URL` (`client.go`). Fine for proving the
two-service trace works (roadmap #12); not fine as a real topology. Where it
actually breaks:
- **Single point of failure** — `forwardCountHandler` (`client.go`) returns
  `502` the instant the one upstream is unreachable
  (`TestForwardCountHandlerBadGatewayOnUpstreamDown` proves this on purpose).
  There's no second instance to fail over to.
- **No back-pressure** — the edge calls upstream synchronously, inline with
  the client's request (`upstreamClient.count`, `client.go`). A slow upstream
  makes the edge slow; there's nothing absorbing a burst.
- **No load spreading** — even with more upstream replicas, one URL means
  every request hits the same instance. `docker-compose.yml`'s DNS round-robin
  doesn't apply here since it's a single named service, not a pool.

This is a **load balancer** problem, not a **queue** problem, and the two
solve different failure modes:
- A **queue** in front of the forward hop (SQS/RabbitMQ-shaped) would decouple
  the edge from upstream latency and let a burst queue up instead of
  timing out — but `/count` is a synchronous request/response API: the
  client is waiting on the HTTP connection for an answer. Queueing the
  forward hop means either blocking the client anyway (defeats the point) or
  turning `/count` into "202 Accepted, poll a result later" — a real API
  contract change, not a drop-in fix.
- A **load balancer** — multiple upstream replicas behind one address, health
  checks ejecting a dead one, round-robin or least-connections spreading load
  — fixes the actual problems here (SPOF, no spreading) *without* changing
  the request/response contract. `client.go` already does the hard part
  (propagating the trace context, handling upstream errors as `502`); it just
  needs `WORDCOUNT_UPSTREAM_URL` to resolve to more than one healthy backend.
- **Verdict for this shape of service**: load balancer, not queue. The queue
  answer is right for background/async work (the Alertmanager webhook sink,
  `webhook.go`, is closer to that shape); a synchronous count-and-respond API
  wants spread-and-failover, not decoupling.
