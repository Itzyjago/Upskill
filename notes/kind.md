# kind — Kubernetes in Docker

Wrote these to finally *run* wordcount's `deploy/k8s.yaml` (roadmap #5) instead
of just writing the manifest. **kind** spins up a throwaway cluster as Docker
containers — each "node" is a container running kubelet. Perfect for local
probe/rollout experiments.

## The loop
```sh
kind create cluster --name upskill          # ~30s; sets kubectl context
kubectl cluster-info --context kind-upskill
kind delete cluster --name upskill          # tear down, leaves nothing behind
```

## The gotcha: images must be *loaded*
A kind node can't see your host's Docker images — it has its own container
runtime. A locally-built `wordcount:latest` will sit in `ImagePullBackOff`
forever unless you load it in:
```sh
docker build -t wordcount:latest .
kind load docker-image wordcount:latest --name upskill
```
And set `imagePullPolicy: IfNotPresent` (or `Never`) so it doesn't try to pull
`:latest` from a registry that doesn't have it. This bit me — the default pull
policy for the `:latest` tag is `Always`.

## Watching the readiness probe gate traffic
The whole point of #5 — see readiness actually hold back a rollout:
```sh
kubectl apply -f deploy/k8s.yaml
kubectl get pods -w                         # watch them go Pending → Running → Ready
kubectl rollout status deploy/wordcount     # blocks until the new ReplicaSet is Ready
kubectl port-forward svc/wordcount 8080:80  # reach it from the host
```
- A pod is only added to the Service's endpoints once **readiness** passes, so
  during a rolling update traffic never hits a not-yet-ready pod.
- Force a failure (point readiness at a bad path) and watch the rollout **stall**
  instead of cutting over — exactly the safety the probe buys.

## What clicked
- kind is just Docker — `docker ps` shows the node containers. The "load the
  image" step is the thing everyone trips on, because locally the registry and
  the cluster *feel* like the same machine but aren't.

See [kubernetes.md](kubernetes.md) for the object model and
[health-checks.md](health-checks.md) for probe semantics.
