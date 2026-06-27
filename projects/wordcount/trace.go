package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"strconv"
	"strings"
)

// W3C Trace Context, hand-rolled — same spirit as the hand-rolled Prometheus
// exposition in metrics.go: implement the wire format to learn it, no
// OpenTelemetry SDK. See notes/trace-context.md.
//
// traceparent: version "-" trace-id "-" parent-id "-" flags
//   00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01
//   version  : 2 hex   (we only speak "00")
//   trace-id : 32 hex  (16 bytes), never all-zero
//   span-id  : 16 hex  (8 bytes),  never all-zero
//   flags    : 2 hex   (bit 0 = sampled)

// spanContext is the trace state for the current span: the trace it belongs to,
// this span's own id, and whether the trace is sampled.
type spanContext struct {
	traceID string // 32 lowercase hex
	spanID  string // 16 lowercase hex
	sampled bool
}

// ctxKey is an unexported type so our context value can't collide with another
// package's key — the standard context.WithValue idiom.
type ctxKey struct{}

// randHex returns n random bytes as lowercase hex (2n chars), from a CSPRNG.
func randHex(n int) string {
	b := make([]byte, n)
	_, _ = rand.Read(b) // crypto/rand.Read never returns a short read or error
	return hex.EncodeToString(b)
}

// newSpanContext starts a fresh root span: a new trace and a new span, sampled.
func newSpanContext() spanContext {
	return spanContext{traceID: randHex(16), spanID: randHex(8), sampled: true}
}

// child returns a new span in the same trace as the parent — same trace-id, a
// fresh span-id, carrying the sampled flag. This is the "keep the trace-id,
// mint a new span-id" rule that builds the span tree across hops.
func (sc spanContext) child() spanContext {
	return spanContext{traceID: sc.traceID, spanID: randHex(8), sampled: sc.sampled}
}

// traceparent renders the header for propagation to a downstream service.
func (sc spanContext) traceparent() string {
	flags := "00"
	if sc.sampled {
		flags = "01"
	}
	return "00-" + sc.traceID + "-" + sc.spanID + "-" + flags
}

// parseTraceparent parses an inbound traceparent header. The returned context's
// spanID is the *sender's* span (our parent-to-be). On any malformed input it
// returns ok=false so the caller starts a fresh trace rather than failing the
// request — tracing is best-effort, never part of the contract.
func parseTraceparent(h string) (spanContext, bool) {
	parts := strings.Split(h, "-")
	if len(parts) != 4 {
		return spanContext{}, false
	}
	version, traceID, parentID, flags := parts[0], parts[1], parts[2], parts[3]

	if version != "00" {
		return spanContext{}, false
	}
	if len(traceID) != 32 || !isLowerHex(traceID) || isAllZero(traceID) {
		return spanContext{}, false
	}
	if len(parentID) != 16 || !isLowerHex(parentID) || isAllZero(parentID) {
		return spanContext{}, false
	}
	if len(flags) != 2 || !isLowerHex(flags) {
		return spanContext{}, false
	}

	v, _ := strconv.ParseUint(flags, 16, 8) // validated as 2 hex chars above
	return spanContext{traceID: traceID, spanID: parentID, sampled: v&1 == 1}, true
}

// isLowerHex reports whether s is non-empty and all lowercase hex digits. W3C
// requires lowercase, so an uppercase header is treated as malformed.
func isLowerHex(s string) bool {
	if s == "" {
		return false
	}
	for i := 0; i < len(s); i++ {
		c := s[i]
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f')) {
			return false
		}
	}
	return true
}

// isAllZero reports whether every character is '0' — an all-zero trace or span
// id is invalid per the spec.
func isAllZero(s string) bool {
	return strings.Trim(s, "0") == ""
}

// withSpan stashes a spanContext in the request context so handlers (and the log
// line) downstream can read it. Mirrors how go-context.md threads a deadline.
func withSpan(ctx context.Context, sc spanContext) context.Context {
	return context.WithValue(ctx, ctxKey{}, sc)
}

// spanFrom pulls the spanContext back out, if one was set.
func spanFrom(ctx context.Context) (spanContext, bool) {
	sc, ok := ctx.Value(ctxKey{}).(spanContext)
	return sc, ok
}
