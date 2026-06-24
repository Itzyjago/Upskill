# wordcount

A tiny `wc` clone in Go — my "build a small CLI to make the concurrency/stdlib
notes stick" roadmap goal, step one.

## Build & run
```sh
go build -o bin/wc .
echo "hello world" | ./bin/wc          # 1  2  12
./bin/wc -w notes.md                    # words only
./bin/wc *.md                           # per-file + total
```

## Test
```sh
go test ./...
```

## What I practiced
- `flag` for CLI parsing and the "no flags → do everything" default.
- Reading from stdin *or* file args; streaming with `bufio` instead of slurping.
- Explicit error returns, non-zero exit on failure, errors to stderr.
- Table-driven tests (see `main_test.go`).

Next: containerize it (see the `Dockerfile`) so it runs anywhere.
