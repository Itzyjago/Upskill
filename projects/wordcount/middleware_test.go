package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// fakeCollector is a real OTLP/HTTP collector double: an httptest.Server that
// records the payload of the first POST it receives and signals a channel, so
// a test can deterministically wait on middleware.go's fire-and-forget export
// goroutine instead of sleeping and racing it (notes/testing.md, "the
// fake-collector-shaped hole").
type fakeCollector struct {
	*httptest.Server
	received chan otlpPayload
}

func newFakeCollector() *fakeCollector {
	fc := &fakeCollector{received: make(chan otlpPayload, 1)}
	fc.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var p otlpPayload
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		fc.received <- p
		w.WriteHeader(http.StatusOK)
	}))
	return fc
}

// waitForSpan blocks for the exported payload or fails after a bound — the
// export goroutine in instrument() is never joined, so this is the only way
// to observe it without an arbitrary sleep racing the goroutine.
func (fc *fakeCollector) waitForSpan(t *testing.T) otlpPayload {
	t.Helper()
	select {
	case p := <-fc.received:
		return p
	case <-time.After(2 * time.Second):
		t.Fatal("collector never received an exported span")
		return otlpPayload{}
	}
}

// TestInstrumentExportsSpanToCollector closes the gap testing.md calls out:
// every other test passes a nil *otlpExporter and skips export entirely, so a
// payload-shape or wiring bug in the goroutine at the bottom of instrument()
// could ship untested. Here tr points at a real fake, so a passing test means
// export actually fired with the right ids and attributes, not just that the
// handler ran.
func TestInstrumentExportsSpanToCollector(t *testing.T) {
	fc := newFakeCollector()
	defer fc.Close()

	tr := newOTLPExporter(fc.URL, "wordcount")
	mux := newMux(newMetrics(), tr, nil, nil)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/count", strings.NewReader("hello world\n"))
	mux.ServeHTTP(rec, req)

	p := fc.waitForSpan(t)
	if len(p.ResourceSpans) != 1 {
		t.Fatalf("resourceSpans = %d, want 1", len(p.ResourceSpans))
	}
	rs := p.ResourceSpans[0]
	if got := rs.Resource.Attributes[0].Value.StringValue; got != "wordcount" {
		t.Errorf("service.name = %q, want wordcount", got)
	}
	if len(rs.ScopeSpans) != 1 || len(rs.ScopeSpans[0].Spans) != 1 {
		t.Fatalf("want exactly one scopeSpan with one span")
	}
	sp := rs.ScopeSpans[0].Spans[0]
	if sp.TraceID == "" || sp.SpanID == "" {
		t.Errorf("exported span missing ids: %+v", sp)
	}

	var gotRoute, gotStatus string
	for _, a := range sp.Attributes {
		switch a.Key {
		case "http.route":
			gotRoute = a.Value.StringValue
		case "http.response.status_code":
			gotStatus = a.Value.StringValue
		}
	}
	if gotRoute != "/count" {
		t.Errorf("http.route attr = %q, want /count", gotRoute)
	}
	if gotStatus != "200" {
		t.Errorf("http.response.status_code attr = %q, want 200", gotStatus)
	}
}

// TestInstrumentSkipsExportWhenDisabled pins the other half of the contract:
// a nil tr (export disabled) must never dial out, even though every field
// needed to build a payload is otherwise available.
func TestInstrumentSkipsExportWhenDisabled(t *testing.T) {
	fc := newFakeCollector()
	defer fc.Close()

	tr := newOTLPExporter(fc.URL, "wordcount")
	mux := newMux(newMetrics(), tr, nil, nil)

	// Confirm the collector is reachable and instrumented requests do export
	// (same assertion as the test above, kept minimal) before trusting silence
	// from the nil-tr mux below as "disabled" rather than "collector unreachable".
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/count", strings.NewReader("warmup\n"))
	mux.ServeHTTP(rec, req)
	fc.waitForSpan(t) // drain the warmup export

	disabled := newMux(newMetrics(), nil, nil, nil)
	rec2 := httptest.NewRecorder()
	req2 := httptest.NewRequest(http.MethodPost, "/count", strings.NewReader("hi\n"))
	disabled.ServeHTTP(rec2, req2)

	select {
	case p := <-fc.received:
		t.Fatalf("collector received a span with nil tr: %+v", p)
	case <-time.After(200 * time.Millisecond):
		// expected: nothing arrives
	}
}
