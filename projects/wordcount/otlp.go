package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// otlpExporter ships finished spans to an OTLP/HTTP collector as JSON — the
// dependency-free half of distributed tracing (notes/otlp.md). No OpenTelemetry
// SDK: we marshal the ResourceSpans envelope and POST it to <endpoint>/v1/traces
// ourselves, the same "implement the wire format to learn it" spirit as the
// hand-rolled Prometheus exposition in metrics.go.
type otlpExporter struct {
	endpoint string // collector base URL, e.g. http://jaeger:4318
	service  string // service.name reported on the Resource
	client   *http.Client
}

func newOTLPExporter(endpoint, service string) *otlpExporter {
	return &otlpExporter{
		endpoint: strings.TrimRight(endpoint, "/"),
		service:  service,
		client:   &http.Client{Timeout: 2 * time.Second},
	}
}

// kv is one string-valued OTLP attribute (the only value type we emit).
type kv struct{ key, val string }

// --- OTLP/JSON wire structs: only the fields we populate (notes/otlp.md) ---

type otlpPayload struct {
	ResourceSpans []otlpResourceSpans `json:"resourceSpans"`
}
type otlpResourceSpans struct {
	Resource   otlpResource     `json:"resource"`
	ScopeSpans []otlpScopeSpans `json:"scopeSpans"`
}
type otlpResource struct {
	Attributes []otlpKeyValue `json:"attributes"`
}
type otlpScopeSpans struct {
	Scope otlpScope  `json:"scope"`
	Spans []otlpSpan `json:"spans"`
}
type otlpScope struct {
	Name string `json:"name"`
}
type otlpSpan struct {
	TraceID           string         `json:"traceId"`
	SpanID            string         `json:"spanId"`
	ParentSpanID      string         `json:"parentSpanId,omitempty"` // omitted on a root span
	Name              string         `json:"name"`
	Kind              int            `json:"kind"`
	StartTimeUnixNano string         `json:"startTimeUnixNano"`
	EndTimeUnixNano   string         `json:"endTimeUnixNano"`
	Attributes        []otlpKeyValue `json:"attributes,omitempty"`
	Status            otlpStatus     `json:"status"`
}
type otlpStatus struct {
	Code int `json:"code"`
}
type otlpKeyValue struct {
	Key   string       `json:"key"`
	Value otlpAnyValue `json:"value"`
}
type otlpAnyValue struct {
	StringValue string `json:"stringValue"`
}

const (
	spanKindServer  = 2 // OTLP SpanKind.SERVER — an inbound request
	spanKindClient  = 3 // OTLP SpanKind.CLIENT — an outbound call we made
	statusCodeUnset = 0 // leave success UNSET, per OTel HTTP semconv
	statusCodeError = 2 // ERROR
)

// payloadFor builds the OTLP envelope for one finished span. Split out from the
// HTTP send so tests can assert on the JSON without standing up a collector.
func (e *otlpExporter) payloadFor(s span, attrs []kv) otlpPayload {
	code := statusCodeUnset
	if s.failed {
		code = statusCodeError
	}
	return otlpPayload{ResourceSpans: []otlpResourceSpans{{
		Resource: otlpResource{Attributes: []otlpKeyValue{
			strAttr("service.name", e.service),
		}},
		ScopeSpans: []otlpScopeSpans{{
			Scope: otlpScope{Name: "wordcount/serve"},
			Spans: []otlpSpan{{
				TraceID:           s.sc.traceID,
				SpanID:            s.sc.spanID,
				ParentSpanID:      s.parentID,
				Name:              s.name,
				Kind:              s.kind,
				StartTimeUnixNano: nano(s.start),
				EndTimeUnixNano:   nano(s.end),
				Attributes:        toAttrs(attrs),
				Status:            otlpStatus{Code: code},
			}},
		}},
	}}}
}

// export POSTs the span to <endpoint>/v1/traces. Best-effort: a collector that's
// down returns an error the caller logs and ignores — tracing never fails a
// request (the same rule parseTraceparent follows in trace.go).
func (e *otlpExporter) export(ctx context.Context, s span, attrs []kv) error {
	body, err := json.Marshal(e.payloadFor(s, attrs))
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		e.endpoint+"/v1/traces", bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := e.client.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode/100 != 2 {
		return fmt.Errorf("otlp: collector returned %s", resp.Status)
	}
	return nil
}

// nano renders a time the way OTLP wants it: unix nanoseconds as a *string*. A
// JSON number would silently lose precision on a 64-bit value (notes/otlp.md).
func nano(t time.Time) string { return strconv.FormatInt(t.UnixNano(), 10) }

func strAttr(k, v string) otlpKeyValue {
	return otlpKeyValue{Key: k, Value: otlpAnyValue{StringValue: v}}
}

func toAttrs(attrs []kv) []otlpKeyValue {
	out := make([]otlpKeyValue, 0, len(attrs))
	for _, a := range attrs {
		out = append(out, strAttr(a.key, a.val))
	}
	return out
}
