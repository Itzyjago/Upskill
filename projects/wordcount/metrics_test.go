package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestObserveCountsAndBuckets(t *testing.T) {
	m := newMetrics()
	m.observe("POST", "/count", "200", 0.002) // 2ms
	m.observe("POST", "/count", "200", 0.2)   // 200ms

	route := labelKey{method: "POST", path: "/count"}
	if got := m.cnt[route]; got != 2 {
		t.Fatalf("count = %d, want 2", got)
	}
	// Buckets are cumulative. durBounds index 3 is le=0.005s, which should hold
	// only the 2ms sample; index 7 is le=0.5s, which should hold both.
	if got := m.buckets[route][3]; got != 1 {
		t.Errorf("le=0.005 bucket = %d, want 1", got)
	}
	if got := m.buckets[route][7]; got != 2 {
		t.Errorf("le=0.5 bucket = %d, want 2", got)
	}
	if c := m.reqTotal[labelKey{"POST", "/count", "200"}]; c != 2 {
		t.Errorf("requests_total = %d, want 2", c)
	}
}

func TestRenderExposition(t *testing.T) {
	m := newMetrics()
	m.observe("POST", "/count", "200", 0.003)
	out := m.render()

	want := []string{
		"# TYPE http_requests_total counter",
		`http_requests_total{method="POST",path="/count",status="200"} 1`,
		"# TYPE http_request_duration_seconds histogram",
		`http_request_duration_seconds_count{method="POST",path="/count"} 1`,
		`http_request_duration_seconds_bucket{method="POST",path="/count",le="+Inf"} 1`,
		"# TYPE http_requests_in_flight gauge",
	}
	for _, s := range want {
		if !strings.Contains(out, s) {
			t.Errorf("exposition missing %q\n--- got ---\n%s", s, out)
		}
	}
}

func TestInstrumentedEndpoint(t *testing.T) {
	m := newMetrics()
	mux := newMux(m, nil) // nil exporter: span export disabled for this test

	// Drive one /count request through the instrumented mux.
	req := httptest.NewRequest(http.MethodPost, "/count", strings.NewReader("hello world\n"))
	mux.ServeHTTP(httptest.NewRecorder(), req)

	// Scrape /metrics and confirm the request was recorded.
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/metrics", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); !strings.HasPrefix(ct, "text/plain") {
		t.Errorf("content-type = %q, want text/plain prefix", ct)
	}
	if body := rec.Body.String(); !strings.Contains(body, `path="/count",status="200"`) {
		t.Errorf("metrics missing the /count series:\n%s", body)
	}
}
