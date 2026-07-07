# netcheck

Resolves a hostname and times a raw TCP connect — makes DNS resolution and
the TCP handshake (`notes/networking.md`) concrete instead of prose.

```
go run . -host example.com -port 80
example.com -> [104.20.23.154 172.66.147.243], TCP connect to 104.20.23.154:80 in 4.5ms
```

`check()` takes the resolver and dialer as function parameters so
`main_test.go` can inject fakes — no test depends on a real network call,
same double-instead-of-nil-out pattern as `wordcount`'s `upstreamClient`.
