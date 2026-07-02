# Grafana — turning metrics into dashboards

Prometheus stores and queries the time series (see [prometheus.md](prometheus.md)
and [promql.md](promql.md)); Grafana is the front end that draws them. Wrote
this while standing up a local stack to scrape wordcount's `/metrics`.

## The pieces
- **Data source** — where the numbers come from. Point Grafana at the Prometheus
  HTTP API (`http://prometheus:9090`) and it runs PromQL on your behalf.
- **Dashboard** — a grid of **panels**. Each panel is one (or more) PromQL
  queries plus a visualization (time series, stat, gauge, table, heatmap).
- **Panel target** — the query. The same `rate(...)`/`histogram_quantile(...)`
  expressions from the PromQL notes, just typed into a panel instead of the
  Prometheus expression browser.

## Provisioning — dashboards as code
Clicking dashboards together by hand doesn't survive a container restart. The
fix is **provisioning**: drop YAML + JSON files in and Grafana loads them on boot.
- `provisioning/datasources/*.yml` — declare the Prometheus data source (so you
  never click "Add data source").
- `provisioning/dashboards/*.yml` — a *provider* that tells Grafana which folder
  to watch for dashboard JSON.
- `dashboards/*.json` — the dashboards themselves, exported from the UI or
  hand-written. Version-controlled, reviewable, reproducible.
- The key glue: a panel's `datasource` references the data source by **uid**, so
  pin a known `uid` in the datasource YAML and use it in the JSON.

## Variables (templating)
- A dashboard **variable** like `$path` becomes a dropdown; queries interpolate
  it (`...{path="$path"}`). One dashboard, filterable across routes/instances,
  instead of one per target.
- `label_values(http_requests_total, path)` populates the dropdown straight from
  the data — the menu can't drift from reality.

## Grafana alerting vs. Prometheus rules + Alertmanager (roadmap #14)
wordcount already has one alerting path: Prometheus rules decide *whether*
(`rules/alerts.yml`, the `for:` window), Alertmanager decides *how a human
hears about it* (`alertmanager.yml`, the routing tree — see
[alertmanager.md](alertmanager.md)). Grafana ships its own, separate alerting
engine that can do the same job end to end. Same three questions, different
pieces:
- **Whether it's firing** — a Grafana **alert rule** is a query (any data
  source, not just Prometheus — the PromQL from `promql.md` works unmodified)
  plus a threshold and a `for:`-equivalent evaluation window. Functionally the
  same decision as a Prometheus rule; just evaluated *by Grafana* instead of
  *by Prometheus*, so it works even against a data source with no rule engine
  of its own (a plain SQL database, say).
- **Where it goes** — a **contact point** is Grafana's receiver: webhook,
  email, Slack, PagerDuty. Direct analog of Alertmanager's `receivers`.
- **Routing + noise control** — a **notification policy** tree matches on
  labels and routes to a contact point, with its own grouping and
  `repeat_interval`. Direct analog of Alertmanager's `route` tree — same
  shape, same job, different config format.
- The real difference isn't the feature set, it's *where the state lives*.
  Prometheus + Alertmanager is data-source-native: the alert is defined right
  next to the metric, versioned with the same rule files. Grafana alerting is
  data-source-agnostic: one alerting engine over everything Grafana can query,
  which is the whole pitch if the stack isn't all-Prometheus. For a
  single-Prometheus setup like this one, running *both* is redundant — pick a
  side so the same "5xx ratio > 5%" condition doesn't quietly drift into two
  different thresholds nobody remembers to keep in sync.
- wordcount keeps Prometheus + Alertmanager as the source of truth (already
  built, already provisioned) and adds a **second**, Grafana-native copy of
  the same `HighErrorRate` alert purely to see the provisioning shape
  side by side — not because a real setup should run both.

## What clicked
- Grafana doesn't store metrics — it's a *query-and-draw* layer over Prometheus.
  All the intelligence is still PromQL; Grafana just arranges the answers and
  refreshes them on a timer. And provisioning is the same lesson as the rest of
  this repo: if it isn't a file in git, it isn't real — it's a snowflake waiting
  to be lost on the next restart.
- Grafana alerting and Prometheus/Alertmanager solve the identical problem with
  near-identical shapes (rule → route → receiver) — which makes sense once you
  see routing trees and grouping/`repeat_interval` aren't a Prometheus-specific
  idea, they're just what *any* alerting system needs to avoid paging a human
  once per firing series.
