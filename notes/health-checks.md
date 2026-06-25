# Health checks & probes

Wrote these while adding `/healthz` to wordcount and wiring it to a Kubernetes
probe. The three probe types do *different* jobs — conflating them causes
restart loops.

## The three probes (Kubernetes)
- **Liveness** — "is the process wedged?" Fail → **kill and restart** the
  container. Use for deadlocks, not for missing dependencies.
- **Readiness** — "can it serve traffic right now?" Fail → **remove from the
  Service endpoints** (stop sending requests) but **don't** restart. Recovers
  on its own when the dep comes back.
- **Startup** — "has a slow starter finished booting?" Holds off liveness/
  readiness until it passes, so slow JVMs/migrations aren't killed early.

## Why the distinction matters
- If a readiness check that depends on a database is used as **liveness**, a
  brief DB blip restarts every pod — turning a small outage into a big one.
- Liveness should be **cheap and local** (process responsive?). Readiness may
  check downstreams (DB, cache) the instance genuinely needs.

## Probe knobs
```yaml
readinessProbe:
  httpGet: { path: /healthz, port: 8080 }
  initialDelaySeconds: 2
  periodSeconds: 5
  failureThreshold: 3      # consecutive fails before it's marked unready
  timeoutSeconds: 1
```
- `failureThreshold * periodSeconds` ≈ how long a flap is tolerated.
- `initialDelaySeconds` (or better, a startup probe) avoids early false fails.

## Graceful shutdown ties in
- On `SIGTERM`: flip readiness to failing, **stop accepting new** connections,
  then **drain** in-flight requests before exit (wordcount does this via
  `signal.NotifyContext` — see [go-context.md](go-context.md)).
- Kubernetes removes the pod from endpoints in parallel with sending SIGTERM,
  so a short drain window prevents dropped requests.

## Beyond k8s
- Same idea behind load-balancer health checks and `docker HEALTHCHECK`.
- Keep the endpoint unauthenticated and dependency-light, or the probe itself
  becomes the outage.
