package main

import (
	"encoding/json"
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

// webhookSink runs a tiny HTTP receiver that logs the alerts Alertmanager POSTs
// to it — the "route an alert somewhere real" half of roadmap #10, reusing this
// same binary instead of standing up another service. It's not a pager: it's
// proof that a firing rule made it all the way through routing to a receiver,
// and the structured line is something downstream tooling could act on.
func webhookSink(addr string) error {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", healthHandler) // reuse the serve-mode probe handler
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "POST Alertmanager webhook payloads here", http.StatusMethodNotAllowed)
			return
		}
		var n alertNotification
		if err := json.NewDecoder(r.Body).Decode(&n); err != nil {
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
	})

	srv := &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}
	slog.Info("webhook sink listening", "addr", addr)
	return srv.ListenAndServe()
}
