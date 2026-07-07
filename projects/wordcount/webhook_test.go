package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAlertWebhookHandlerAcceptsAValidPayload(t *testing.T) {
	rec := httptest.NewRecorder()
	body := `{"status":"firing","alerts":[{"status":"firing",` +
		`"labels":{"alertname":"HighErrorRate","severity":"page"},` +
		`"annotations":{"summary":"5xx ratio above threshold"}}]}`
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))

	alertWebhookHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
}

func TestAlertWebhookHandlerRejectsGET(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	alertWebhookHandler(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusMethodNotAllowed)
	}
}

func TestAlertWebhookHandlerRejectsMalformedJSON(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("not json"))

	alertWebhookHandler(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

// TestAlertWebhookHandlerRejectsOversizedBody is the actual security-audit
// fix: this endpoint used to json.NewDecoder(r.Body) an inbound payload with
// no cap at all — the same resource-exhaustion shape countHandler and
// forwardCountHandler were already fixed against (roadmap #16), just missed
// here (notes/security.md).
func TestAlertWebhookHandlerRejectsOversizedBody(t *testing.T) {
	rec := httptest.NewRecorder()
	oversized := `{"status":"` + strings.Repeat("a", maxCountBodyBytes+1) + `"}`
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(oversized))

	alertWebhookHandler(rec, req)

	if rec.Code != http.StatusRequestEntityTooLarge {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusRequestEntityTooLarge)
	}
}
