# System design notes

Building blocks I keep reaching for, and the tradeoffs behind each.

## Caching
- Layers: client → CDN → app/in-memory (Redis) → DB. Cache closest to the read.
- Invalidation is the hard part. Strategies: TTL expiry, write-through,
  write-back, cache-aside (lazy load on miss).
- Watch for **stampede** (many misses at once) → use locks/single-flight or
  jittered TTLs. Watch for stale reads.

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
