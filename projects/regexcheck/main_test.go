package main

import (
	"regexp"
	"strings"
	"testing"
)

func TestParsePrice(t *testing.T) {
	cases := []struct {
		in   string
		want parsedPrice
		ok   bool
	}{
		{"$42.50 each", parsedPrice{hasDollar: true, amount: "42.50", hasEach: true}, true},
		{"$42.50", parsedPrice{hasDollar: true, amount: "42.50", hasEach: false}, true},
		{"42.50 each", parsedPrice{hasDollar: false, amount: "42.50", hasEach: true}, true},
		{"42.50", parsedPrice{hasDollar: false, amount: "42.50", hasEach: false}, true},
		{"not a price", parsedPrice{}, false},
	}
	for _, c := range cases {
		t.Run(c.in, func(t *testing.T) {
			got, ok := parsePrice(c.in)
			if ok != c.ok {
				t.Fatalf("ok = %v, want %v", ok, c.ok)
			}
			if ok && got != c.want {
				t.Errorf("parsePrice(%q) = %+v, want %+v", c.in, got, c.want)
			}
		})
	}
}

// TestRE2RejectsLookahead pins the claim in notes/regex.md as a real
// regression test instead of prose someone has to trust: RE2 doesn't just
// handle lookahead differently, it refuses to compile it at all.
func TestRE2RejectsLookahead(t *testing.T) {
	//nolint:staticcheck // SA1000: the whole point is that this doesn't compile.
	_, err := regexp.Compile(`\d+(?= USD)`)
	if err == nil {
		t.Fatal("expected a compile error for lookahead syntax, got nil -- did RE2 gain lookahead support?")
	}
	if !strings.Contains(err.Error(), "invalid or unsupported Perl syntax") {
		t.Errorf("error = %q, want it to mention unsupported Perl syntax", err.Error())
	}
}

// TestRE2RejectsLookbehind pins the other half: RE2 parses (?<=...) as a
// malformed *named capture group* attempt, a different error for the same
// missing feature -- worth its own test since the error text differs enough
// that someone could mistake it for an unrelated bug.
func TestRE2RejectsLookbehind(t *testing.T) {
	//nolint:staticcheck // SA1000: the whole point is that this doesn't compile.
	_, err := regexp.Compile(`(?<=\$)\d+`)
	if err == nil {
		t.Fatal("expected a compile error for lookbehind syntax, got nil -- did RE2 gain lookbehind support?")
	}
	if !strings.Contains(err.Error(), "invalid named capture") {
		t.Errorf("error = %q, want it to mention invalid named capture", err.Error())
	}
}
