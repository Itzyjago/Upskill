package main

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
)

// durBounds are the histogram bucket upper bounds, in seconds. Cumulative: an
// observation lands in every bucket whose bound it's <=. Trimmed to what a fast
// in-memory handler actually sees (sub-millisecond to ~1s).
var durBounds = []float64{0.0001, 0.0005, 0.001, 0.005, 0.01, 0.05, 0.1, 0.5, 1}

// labelKey identifies one labeled time series. For the duration histogram the
// status field is left empty — we bucket latency by route, not by outcome.
type labelKey struct {
	method string
	path   string
	status string
}

// metrics is a tiny, dependency-free metrics registry. It emits the Prometheus
// text exposition format by hand (see render) — the goal is to understand the
// format, not to pull in client_golang.
type metrics struct {
	mu       sync.Mutex
	reqTotal map[labelKey]int64   // http_requests_total, by method+path+status
	buckets  map[labelKey][]int64 // cumulative bucket counts, parallel to durBounds
	sum      map[labelKey]float64 // sum of observed seconds, by route
	cnt      map[labelKey]int64   // observation count, by route
	inFlight int64                // gauge — touched atomically, no lock
}

func newMetrics() *metrics {
	return &metrics{
		reqTotal: make(map[labelKey]int64),
		buckets:  make(map[labelKey][]int64),
		sum:      make(map[labelKey]float64),
		cnt:      make(map[labelKey]int64),
	}
}

// incInFlight / decInFlight move the in-flight gauge. Atomic so the hot path
// doesn't take the registry lock just to bump a counter.
func (m *metrics) incInFlight() { atomic.AddInt64(&m.inFlight, 1) }
func (m *metrics) decInFlight() { atomic.AddInt64(&m.inFlight, -1) }

// observe records one finished request: bumps the labeled counter and files the
// latency into the per-route cumulative histogram.
func (m *metrics) observe(method, path, status string, seconds float64) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.reqTotal[labelKey{method, path, status}]++

	route := labelKey{method: method, path: path}
	if m.buckets[route] == nil {
		m.buckets[route] = make([]int64, len(durBounds))
	}
	for i, bound := range durBounds {
		if seconds <= bound {
			m.buckets[route][i]++
		}
	}
	m.sum[route] += seconds
	m.cnt[route]++
}

// sortedKeys returns the registry keys in a stable order so the exposition
// output is deterministic — Go map iteration order isn't.
func sortedKeys(series map[labelKey]int64) []labelKey {
	keys := make([]labelKey, 0, len(series))
	for k := range series {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		if keys[i].path != keys[j].path {
			return keys[i].path < keys[j].path
		}
		if keys[i].method != keys[j].method {
			return keys[i].method < keys[j].method
		}
		return keys[i].status < keys[j].status
	})
	return keys
}

// render produces the full Prometheus text exposition for the registry.
func (m *metrics) render() string {
	m.mu.Lock()
	defer m.mu.Unlock()

	var b strings.Builder
	p := func(format string, args ...any) {
		_, _ = fmt.Fprintf(&b, format, args...)
	}

	p("# HELP http_requests_total Total HTTP requests handled.\n")
	p("# TYPE http_requests_total counter\n")
	for _, k := range sortedKeys(m.reqTotal) {
		p("http_requests_total{method=%q,path=%q,status=%q} %d\n",
			k.method, k.path, k.status, m.reqTotal[k])
	}

	p("# HELP http_request_duration_seconds Request latency by route.\n")
	p("# TYPE http_request_duration_seconds histogram\n")
	for _, k := range sortedKeys(m.cnt) {
		for i, bound := range durBounds {
			p("http_request_duration_seconds_bucket{method=%q,path=%q,le=%q} %d\n",
				k.method, k.path, strconv.FormatFloat(bound, 'g', -1, 64), m.buckets[k][i])
		}
		p("http_request_duration_seconds_bucket{method=%q,path=%q,le=\"+Inf\"} %d\n",
			k.method, k.path, m.cnt[k])
		p("http_request_duration_seconds_sum{method=%q,path=%q} %g\n", k.method, k.path, m.sum[k])
		p("http_request_duration_seconds_count{method=%q,path=%q} %d\n", k.method, k.path, m.cnt[k])
	}

	p("# HELP http_requests_in_flight Requests currently being served.\n")
	p("# TYPE http_requests_in_flight gauge\n")
	p("http_requests_in_flight %d\n", atomic.LoadInt64(&m.inFlight))

	return b.String()
}

// metricsHandler serves the registry in Prometheus text format. Kept on GET and
// dependency-light so a scraper can hit it cheaply.
func (m *metrics) metricsHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; version=0.0.4; charset=utf-8")
	_, _ = w.Write([]byte(m.render()))
}
