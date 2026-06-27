# Prometheus — the scrape side (config, targets, relabeling)

[prometheus.md](prometheus.md) is about what a service *exposes* (metric types,
the text format). This is the other half: how the Prometheus **server** finds
targets and pulls those numbers. Wrote it wiring a real `prometheus.yml` to
scrape wordcount.

## Pull, not push
- Prometheus **scrapes** — it reaches out to each target's `/metrics` on a timer
  (`scrape_interval`) and stores what it gets. Targets don't push.
- Why pull wins for monitoring: Prometheus knows the full target list, so a
  *missing* scrape is itself a signal (`up == 0`). A push model can't tell "all
  healthy" from "silently dead."

## A minimal config
```yaml
global:
  scrape_interval: 5s        # how often to pull
  evaluation_interval: 5s    # how often to evaluate alert/recording rules

scrape_configs:
  - job_name: wordcount
    metrics_path: /metrics
    static_configs:
      - targets: ["wordcount:8080"]
```
- A **job** is a set of like targets; each target becomes a series with `job`
  and `instance` labels attached automatically.
- `static_configs` hard-codes targets — fine for a compose stack. In a cluster
  you'd swap in **service discovery** (`kubernetes_sd_configs`, etc.) so targets
  appear/vanish with pods (the `prometheus.io/scrape` annotations in
  `deploy/k8s.yaml` are exactly what that SD reads).

## Relabeling — the sharp tool
- `relabel_configs` rewrite a target's labels *before* scraping (pick/drop
  targets, build the address); `metric_relabel_configs` rewrite *after* (drop
  noisy series). Powerful and where most "why isn't it scraping?" lives.

## `up` and staleness
- Every scrape writes a synthetic `up` sample: `1` success, `0` failure. The
  first dashboard panel / alert you ever want is `up == 0`.
- Miss a few scrapes and a series goes **stale** — it stops returning in instant
  queries instead of flatlining a wrong value. Keeps `rate()` honest.

## What clicked
- The exposition endpoint and the scrape config are two ends of one contract:
  the app just keeps a `/metrics` page current, and *pull* + `up` turn "can I
  even reach it?" into a first-class metric — something a push pipeline can't
  give you for free.
