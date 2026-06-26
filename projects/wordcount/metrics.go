package main

import (
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
