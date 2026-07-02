package main

import (
	"encoding/json"
	"strings"
	"testing"
	"time"
)

// testSpan builds a finished span with fixed times so payload assertions are
// deterministic (validTrace/validSpan come from trace_test.go).
func testSpan(failed bool, parentID string) span {
	return span{
		sc:       spanContext{traceID: validTrace, spanID: validSpan, sampled: true},
		parentID: parentID,
		name:     "POST /count",
		kind:     spanKindServer,
		start:    time.Unix(0, 1000),
		end:      time.Unix(0, 5000),
		failed:   failed,
	}
}

func TestPayloadForMapsSpanFields(t *testing.T) {
	e := newOTLPExporter("http://collector:4318/", "wordcount")
	attrs := []kv{{"http.route", "/count"}, {"http.response.status_code", "200"}}
	p := e.payloadFor(testSpan(false, validSpan), attrs)

	if len(p.ResourceSpans) != 1 {
		t.Fatalf("resourceSpans = %d, want 1", len(p.ResourceSpans))
	}
	rs := p.ResourceSpans[0]
	if got := rs.Resource.Attributes[0]; got.Key != "service.name" || got.Value.StringValue != "wordcount" {
		t.Errorf("resource attr = %+v, want service.name=wordcount", got)
	}
	if len(rs.ScopeSpans) != 1 || len(rs.ScopeSpans[0].Spans) != 1 {
		t.Fatalf("want exactly one scopeSpan with one span")
	}
	sp := rs.ScopeSpans[0].Spans[0]
	if sp.TraceID != validTrace || sp.SpanID != validSpan {
		t.Errorf("ids = %s/%s, want %s/%s", sp.TraceID, sp.SpanID, validTrace, validSpan)
	}
	if sp.ParentSpanID != validSpan {
		t.Errorf("parentSpanId = %q, want %q", sp.ParentSpanID, validSpan)
	}
	if sp.Kind != spanKindServer {
		t.Errorf("kind = %d, want %d (SERVER)", sp.Kind, spanKindServer)
	}
	if sp.StartTimeUnixNano != "1000" || sp.EndTimeUnixNano != "5000" {
		t.Errorf("times = %s/%s, want 1000/5000", sp.StartTimeUnixNano, sp.EndTimeUnixNano)
	}
	if sp.Status.Code != statusCodeUnset {
		t.Errorf("status = %d, want unset (%d) on success", sp.Status.Code, statusCodeUnset)
	}
	if len(sp.Attributes) != len(attrs) {
		t.Errorf("attributes = %d, want %d", len(sp.Attributes), len(attrs))
	}
}

func TestPayloadForFailedSpanIsError(t *testing.T) {
	e := newOTLPExporter("http://c:4318", "wordcount")
	p := e.payloadFor(testSpan(true, validSpan), nil)
	if got := p.ResourceSpans[0].ScopeSpans[0].Spans[0].Status.Code; got != statusCodeError {
		t.Errorf("status = %d, want error (%d)", got, statusCodeError)
	}
}

// A root span (no parent) must omit parentSpanId entirely — a present-but-empty
// parent reads as a dangling reference to the backend, not "no parent".
func TestRootSpanOmitsParent(t *testing.T) {
	e := newOTLPExporter("http://c:4318", "wordcount")
	body, err := json.Marshal(e.payloadFor(testSpan(false, ""), nil))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(body), "parentSpanId") {
		t.Errorf("root span JSON should omit parentSpanId:\n%s", body)
	}
}

func TestNanoIsStringNanoseconds(t *testing.T) {
	if got := nano(time.Unix(1, 500)); got != "1000000500" {
		t.Errorf("nano = %s, want 1000000500", got)
	}
}

func TestNewOTLPExporterTrimsTrailingSlash(t *testing.T) {
	e := newOTLPExporter("http://collector:4318/", "wordcount")
	if e.endpoint != "http://collector:4318" {
		t.Errorf("endpoint = %q, want trailing slash trimmed", e.endpoint)
	}
}
