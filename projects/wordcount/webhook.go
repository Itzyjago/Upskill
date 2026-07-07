package main

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"os"
	"time"
)

// alertNotification is the slice of Alertmanager's webhook payload we care about:
// the envelope it POSTs to a webhook receiver, with one entry per alert plus its
// labels and annotations (notes/alertmanager.md).
type alertNotification struct {
	Status string `json:"status"`
	Alerts []struct {
		Status      string            `json:"status"`
		Labels      map[string]string `json:"labels"`
		Annotations map[string]string `json:"annotations"`
	} `json:"alerts"`
}

// alertWebhookHandler logs the alerts Alertmanager POSTs to it. Split out from
// webhookSink, same reason newMux is split from main.go's serve mode: so a
// test can exercise it directly (httptest.NewRecorder) without binding a port.
func alertWebhookHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "POST Alertmanager webhook payloads here", http.StatusMethodNotAllowed)
		return
	}
	// Same cap as countHandler/forwardCountHandler (notes/security.md) — this
	// endpoint faces Alertmanager, not the public internet, but an unbounded
	// read is unbounded regardless of who's supposed to be on the other end.
	// Read fully *before* decoding rather than wiring json.NewDecoder straight
	// onto MaxBytesReader: whether a mid-decode cap trip still satisfies
	// errors.As(err, *http.MaxBytesError) depends on exactly where in the
	// JSON it trips (verified — it doesn't always), so statusForBodyErr needs
	// the read's own error, not whatever the decoder reshapes it into.
	body, err := io.ReadAll(http.MaxBytesReader(w, r.Body, maxCountBodyBytes))
	if err != nil {
		http.Error(w, err.Error(), statusForBodyErr(err))
		return
	}
	var n alertNotification
	if err := json.Unmarshal(body, &n); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	for _, a := range n.Alerts {
		// status flips to "resolved" when the alert clears (send_resolved).
		slog.Info("alert received",
			"group_status", n.Status,
			"status", a.Status,
			"alertname", a.Labels["alertname"],
			"severity", a.Labels["severity"],
			"summary", a.Annotations["summary"],
		)
	}
	w.WriteHeader(http.StatusOK)
}

// newWebhookMux wires the webhook sink's routes. Split out from webhookSink,
// same reason newMux is split from server.go's serve() — so a test can drive
// the real mux (a real HTTP round trip, not just calling a handler func
// directly) without binding webhookSink's own fixed addr.
func newWebhookMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", healthHandler) // reuse the serve-mode probe handler
	mux.HandleFunc("/", alertWebhookHandler)
	return mux
}

// webhookSink runs a tiny HTTP receiver that logs the alerts Alertmanager POSTs
// to it — the "route an alert somewhere real" half of roadmap #10, reusing this
// same binary instead of standing up another service. It's not a pager: it's
// proof that a firing rule made it all the way through routing to a receiver,
// and the structured line is something downstream tooling could act on.
func webhookSink(addr string) error {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	srv := &http.Server{
		Addr:              addr,
		Handler:           newWebhookMux(),
		ReadHeaderTimeout: 5 * time.Second,
	}
	slog.Info("webhook sink listening", "addr", addr)
	return srv.ListenAndServe()
}
