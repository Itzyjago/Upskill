package main

import (
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

// instrument wraps a handler so every request updates the registry: bump the
// in-flight gauge for the duration, then on completion record the labeled
// counter (by status) and the latency histogram. path is passed in rather than
// read from the URL so we label by route, not by high-cardinality raw paths.
func (m *metrics) instrument(path string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m.incInFlight()
		defer m.decInFlight()

		// Trace context: continue the inbound trace if there's a valid
		// traceparent, else start a fresh root. Either way we run our own child
		// span (see notes/trace-context.md). The span rides in the request
		// context so handlers can read it.
		var sc spanContext
		if parent, ok := parseTraceparent(r.Header.Get("traceparent")); ok {
			sc = parent.child()
		} else {
			sc = newSpanContext()
		}
		r = r.WithContext(withSpan(r.Context(), sc))
		// Echo our span back so a caller (curl -i) can see the trace id. With a
		// real downstream call you'd inject this into the *outbound* request
		// instead, making our span the next hop's parent.
		w.Header().Set("traceparent", sc.traceparent())

		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next(rec, r)

		dur := time.Since(start)
		m.observe(r.Method, path, strconv.Itoa(rec.status), dur.Seconds())
		// One structured line per request — the same dimensions the metrics use,
		// plus trace_id/span_id so a log and its trace cross-link
		// (see notes/structured-logging.md, notes/trace-context.md).
		slog.Info("request",
			"method", r.Method,
			"path", path,
			"status", rec.status,
			"dur_ms", dur.Milliseconds(),
			"trace_id", sc.traceID,
			"span_id", sc.spanID,
		)
	}
}
