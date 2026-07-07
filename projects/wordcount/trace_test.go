package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const (
	validTrace = "4bf92f3577b34da6a3ce929d0e0e4736" // 32 hex
	validSpan  = "00f067aa0ba902b7"                 // 16 hex
)

func TestParseTraceparentValid(t *testing.T) {
	in := "00-" + validTrace + "-" + validSpan + "-01"
	sc, ok := parseTraceparent(in)
	if !ok {
		t.Fatalf("parseTraceparent(%q) ok=false, want true", in)
	}
	if sc.traceID != validTrace {
		t.Errorf("traceID = %q, want %q", sc.traceID, validTrace)
	}
	// The parsed span-id is the sender's span — our parent-to-be.
	if sc.spanID != validSpan {
		t.Errorf("spanID = %q, want %q", sc.spanID, validSpan)
	}
	if !sc.sampled {
		t.Errorf("sampled = false, want true (flags 01)")
	}
}

func TestParseTraceparentRejects(t *testing.T) {
	cases := map[string]string{
		"empty":           "",
		"too few fields":  "00-" + validTrace + "-" + validSpan,
		"bad version":     "ff-" + validTrace + "-" + validSpan + "-01",
		"short trace":     "00-abc-" + validSpan + "-01",
		"short span":      "00-" + validTrace + "-abc-01",
		"all-zero trace":  "00-" + strings.Repeat("0", 32) + "-" + validSpan + "-01",
		"all-zero span":   "00-" + validTrace + "-" + strings.Repeat("0", 16) + "-01",
		"uppercase hex":   "00-" + strings.ToUpper(validTrace) + "-" + validSpan + "-01",
		"non-hex in span": "00-" + validTrace + "-00f067aa0ba902gz-01",
		"bad flags len":   "00-" + validTrace + "-" + validSpan + "-1",
	}
	for name, in := range cases {
		if _, ok := parseTraceparent(in); ok {
			t.Errorf("%s: parseTraceparent(%q) ok=true, want false", name, in)
		}
	}
}

func TestChildKeepsTraceNewSpan(t *testing.T) {
	parent := spanContext{traceID: validTrace, spanID: validSpan, sampled: true}
	child := parent.child()
	if child.traceID != parent.traceID {
		t.Errorf("child traceID = %q, want same as parent %q", child.traceID, parent.traceID)
	}
	if child.spanID == parent.spanID {
		t.Errorf("child spanID = %q, want a fresh id, not the parent's", child.spanID)
	}
	if len(child.spanID) != 16 {
		t.Errorf("child spanID len = %d, want 16 hex", len(child.spanID))
	}
	if !child.sampled {
		t.Errorf("child dropped the sampled flag")
	}
}

func TestTraceparentRoundTrip(t *testing.T) {
	sc := spanContext{traceID: validTrace, spanID: validSpan, sampled: true}
	got, ok := parseTraceparent(sc.traceparent())
	if !ok {
		t.Fatalf("round-trip parse of %q failed", sc.traceparent())
	}
	if got != sc {
		t.Errorf("round trip = %+v, want %+v", got, sc)
	}
}

func TestMiddlewarePropagatesTrace(t *testing.T) {
	in := "00-" + validTrace + "-" + validSpan + "-01"
	req := httptest.NewRequest(http.MethodPost, "/count", strings.NewReader("hello world\n"))
	req.Header.Set("traceparent", in)

	rec := httptest.NewRecorder()
	newMux(newMetrics(), nil, nil, nil).ServeHTTP(rec, req) // nil exporter, nil upstream

	out, ok := parseTraceparent(rec.Header().Get("traceparent"))
	if !ok {
		t.Fatalf("response traceparent %q did not parse", rec.Header().Get("traceparent"))
	}
	// Same trace, but our own (child) span — not the inbound parent's span.
	if out.traceID != validTrace {
		t.Errorf("response traceID = %q, want the inbound %q", out.traceID, validTrace)
	}
	if out.spanID == validSpan {
		t.Errorf("response spanID = %q, want a fresh child span", out.spanID)
	}
}

func TestMiddlewareStartsFreshTrace(t *testing.T) {
	// No inbound traceparent → a brand-new root trace, still a valid header out.
	req := httptest.NewRequest(http.MethodPost, "/count", strings.NewReader("hi"))
	rec := httptest.NewRecorder()
	newMux(newMetrics(), nil, nil, nil).ServeHTTP(rec, req) // nil exporter, nil upstream

	if _, ok := parseTraceparent(rec.Header().Get("traceparent")); !ok {
		t.Errorf("no inbound trace: response traceparent %q is not valid", rec.Header().Get("traceparent"))
	}
}
