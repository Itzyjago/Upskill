# Alertmanager — from a firing rule to an actual page

[alerting.md](alerting.md) ended on the split: Prometheus *decides* an alert
fires, Alertmanager *delivers* it. This is the delivery half — wrote it adding
Alertmanager to the local stack so a firing rule POSTs to a real webhook instead
of just lighting up the Prometheus UI (`deploy/observability/alertmanager.yml`).

## The handoff
- Prometheus evaluates rules, and for each *firing* alert it **pushes** the alert
  to Alertmanager's API (`/api/v2/alerts`). You wire this in `prometheus.yml`:
  ```yaml
  alerting:
    alertmanagers:
      - static_configs:
          - targets: ["alertmanager:9093"]
  ```
- Alertmanager then owns everything *after* "this is firing": grouping, dedup,
  silencing, throttling, and routing to a receiver.

## The routing tree
`route` is a tree. An alert enters at the root and walks down to the most
specific matching child; that node's `receiver` gets it.
```yaml
route:
  receiver: default            # fallback if nothing more specific matches
  group_by: [alertname]        # collapse alerts sharing these labels into one notification
  group_wait: 30s              # wait a beat to batch the first alert in a group
  group_interval: 5m           # how long before adding new alerts to an existing group's notification
  repeat_interval: 4h          # re-send a still-firing alert this often (anti-spam)
  routes:
    - matchers: [severity="page"]
      receiver: oncall-webhook # pages go here; everything else falls to default
```
- **Matching by label** is why `alerting.md` insisted labels are for machines:
  `severity: page` on the *rule* is what steers it to the oncall receiver here.

## Grouping, inhibition, silences — the noise controls
- **Grouping** — one notification for "10 instances down" instead of 10 pages.
  `group_by` picks the dimensions; everything else gets bundled.
- **Inhibition** — suppress alert B while alert A fires (don't page "high latency"
  when "whole datacenter down" is already firing — the cause already paged).
- **Silences** — a time-boxed mute by label matcher, set in the UI during a known
  maintenance window. Temporary, unlike inhibition's standing rule.

## Receivers
The leaf of the tree — *where* it goes: `webhook_configs`, `email_configs`,
`slack_configs`, `pagerduty_configs`, .... The webhook is the lowest common
denominator: Alertmanager `POST`s a JSON envelope (`status`, `alerts[]`, their
labels + annotations) to a URL, and anything that speaks HTTP can receive it.

## What clicked
- The `for:` window in [alerting.md](alerting.md) is *Prometheus'* anti-flap; the
  `group_*`/`repeat_interval` dials are *Alertmanager's* anti-spam. Two different
  layers fighting noise — one decides *whether* it's real, the other decides
  *how often a human hears about it*. Confusing them is how you get either missed
  incidents or pager fatigue.
