package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// newMux wires the HTTP routes. Split out so tests can exercise the handlers
// without binding a port. Takes the metrics registry so /metrics can expose it.
func newMux(m *metrics) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", healthHandler)
	mux.HandleFunc("/count", countHandler)
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
	m := newMetrics()
	srv := &http.Server{
		Addr:              addr,
		Handler:           newMux(m),
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
