# The four golden signals

From Google's SRE book — if you can only instrument four things, instrument
these. Wrote this to decide *what* wordcount's `/metrics` should expose, not just
*how* (see [prometheus.md](prometheus.md)).

## The four
1. **Latency** — how long requests take. Measure successful and failed requests
   **separately**: a fast stream of 500s can otherwise hide behind a healthy
   average. Track distributions (histograms → p50/p95/p99), never just the mean —
   the mean lies when there's a long tail.
2. **Traffic** — demand on the system. Requests/sec, transactions/sec. This is
   the denominator for almost everything else (error *rate*, saturation).
3. **Errors** — rate of failed requests. Explicit (5xx), implicit (200 with the
   wrong body), or policy (too slow counts as failed). Often the hardest to
   define honestly.
4. **Saturation** — how "full" the system is. The constrained resource (CPU,
   memory, I/O, connection pool, queue depth). Utilization predicts latency
   cliffs *before* they show up as errors.

## How wordcount maps onto them
- **Latency** → `http_request_duration_seconds` histogram, by route.
- **Traffic** → `rate(http_requests_total[1m])`.
- **Errors** → the same counter filtered to `status=~"5.."`, over total.
- **Saturation** → `http_requests_in_flight` gauge (a stand-in for "how busy");
  in a real service you'd also watch CPU and the goroutine count.

## RED and USE — two cousins
- **RED** (for *request-driven services*): **R**ate, **E**rrors, **D**uration.
  Basically signals 2,3,1 — the per-endpoint view.
- **USE** (for *resources*): **U**tilization, **S**aturation, **E**rrors. Aim
  this at CPUs, disks, pools.
- Rule of thumb: RED for the service, USE for the machine under it.

## What clicked
- Averages are a trap — p99 latency is a different *question* than mean latency,
  and SLOs live at the tail. A histogram answers both; a single gauge of "avg
  latency" answers neither honestly.
- Saturation is the leading indicator. Errors and latency are lagging — by the
  time they spike, saturation already told you.
