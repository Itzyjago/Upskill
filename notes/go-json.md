# Go JSON notes

From wiring wordcount's `/count` endpoint to return a JSON tally with
`encoding/json`.

## Struct tags control the wire shape
```go
type counts struct {
    Lines int `json:"lines"`
    Words int `json:"words"`
    Bytes int `json:"bytes"`
}
```
- Fields **must be exported** (capitalized) to be marshaled — unexported fields
  are skipped silently. (This is why wordcount's `counts` fields are uppercase.)
- The tag sets the JSON key; without it the key is the Go field name verbatim.
- `omitempty` drops zero-valued fields; `json:"-"` never serializes a field.

## Encode/decode vs marshal/unmarshal
- `json.Marshal(v)` → `[]byte`; `json.Unmarshal(b, &v)` → fills a struct.
- `json.NewEncoder(w).Encode(v)` streams straight to an `io.Writer` (e.g. an
  `http.ResponseWriter`) — no intermediate buffer. `NewDecoder(r).Decode(&v)`
  is the streaming read side.
- Prefer the streaming Encoder/Decoder for HTTP bodies and files; reach for
  Marshal/Unmarshal when you already have/need the bytes.

## Gotchas
- Decoding into `interface{}` makes every number a `float64` — surprises with
  large ints. Decode into a typed struct when you can.
- Unmarshal **ignores unknown fields** by default; call
  `dec.DisallowUnknownFields()` to make typos in input loud.
- A missing JSON key leaves the Go field at its zero value — you can't tell
  "absent" from "sent as zero" unless you use a pointer (`*int`).
- Always set `Content-Type: application/json` before writing the body.

## Ties in
- The HTTP plumbing lives in [http.md](http.md); the handler that uses this is
  wordcount's `countHandler`.
