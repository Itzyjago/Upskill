package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHealthHandler(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)

	healthHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if got := strings.TrimSpace(rec.Body.String()); got != "ok" {
		t.Errorf("body = %q, want %q", got, "ok")
	}
}

func TestCountHandler(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/count", strings.NewReader("hello world\n"))

	countHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	var c counts
	if err := json.NewDecoder(rec.Body).Decode(&c); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if c.Lines != 1 || c.Words != 2 || c.Bytes != 12 {
		t.Errorf("got %+v, want {Lines:1 Words:2 Bytes:12}", c)
	}
}

func TestCountHandlerRejectsGET(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/count", nil)

	countHandler(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusMethodNotAllowed)
	}
}

func TestCountHandlerRejectsOversizedBody(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/count", strings.NewReader(strings.Repeat("a", maxCountBodyBytes+1)))

	countHandler(rec, req)

	if rec.Code != http.StatusRequestEntityTooLarge {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusRequestEntityTooLarge)
	}
}

func TestStatusForBodyErr(t *testing.T) {
	if got := statusForBodyErr(&http.MaxBytesError{Limit: 10}); got != http.StatusRequestEntityTooLarge {
		t.Errorf("status = %d, want %d for a MaxBytesError", got, http.StatusRequestEntityTooLarge)
	}
	if got := statusForBodyErr(errors.New("boom")); got != http.StatusBadRequest {
		t.Errorf("status = %d, want %d for any other error", got, http.StatusBadRequest)
	}
}
