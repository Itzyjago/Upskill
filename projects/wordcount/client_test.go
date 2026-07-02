package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestUpstreamClientInjectsChildTraceparent(t *testing.T) {
	var gotHeader string
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotHeader = r.Header.Get("traceparent")
		_ = json.NewEncoder(w).Encode(counts{Lines: 1, Words: 2, Bytes: 12})
	}))
	defer upstream.Close()

	u := newUpstreamClient(upstream.URL, nil)
	ctx := withSpan(context.Background(), spanContext{traceID: validTrace, spanID: validSpan, sampled: true})

	c, err := u.count(ctx, []byte("hello world\n"))
	if err != nil {
		t.Fatalf("count: %v", err)
	}
	if c.Lines != 1 || c.Words != 2 || c.Bytes != 12 {
		t.Errorf("got %+v, want the upstream's tally, not a local recount", c)
	}

	out, ok := parseTraceparent(gotHeader)
	if !ok {
		t.Fatalf("outbound traceparent %q did not parse", gotHeader)
	}
	if out.traceID != validTrace {
		t.Errorf("traceID = %q, want the same trace %q", out.traceID, validTrace)
	}
	// The injected id must be the *client* span's own id, not the parent's —
	// that's what makes this span the upstream's parent, one level deeper.
	if out.spanID == validSpan {
		t.Errorf("injected spanID = %q, want a fresh child id, not the parent's %q", out.spanID, validSpan)
	}
}

func TestUpstreamClientNoParentStartsFreshTrace(t *testing.T) {
	var gotHeader string
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotHeader = r.Header.Get("traceparent")
		_ = json.NewEncoder(w).Encode(counts{})
	}))
	defer upstream.Close()

	u := newUpstreamClient(upstream.URL, nil)
	if _, err := u.count(context.Background(), []byte("hi")); err != nil {
		t.Fatalf("count: %v", err)
	}
	if _, ok := parseTraceparent(gotHeader); !ok {
		t.Errorf("no span in ctx: outbound traceparent %q should still be valid", gotHeader)
	}
}

func TestUpstreamClientErrorsOnUpstreamFailure(t *testing.T) {
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "boom", http.StatusInternalServerError)
	}))
	defer upstream.Close()

	u := newUpstreamClient(upstream.URL, nil)
	if _, err := u.count(context.Background(), []byte("hi")); err == nil {
		t.Error("want an error when the upstream returns 500")
	}
}

func TestForwardCountHandlerReturnsUpstreamTally(t *testing.T) {
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_ = json.NewEncoder(w).Encode(counts{Lines: 0, Words: 2, Bytes: 11})
	}))
	defer upstream.Close()

	up := newUpstreamClient(upstream.URL, nil)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/count", strings.NewReader("hello world"))

	forwardCountHandler(up)(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	var c counts
	if err := json.NewDecoder(rec.Body).Decode(&c); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if c.Words != 2 {
		t.Errorf("words = %d, want the upstream's answer (2)", c.Words)
	}
}

func TestForwardCountHandlerRejectsGET(t *testing.T) {
	up := newUpstreamClient("http://unused.invalid", nil)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/count", nil)

	forwardCountHandler(up)(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusMethodNotAllowed)
	}
}

func TestForwardCountHandlerRejectsOversizedBody(t *testing.T) {
	up := newUpstreamClient("http://unused.invalid", nil)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/count", strings.NewReader(strings.Repeat("a", maxCountBodyBytes+1)))

	forwardCountHandler(up)(rec, req)

	if rec.Code != http.StatusRequestEntityTooLarge {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusRequestEntityTooLarge)
	}
}

func TestForwardCountHandlerBadGatewayOnUpstreamDown(t *testing.T) {
	up := newUpstreamClient("http://127.0.0.1:1", nil) // nothing listens here
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/count", strings.NewReader("hi"))

	forwardCountHandler(up)(rec, req)

	if rec.Code != http.StatusBadGateway {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusBadGateway)
	}
}
