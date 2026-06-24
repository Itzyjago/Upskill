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

## Gotchas
- `CrashLoopBackOff` → check `logs --previous` and the readiness/liveness config.
- No `limits` set → noisy-neighbor pods can starve a node.
