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

## Debugging
- `dig name +short`, `nslookup`. `curl -v` shows DNS → TCP → TLS → HTTP stages.
- `ping` (ICMP reachability), `traceroute`/`mtr` (path + latency per hop).
