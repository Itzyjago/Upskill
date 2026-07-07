# Kubernetes notes

The mental model: you declare desired state; controllers reconcile reality to it.

## Core objects
- **Pod** — one or more containers sharing network + storage. The smallest unit;
  usually not created directly.
- **Deployment** — manages a ReplicaSet → keeps N pod replicas running, handles
  rolling updates and rollbacks.
- **Service** — stable virtual IP + DNS name in front of a set of pods (selected
  by labels). Types: `ClusterIP` (internal), `NodePort`, `LoadBalancer`.
- **Ingress** — HTTP routing (host/path) to Services; needs an ingress controller.
- **ConfigMap** / **Secret** — config and credentials injected as env or files.

## Reconciliation loop
- `kubectl apply -f` writes desired state to the API server (etcd).
- Controllers watch and converge actual → desired. Self-healing: kill a pod and
  the Deployment recreates it.

## Health & rollout
- `livenessProbe` (restart if unhealthy), `readinessProbe` (remove from Service
  until ready).
- `kubectl rollout status/undo deployment/<name>` — watch or revert a rollout.
- Set resource `requests` (scheduling) and `limits` (cap) per container.

## Day-to-day
- `kubectl get pods -o wide`, `describe pod` (events!), `logs -f`, `exec -it`.
- Namespaces isolate environments; `kubectl config set-context` to switch.

## Autoscaling
- **HorizontalPodAutoscaler (HPA)** — watches a metric (CPU/memory %, or a
  custom metric) on a Deployment and adjusts `replicas` to hold it near a
  target. Reconciliation, same as everything else here: HPA doesn't scale
  directly, it just edits `spec.replicas` and the Deployment controller does
  the rest.
- It needs something serving the metrics API to watch — `resources.requests`
  on the container (the denominator for a CPU **percentage**) *and* a metrics
  source. On a real cluster that's usually already there; on **kind** it isn't
  — the **metrics-server** add-on has to be installed separately, and its
  default kubelet TLS verification needs `--kubelet-insecure-tls` to work
  inside kind's non-standard networking. `kubectl top pods` returning nothing
  ("metrics not available yet") is the tell.
- `minReplicas`/`maxReplicas` bound it; `targetCPUUtilizationPercentage: 70`
  means "average CPU usage across pods; add a replica if usage runs high,
  remove one if it runs low, don't go outside the bounds."
- Without `resources.requests.cpu` set, a CPU-based HPA has nothing to compute
  a percentage *of* — this is the other reason `requests` matter beyond
  scheduling.

### The fuzzy bit, revisited: what actually happens between polls
Going back over this after `deploy/k8s.yaml`'s HPA had been running a while,
the part that was fuzzy wasn't metrics-server itself, it was the *timing* —
`kubectl describe hpa` showing a stable replica count even while load
visibly changed looked like it was stuck, not working.
- The HPA controller polls the metrics API on a fixed interval (15s default,
  cluster-wide, not per-HPA) — it isn't reacting to every scrape, it's
  sampling. A CPU spike between polls can be invisible for up to that long.
- `deploy/k8s.yaml`'s HPA sets no `behavior` block, so both directions use the
  API's defaults — and the two defaults are deliberately asymmetric:
  scale-**up** has effectively no stabilization window (react fast, a
  loaded pod needs help now), scale-**down** stabilizes over the **last 5
  minutes** of recommendations and picks the *highest* replica count seen in
  that window before shrinking. That asymmetry is the actual answer to "why
  did it scale up in 30s but take 5 minutes to scale back down" — it's not
  lag, it's a deliberate anti-flap guard so a brief traffic dip doesn't yank
  capacity right before the next spike.
- `minReplicas: 2` means CPU usage below 70% never drops replicas below 2 —
  the floor isn't a suggestion, `kubectl top pods` idle at ~1% CPU with 2
  replicas sitting there the whole time is expected, not a sign the HPA gave
  up.

## Gotchas
- `CrashLoopBackOff` → check `logs --previous` and the readiness/liveness config.
- No `limits` set → noisy-neighbor pods can starve a node.
- HPA silently does nothing if metrics-server isn't installed — `kubectl
  describe hpa` shows `unable to get metrics`, not an error loud enough to
  notice by accident.
