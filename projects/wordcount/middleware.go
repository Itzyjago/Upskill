package main

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

// statusRecorder wraps a ResponseWriter to remember the status code the handler
// wrote — net/http gives no way to read it back after the fact. Defaults to 200
// because a handler that only calls Write (never WriteHeader) implies 200 OK.
type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

// instrument wraps a handler so every request updates the registry, logs one
// structured line, and emits a trace span. It bumps the in-flight gauge for the
// duration, then on completion records the labeled counter (by status) and the
// latency histogram. path is passed in rather than read from the URL so we label
// by route, not by high-cardinality raw paths. tr may be nil (export disabled).
func (m *metrics) instrument(path string, tr *otlpExporter, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m.incInFlight()
		defer m.decInFlight()

		// Start the timed server span: continue the inbound trace if there's a
		// valid traceparent (the sender becomes our parent), else start a fresh
		// root. The span's ids ride in the request context so handlers can read
		// them; the timing rides with the span for export (trace.go, notes/otlp.md).
		s := startServerSpan(r.Header.Get("traceparent"), r.Method+" "+path, time.Now())
		r = r.WithContext(withSpan(r.Context(), s.sc))
		// Echo our span back so a caller (curl -i) can see the trace id. With a
		// real downstream call you'd inject this into the *outbound* request
		// instead, making our span the next hop's parent.
		w.Header().Set("traceparent", s.sc.traceparent())

		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next(rec, r)

		s = s.finish(time.Now(), rec.status >= 500)
		dur := s.duration()
		m.observe(r.Method, path, strconv.Itoa(rec.status), dur.Seconds())
		// One structured line per request — the same dimensions the metrics use,
		// plus trace_id/span_id so a log and its trace cross-link
		// (see notes/structured-logging.md, notes/trace-context.md).
		slog.Info("request",
			"method", r.Method,
			"path", path,
			"status", rec.status,
			"dur_ms", dur.Milliseconds(),
			"trace_id", s.sc.traceID,
			"span_id", s.sc.spanID,
		)

		// Export the finished span out-of-band: never block the response on the
		// collector, and a collector that's down must not fail the request
		// (best-effort, like trace propagation). nil tr → export disabled.
		if tr != nil {
			attrs := []kv{
				{"http.request.method", r.Method},
				{"http.route", path},
				{"http.response.status_code", strconv.Itoa(rec.status)},
			}
			go func() {
				ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
				defer cancel()
				if err := tr.export(ctx, s, attrs); err != nil {
					slog.Debug("otlp export failed", "err", err.Error(), "trace_id", s.sc.traceID)
				}
			}()
		}
	}
}
