// wc — a tiny word-count CLI, built to make the Go notes (flag, io, errors) stick.
//
// Usage:
//
//	wc [-l] [-w] [-c] [file...]   # reads stdin when no files given
//	wc -serve :8080              # run as an HTTP service (POST body to /count)
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
	Lines int `json:"lines"`
	Words int `json:"words"`
	Bytes int `json:"bytes"`
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
		c.Bytes++
		if b == '\n' {
			c.Lines++
		}
		if b == ' ' || b == '\t' || b == '\n' || b == '\r' {
			inWord = false
		} else if !inWord {
			inWord = true
			c.Words++
		}
	}
}

func main() {
	var showL, showW, showC bool
	var serveAddr string
	flag.BoolVar(&showL, "l", false, "count lines")
	flag.BoolVar(&showW, "w", false, "count words")
	flag.BoolVar(&showC, "c", false, "count bytes")
	flag.StringVar(&serveAddr, "serve", "", "run as an HTTP service on this address (e.g. :8080)")
	flag.Parse()

	if serveAddr != "" {
		if err := serve(serveAddr); err != nil {
			fmt.Fprintln(os.Stderr, "wc:", err)
			os.Exit(1)
		}
		return
	}

	// No selector flags → show everything, like real wc.
	if !showL && !showW && !showC {
		showL, showW, showC = true, true, true
	}

	print := func(name string, c counts) {
		if showL {
			fmt.Printf("%8d", c.Lines)
		}
		if showW {
			fmt.Printf("%8d", c.Words)
		}
		if showC {
			fmt.Printf("%8d", c.Bytes)
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
		total.Lines += c.Lines
		total.Words += c.Words
		total.Bytes += c.Bytes
	}
	if len(files) > 1 {
		print("total", total)
	}
	os.Exit(exit)
}
