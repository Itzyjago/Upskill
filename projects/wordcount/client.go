package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"
)

// upstreamClient forwards /count requests to another wordcount instance,
// wrapping the outbound call in a client span so the two services' spans
// stitch into one trace — roadmap #12, notes/distributed-tracing.md
// "server spans vs. client spans". tr may be nil (export disabled), matching
// the server-side middleware in middleware.go.
type upstreamClient struct {
	url    string // full URL of the upstream /count endpoint
	client *http.Client
	tr     *otlpExporter
}

func newUpstreamClient(url string, tr *otlpExporter) *upstreamClient {
	return &upstreamClient{url: url, client: &http.Client{Timeout: 5 * time.Second}, tr: tr}
}

// count POSTs body to the upstream instance and returns its tally. ctx should
// carry the caller's current span (see withSpan in trace.go) so the outbound
// call becomes a child of it; with no span in ctx it starts a fresh trace
// rather than failing the request, the same best-effort rule trace.go follows
// for a malformed inbound traceparent.
func (u *upstreamClient) count(ctx context.Context, body []byte) (counts, error) {
	parent, ok := spanFrom(ctx)
	if !ok {
		parent = newSpanContext()
	}
	s := startClientSpan(parent, "POST /count (upstream)", time.Now())

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.url, bytes.NewReader(body))
	if err != nil {
		return counts{}, err
	}
	// Inject *this* span's id, not the parent's — it becomes the upstream
	// service's parent span id, one level deeper in the waterfall.
	req.Header.Set("traceparent", s.sc.traceparent())

	resp, err := u.client.Do(req)
	var c counts
	failed := err != nil
	if err == nil {
		defer func() { _ = resp.Body.Close() }()
		if resp.StatusCode/100 != 2 {
			err = fmt.Errorf("upstream returned %s", resp.Status)
			failed = true
		} else if decErr := json.NewDecoder(resp.Body).Decode(&c); decErr != nil {
			err = decErr
			failed = true
		}
	}
	s = s.finish(time.Now(), failed)
	u.export(s)

	if err != nil {
		return counts{}, fmt.Errorf("upstream count: %w", err)
	}
	return c, nil
}

// export ships the finished client span, off the hot path — best-effort, like
// the server-side export in middleware.go. No-op when tr is nil.
func (u *upstreamClient) export(s span) {
	if u.tr == nil {
		return
	}
	attrs := []kv{
		{"http.request.method", http.MethodPost},
		{"url.full", u.url},
	}
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		if err := u.tr.export(ctx, s, attrs); err != nil {
			slog.Debug("otlp export failed", "err", err.Error(), "trace_id", s.sc.traceID)
		}
	}()
}

// forwardCountHandler returns a /count handler that forwards the request body
// to an upstream wordcount instance instead of counting it locally — the
// "edge" side of the two-service trace.
func forwardCountHandler(up *upstreamClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "POST the text to count", http.StatusMethodNotAllowed)
			return
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		c, err := up.count(r.Context(), body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(c); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
