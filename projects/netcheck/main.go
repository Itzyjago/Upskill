// netcheck — resolves a hostname and times a raw TCP connect, to make DNS
// resolution and the TCP handshake (notes/networking.md) concrete instead of
// prose. `go run . -host example.com -port 80`.
package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"time"
)

// resolveFunc and dialFunc abstract the two network calls check() makes, so
// tests can inject fakes instead of hitting a real resolver/socket — same
// reasoning as wordcount's upstreamClient/fakeCollector doubles.
type resolveFunc func(host string) ([]string, error)
type dialFunc func(network, address string, timeout time.Duration) (net.Conn, error)

// checkResult is a plain data value so String()'s formatting logic is
// testable without a network call ever happening.
type checkResult struct {
	host       string
	addrs      []string
	resolveErr error
	dialAddr   string
	dialTook   time.Duration
	dialErr    error
}

func (r checkResult) String() string {
	if r.resolveErr != nil {
		return fmt.Sprintf("%s: DNS lookup failed: %v", r.host, r.resolveErr)
	}
	if r.dialErr != nil {
		return fmt.Sprintf("%s -> %v, TCP connect to %s failed after %v: %v",
			r.host, r.addrs, r.dialAddr, r.dialTook, r.dialErr)
	}
	return fmt.Sprintf("%s -> %v, TCP connect to %s in %v",
		r.host, r.addrs, r.dialAddr, r.dialTook)
}

// check resolves host, then times a TCP connect to the first returned
// address on port. Stops after a resolve failure — dialing an address that
// didn't resolve isn't a TCP problem to report, it's a DNS one.
func check(resolve resolveFunc, dial dialFunc, host, port string) checkResult {
	res := checkResult{host: host}

	addrs, err := resolve(host)
	res.addrs = addrs
	res.resolveErr = err
	if err != nil || len(addrs) == 0 {
		return res
	}

	res.dialAddr = net.JoinHostPort(addrs[0], port)
	start := time.Now()
	conn, dialErr := dial("tcp", res.dialAddr, 5*time.Second)
	res.dialTook = time.Since(start)
	res.dialErr = dialErr
	if conn != nil {
		_ = conn.Close()
	}
	return res
}

func main() {
	host := flag.String("host", "example.com", "hostname to resolve and dial")
	port := flag.String("port", "80", "port to dial")
	flag.Parse()

	result := check(net.LookupHost, dialTCP, *host, *port)
	fmt.Println(result)
	if result.resolveErr != nil || result.dialErr != nil {
		os.Exit(1)
	}
}

// dialTCP adapts net.DialTimeout to dialFunc's signature.
func dialTCP(network, address string, timeout time.Duration) (net.Conn, error) {
	return net.DialTimeout(network, address, timeout)
}
