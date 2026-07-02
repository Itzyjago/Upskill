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

## Gotchas
- `CrashLoopBackOff` → check `logs --previous` and the readiness/liveness config.
- No `limits` set → noisy-neighbor pods can starve a node.
- HPA silently does nothing if metrics-server isn't installed — `kubectl
  describe hpa` shows `unable to get metrics`, not an error loud enough to
  notice by accident.
