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

## What clicked
- Grafana doesn't store metrics — it's a *query-and-draw* layer over Prometheus.
  All the intelligence is still PromQL; Grafana just arranges the answers and
  refreshes them on a timer. And provisioning is the same lesson as the rest of
  this repo: if it isn't a file in git, it isn't real — it's a snowflake waiting
  to be lost on the next restart.
