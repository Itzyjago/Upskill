package main

import (
	"errors"
	"net"
	"testing"
	"time"
)

func TestCheckSuccess(t *testing.T) {
	resolve := func(_ string) ([]string, error) { return []string{"93.184.216.34"}, nil }
	dial := func(_, _ string, _ time.Duration) (net.Conn, error) {
		client, server := net.Pipe()
		_ = server.Close()
		return client, nil
	}

	r := check(resolve, dial, "example.com", "80")

	if r.resolveErr != nil || r.dialErr != nil {
		t.Fatalf("unexpected error: resolveErr=%v dialErr=%v", r.resolveErr, r.dialErr)
	}
	if r.dialAddr != "93.184.216.34:80" {
		t.Errorf("dialAddr = %q, want %q", r.dialAddr, "93.184.216.34:80")
	}
}

func TestCheckResolveFailureSkipsDial(t *testing.T) {
	wantErr := errors.New("no such host")
	resolve := func(_ string) ([]string, error) { return nil, wantErr }
	dial := func(_, _ string, _ time.Duration) (net.Conn, error) {
		t.Fatal("dial should never run after a resolve failure")
		return nil, nil
	}

	r := check(resolve, dial, "nonexistent.invalid", "80")

	if !errors.Is(r.resolveErr, wantErr) {
		t.Errorf("resolveErr = %v, want %v", r.resolveErr, wantErr)
	}
}

func TestCheckNoAddrsSkipsDial(t *testing.T) {
	resolve := func(_ string) ([]string, error) { return nil, nil }
	dial := func(_, _ string, _ time.Duration) (net.Conn, error) {
		t.Fatal("dial should never run with zero resolved addresses")
		return nil, nil
	}

	r := check(resolve, dial, "empty.invalid", "80")

	if r.resolveErr != nil {
		t.Errorf("resolveErr = %v, want nil", r.resolveErr)
	}
}

func TestCheckDialFailure(t *testing.T) {
	resolve := func(_ string) ([]string, error) { return []string{"10.0.0.1"}, nil }
	wantErr := errors.New("connection refused")
	dial := func(_, _ string, _ time.Duration) (net.Conn, error) {
		return nil, wantErr
	}

	r := check(resolve, dial, "host.invalid", "80")

	if !errors.Is(r.dialErr, wantErr) {
		t.Errorf("dialErr = %v, want %v", r.dialErr, wantErr)
	}
}

func TestCheckResultStringVariants(t *testing.T) {
	cases := []struct {
		name string
		r    checkResult
		want string
	}{
		{"resolve failure", checkResult{host: "h", resolveErr: errors.New("boom")}, "h: DNS lookup failed: boom"},
		{"dial failure", checkResult{host: "h", addrs: []string{"1.2.3.4"}, dialAddr: "1.2.3.4:80", dialErr: errors.New("refused")},
			"h -> [1.2.3.4], TCP connect to 1.2.3.4:80 failed after 0s: refused"},
		{"success", checkResult{host: "h", addrs: []string{"1.2.3.4"}, dialAddr: "1.2.3.4:80"},
			"h -> [1.2.3.4], TCP connect to 1.2.3.4:80 in 0s"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := c.r.String(); got != c.want {
				t.Errorf("String() = %q, want %q", got, c.want)
			}
		})
	}
}
