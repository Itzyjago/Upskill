package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// newMux wires the HTTP routes. Split out so tests can exercise the handlers
// without binding a port. Takes the metrics registry so /metrics can expose it,
// the span exporter (nil = export disabled) so handlers emit traces, and an
// upstream client (nil = count locally, set = forward to it — roadmap #12,
// client.go) so /count can act as either the edge or the leaf of a trace.
func newMux(m *metrics, tr *otlpExporter, up *upstreamClient) *http.ServeMux {
	handler := countHandler
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
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok\n"))
}

// countHandler counts the request body and returns the tally as JSON.
func countHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "POST the text to count", http.StatusMethodNotAllowed)
		return
	}
	c, err := count(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(c); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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

	srv := &http.Server{
		Addr:              addr,
		Handler:           newMux(m, tr, up),
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
