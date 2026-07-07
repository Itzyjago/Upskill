// regexcheck makes notes/regex.md's Go-vs-lookaround claim runnable instead
// of asserted: RE2 (Go's regexp) rejects lookaround syntax outright, so
// parsePrice below uses capture groups plus post-processing to pull the same
// information `(?<=\$)\d+\.\d+(?= each)` would in a backtracking engine.
package main

import (
	"fmt"
	"os"
	"regexp"
)

// priceRe captures the optional "$" prefix, the amount, and the optional
// " each" suffix as groups instead of asserting on them with lookaround --
// the "what to do instead in Go" from notes/regex.md, made real.
var priceRe = regexp.MustCompile(`^(\$)?(\d+\.\d+)( each)?$`)

// parsedPrice is what parsePrice extracts from one input string.
type parsedPrice struct {
	hasDollar bool
	amount    string
	hasEach   bool
}

// parsePrice matches s against priceRe and reports whether the optional
// markers were present. ok is false if s isn't a price-shaped string at all.
func parsePrice(s string) (parsedPrice, bool) {
	m := priceRe.FindStringSubmatch(s)
	if m == nil {
		return parsedPrice{}, false
	}
	return parsedPrice{
		hasDollar: m[1] != "",
		amount:    m[2],
		hasEach:   m[3] != "",
	}, true
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: regexcheck <price-string>")
		os.Exit(1)
	}
	p, ok := parsePrice(os.Args[1])
	if !ok {
		fmt.Printf("%q: not a price-shaped string\n", os.Args[1])
		os.Exit(1)
	}
	fmt.Printf("%q -> dollar=%v amount=%s each=%v\n", os.Args[1], p.hasDollar, p.amount, p.hasEach)
}
