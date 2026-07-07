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

func TestCountHandlerReplaysIdempotentRetry(t *testing.T) {
	store := newIdempotencyStore()
	handler := countHandlerFunc(store)

	req1 := httptest.NewRequest(http.MethodPost, "/count", strings.NewReader("hello world\n"))
	req1.Header.Set(idempotencyKeyHeader, "retry-1")
	rec1 := httptest.NewRecorder()
	handler(rec1, req1)

	if rec1.Code != http.StatusOK {
		t.Fatalf("first request status = %d, want %d", rec1.Code, http.StatusOK)
	}
	if rec1.Header().Get("Idempotency-Replayed") != "" {
		t.Error("first request should not be marked as replayed")
	}

	req2 := httptest.NewRequest(http.MethodPost, "/count", strings.NewReader("hello world\n"))
	req2.Header.Set(idempotencyKeyHeader, "retry-1")
	rec2 := httptest.NewRecorder()
	handler(rec2, req2)

	if rec2.Code != http.StatusOK {
		t.Fatalf("replayed request status = %d, want %d", rec2.Code, http.StatusOK)
	}
	if rec2.Header().Get("Idempotency-Replayed") != "true" {
		t.Error("replayed request missing Idempotency-Replayed header")
	}
	if rec1.Body.String() != rec2.Body.String() {
		t.Errorf("replayed body = %q, want it to match the original %q", rec2.Body.String(), rec1.Body.String())
	}
}

func TestCountHandlerRejectsReusedKeyDifferentBody(t *testing.T) {
	store := newIdempotencyStore()
	handler := countHandlerFunc(store)

	req1 := httptest.NewRequest(http.MethodPost, "/count", strings.NewReader("hello\n"))
	req1.Header.Set(idempotencyKeyHeader, "retry-2")
	handler(httptest.NewRecorder(), req1)

	req2 := httptest.NewRequest(http.MethodPost, "/count", strings.NewReader("goodbye\n"))
	req2.Header.Set(idempotencyKeyHeader, "retry-2")
	rec2 := httptest.NewRecorder()
	handler(rec2, req2)

	if rec2.Code != http.StatusConflict {
		t.Errorf("status = %d, want %d for a reused key with a different body", rec2.Code, http.StatusConflict)
	}
}

func TestCountHandlerWithoutKeyAlwaysRecounts(t *testing.T) {
	store := newIdempotencyStore()
	handler := countHandlerFunc(store)

	for i := 0; i < 2; i++ {
		req := httptest.NewRequest(http.MethodPost, "/count", strings.NewReader("hello world\n"))
		rec := httptest.NewRecorder()
		handler(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("request %d status = %d, want %d", i, rec.Code, http.StatusOK)
		}
		if rec.Header().Get("Idempotency-Replayed") != "" {
			t.Errorf("request %d without a key should never be marked replayed", i)
		}
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
