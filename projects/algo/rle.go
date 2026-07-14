package algo

import (
	"strconv"
	"strings"
)

// RunLengthEncode compresses consecutive runs of the same character into
// "<char><count>" pairs (count omitted for runs of length 1). O(n) time.
func RunLengthEncode(s string) string {
	if s == "" {
		return ""
	}
	var b strings.Builder
	count := 1
	for i := 1; i <= len(s); i++ {
		if i < len(s) && s[i] == s[i-1] {
			count++
			continue
		}
		b.WriteByte(s[i-1])
		if count > 1 {
			b.WriteString(strconv.Itoa(count))
		}
		count = 1
	}
	return b.String()
}

// RunLengthDecode reverses RunLengthEncode.
func RunLengthDecode(s string) string {
	var b strings.Builder
	i := 0
	for i < len(s) {
		ch := s[i]
		i++
		start := i
		for i < len(s) && s[i] >= '0' && s[i] <= '9' {
			i++
		}
		count := 1
		if i > start {
			count, _ = strconv.Atoi(s[start:i])
		}
		b.WriteString(strings.Repeat(string(ch), count))
	}
	return b.String()
}
