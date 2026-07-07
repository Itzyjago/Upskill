# Networking notes

## The stack, briefly
- Link → Internet (IP) → Transport (TCP/UDP) → Application (HTTP, DNS, ...).
- IP routes packets between hosts; TCP/UDP multiplex them to a process via ports.

## TCP vs UDP
- TCP: connection-oriented, ordered, reliable (retransmits, flow + congestion
  control). Cost: handshake + head-of-line blocking.
- UDP: fire-and-forget datagrams. No ordering/retransmit. DNS, QUIC, real-time
  media build their own reliability on top.

## TCP handshake & teardown
- Open: SYN → SYN-ACK → ACK (3-way).
- Close: FIN/ACK both directions; `TIME_WAIT` holds the socket briefly to catch
  stragglers.

## DNS resolution
- Recursive resolver → root → TLD (`.com`) → authoritative nameserver.
- Records: `A`/`AAAA` (IP), `CNAME` (alias), `MX` (mail), `TXT` (verification),
  `NS` (delegation).
- TTL controls caching — low TTL before a migration, raise it after.

## TLS handshake (1.3)
- ClientHello (offers ciphers + key share) → ServerHello + certificate →
  keys derived → encrypted. 1-RTT, or 0-RTT on resume.
- Cert chains to a trusted CA; the client verifies hostname + expiry + signature.

### Why 1.3 dropped a round trip, and what 0-RTT actually costs
- **TLS 1.2 needed 2 round trips** to agree on a cipher: ClientHello lists
  supported ciphers, ServerHello *picks one*, and only then can key material
  actually be derived — negotiation and key exchange were sequential.
- **TLS 1.3 collapses this to 1-RTT** by having the client guess: ClientHello
  now *includes* a key share for its most-likely-supported cipher alongside
  the offer list. If the server supports that guess (the common case,
  because TLS 1.3 trimmed the cipher list to a handful of good ones), it
  replies with ServerHello + cert in the same round trip, keys already
  derivable. Negotiation and key exchange happen in parallel instead of in
  sequence — that's the whole speedup, not a shortened algorithm.
- **0-RTT goes further and skips the round trip entirely** on a *resumed*
  connection: the client reuses a key derived from a previous session and
  sends encrypted application data (the actual HTTP request) in the very
  first flight, before the server has said anything. Fast, but it comes with
  a real cost the "0-RTT on resume" one-liner glosses over: that first
  flight has **no replay protection**. An attacker who captured it can
  resend the exact same encrypted bytes, and the server has no way to tell
  the replay from the original — it never got to run its usual
  nonce/freshness check before processing.
- **The mitigation is exactly the idempotency question from
  [http.md](http.md)**: 0-RTT early data should only ever carry
  safe/idempotent requests (`GET`, `HEAD`) — replaying one just re-reads the
  same thing. A non-idempotent request (a `POST` that transfers money,
  charges a card) riding in 0-RTT turns a captured packet into a free
  replay of that side effect. Same underlying question as `/count`'s
  `Idempotency-Key` discussion, one layer down the stack: "is doing this
  twice safe" isn't just an API design concern, TLS 1.3 made it a transport
  *security* concern the moment it added a mode with no replay protection.

## Debugging
- `dig name +short`, `nslookup`. `curl -v` shows DNS → TCP → TLS → HTTP stages.
- `ping` (ICMP reachability), `traceroute`/`mtr` (path + latency per hop).
