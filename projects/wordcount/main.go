// wc — a tiny word-count CLI, built to make the Go notes (flag, io, errors) stick.
//
// Usage:
//
//	wc [-l] [-w] [-c] [file...]   # reads stdin when no files given
//
// Flags mirror Unix wc: -l lines, -w words, -c bytes. With no flag, prints all.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
)

type counts struct {
	lines, words, bytes int
}

// count scans r and tallies lines, words, and bytes in a single pass.
func count(r io.Reader) (counts, error) {
	var c counts
	br := bufio.NewReader(r)
	inWord := false
	for {
		b, err := br.ReadByte()
		if err == io.EOF {
			return c, nil
		}
		if err != nil {
			return c, err
		}
		c.bytes++
		if b == '\n' {
			c.lines++
		}
		if b == ' ' || b == '\t' || b == '\n' || b == '\r' {
			inWord = false
		} else if !inWord {
			inWord = true
			c.words++
		}
	}
}

func main() {
	var showL, showW, showC bool
	flag.BoolVar(&showL, "l", false, "count lines")
	flag.BoolVar(&showW, "w", false, "count words")
	flag.BoolVar(&showC, "c", false, "count bytes")
	flag.Parse()

	// No selector flags → show everything, like real wc.
	if !showL && !showW && !showC {
		showL, showW, showC = true, true, true
	}

	print := func(name string, c counts) {
		if showL {
			fmt.Printf("%8d", c.lines)
		}
		if showW {
			fmt.Printf("%8d", c.words)
		}
		if showC {
			fmt.Printf("%8d", c.bytes)
		}
		if name != "" {
			fmt.Printf(" %s", name)
		}
		fmt.Println()
	}

	files := flag.Args()
	if len(files) == 0 {
		c, err := count(os.Stdin)
		if err != nil {
			fmt.Fprintln(os.Stderr, "wc:", err)
			os.Exit(1)
		}
		print("", c)
		return
	}

	var total counts
	exit := 0
	for _, name := range files {
		f, err := os.Open(name)
		if err != nil {
			fmt.Fprintln(os.Stderr, "wc:", err)
			exit = 1
			continue
		}
		c, err := count(f)
		f.Close()
		if err != nil {
			fmt.Fprintln(os.Stderr, "wc:", err)
			exit = 1
			continue
		}
		print(name, c)
		total.lines += c.lines
		total.words += c.words
		total.bytes += c.bytes
	}
	if len(files) > 1 {
		print("total", total)
	}
	os.Exit(exit)
}
