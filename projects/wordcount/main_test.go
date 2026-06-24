package main

import (
	"strings"
	"testing"
)

func TestCount(t *testing.T) {
	tests := []struct {
		name           string
		in             string
		lines, words   int
		bytes          int
	}{
		{"empty", "", 0, 0, 0},
		{"one word no newline", "hello", 0, 1, 5},
		{"one line", "hello world\n", 1, 2, 12},
		{"extra spaces", "  a   b  \n", 1, 2, 10},
		{"two lines", "a\nb\n", 2, 2, 4},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := count(strings.NewReader(tt.in))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if c.lines != tt.lines || c.words != tt.words || c.bytes != tt.bytes {
				t.Errorf("got {l:%d w:%d c:%d}, want {l:%d w:%d c:%d}",
					c.lines, c.words, c.bytes, tt.lines, tt.words, tt.bytes)
			}
		})
	}
}
