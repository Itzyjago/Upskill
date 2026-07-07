package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// maxCountBodyBytes caps how much of a /count request body count() or
// forwardCountHandler will read — an unbounded read is a resource-exhaustion
// DoS, not just a style nit (notes/security.md). Generous for a text-counting
// tool, small enough that one request can't OOM the process.
const maxCountBodyBytes = 10 << 20 // 10MB

// statusForBodyErr maps a body-read error to an HTTP status: 413 if
// http.MaxBytesReader's cap tripped, 400 for any other read failure.
func statusForBodyErr(err error) int {
	var mbErr *http.MaxBytesError
	if errors.As(err, &mbErr) {
		return http.StatusRequestEntityTooLarge
	}
	return http.StatusBadRequest
}

// newMux wires the HTTP routes. Split out so tests can exercise the handlers
// without binding a port. Takes the metrics registry so /metrics can expose it,
// the span exporter (nil = export disabled) so handlers emit traces, an
// upstream client (nil = count locally, set = forward to it — roadmap #12,
// client.go) so /count can act as either the edge or the leaf of a trace, and
// an idempotency store (nil = caching disabled) so a client-retried /count
// replays the first response instead of counting twice (idempotency.go).
func newMux(m *metrics, tr *otlpExporter, up *upstreamClient, store *idempotencyStore) *http.ServeMux {
	handler := countHandlerFunc(store)
	if up != nil {
		handler = forwardCountHandler(up)
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", m.instrument("/healthz", tr, healthHandler))
	mux.HandleFunc("/count", m.instrument("/count", tr, handler))
	// /metrics is intentionally *not* instrumented — a scraper hitting it every
	// few seconds would swamp the very numbers it's collecting.
	mux.HandleFunc("/metrics", m.metricsHandler)
	return mux
}

// healthHandler is the readiness/liveness probe target — cheap and dependency-free.
func healthHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok\n"))
}

// idempotencyKeyHeader names the client-supplied header for a retry-safe
// request — the header name Stripe's API popularized for this pattern.
const idempotencyKeyHeader = "Idempotency-Key"

// countHandler is the /count handler with idempotency caching disabled —
// the same behavior as before idempotency.go existed. Kept as a package var
// (rather than folding callers over to countHandlerFunc(nil)) so existing
// direct calls in tests didn't need touching.
var countHandler = countHandlerFunc(nil)

// countHandlerFunc builds the /count handler. store may be nil (caching
// off); when set and the client sends an Idempotency-Key header, a repeat
// request with that key and the *same* body replays the cached response
// instead of counting again — a key reused with a *different* body is a
// 409, not a silent cache hit (idempotency.go's lookup()).
func countHandlerFunc(store *idempotencyStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "POST the text to count", http.StatusMethodNotAllowed)
			return
		}

		body, err := io.ReadAll(http.MaxBytesReader(w, r.Body, maxCountBodyBytes))
		if err != nil {
			http.Error(w, err.Error(), statusForBodyErr(err))
			return
		}

		key := r.Header.Get(idempotencyKeyHeader)
		var bodyHash string
		if store != nil && key != "" {
			bodyHash = hashBody(body)
			if status, cached, ok, lookupErr := store.lookup(key, bodyHash); lookupErr != nil {
				http.Error(w, lookupErr.Error(), http.StatusConflict)
				return
			} else if ok {
				w.Header().Set("Content-Type", "application/json")
				w.Header().Set("Idempotency-Replayed", "true")
				w.WriteHeader(status)
				_, _ = w.Write(cached)
				return
			}
		}

		c, err := count(bytes.NewReader(body))
		if err != nil {
			http.Error(w, err.Error(), statusForBodyErr(err))
			return
		}

		var buf bytes.Buffer
		if err := json.NewEncoder(&buf).Encode(c); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if store != nil && key != "" {
			store.store(key, bodyHash, http.StatusOK, buf.Bytes())
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(buf.Bytes())
	}
}

// serve runs the HTTP service until SIGINT/SIGTERM, then drains in-flight
// requests with a bounded grace period (the readiness-probe lesson).
func serve(addr string) error {
	// JSON logs to stdout — structured and machine-parseable, the standard for
	// containerized services where stdout is the log stream.
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
	m := newMetrics()

	// service.name is also a standard OTel env var — set it per instance
	// (compose gives the edge and the upstream different values) so Jaeger's
	// service dropdown actually shows two services, not one instance twice.
	service := "wordcount"
	if s := os.Getenv("OTEL_SERVICE_NAME"); s != "" {
		service = s
	}

	// Span export is opt-in via the standard OTel env var. Unset → tr stays nil
	// and we skip export; spans are still timed and logged. See notes/otlp.md.
	var tr *otlpExporter
	if ep := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"); ep != "" {
		tr = newOTLPExporter(ep, service)
		slog.Info("otlp export enabled", "endpoint", ep, "service", service)
	}

	// Set on the "edge" instance only — it forwards /count to another
	// wordcount instance instead of counting locally, so the trace stitches
	// across two services (roadmap #12, client.go).
	var up *upstreamClient
	if url := os.Getenv("WORDCOUNT_UPSTREAM_URL"); url != "" {
		up = newUpstreamClient(url, tr)
		slog.Info("forwarding /count upstream", "url", url)
	}

	store := newIdempotencyStore()

	srv := &http.Server{
		Addr:              addr,
		Handler:           newMux(m, tr, up, store),
		ReadHeaderTimeout: 5 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	errCh := make(chan error, 1)
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		// Stop accepting new connections, give the rest 10s to finish.
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return srv.Shutdown(shutdownCtx)
	}
}
