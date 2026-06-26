package main

import (
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

		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next(rec, r)

		m.observe(r.Method, path, strconv.Itoa(rec.status), time.Since(start).Seconds())
	}
}
