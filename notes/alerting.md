# Alerting — from a graph to a page

A dashboard only helps when someone's looking at it. Alerting is how the metrics
wake someone up. Wrote this adding alert rules to the local Prometheus stack
(`deploy/observability/rules/`).

## Where rules live
- Prometheus evaluates **alerting rules** on a timer (`evaluation_interval`).
  Each rule is a PromQL expression; when it returns any series, that series is a
  **firing** alert.
- Prometheus only *decides* alerts fire. Routing, grouping, dedup, and actually
  notifying (email, Slack, PagerDuty) is **Alertmanager**, a separate process.
  Clean split: evaluation vs. delivery.

## Anatomy of a rule
```yaml
- alert: HighErrorRate
  expr: |
    sum(rate(http_requests_total{status=~"5.."}[5m]))
      / sum(rate(http_requests_total[5m])) > 0.05
  for: 2m                       # must stay true 2m before firing — kills flapping
  labels:    { severity: page } # labels route it in Alertmanager
  annotations:                  # annotations are the human-readable payload
    summary: "5xx ratio above 5%"
```
- **`for`** is the anti-flap dial: the condition has to hold for the whole window
  before it pages. A one-scrape blip stays `pending` and never fires.
- **labels** are for *machines* (routing/severity); **annotations** are for
  *humans* (what's wrong, what to do). Don't mix them up.

## Symptom, not cause — and the golden signals
- Alert on what the **user feels** (error ratio, p95 latency, no traffic), not on
  causes (CPU 90%). High CPU might be fine; a slow user request never is. This is
  the golden-signals payoff from [golden-signals.md](golden-signals.md): they're
  exactly the symptom-level signals worth paging on.
- Every alert should be **actionable**. If nobody does anything when it fires,
  it's noise — and noise trains people to ignore the pager.

## What clicked
- `for:` is the whole difference between a useful alert and pager fatigue. And
  the rule expression is *just PromQL with a threshold* — the same
  `rate()`/`histogram_quantile()` from [promql.md](promql.md) I already graph,
  now with a `> 0.05` on the end. The dashboard and the alert ask the same
  question; the alert just doesn't need anyone watching.
