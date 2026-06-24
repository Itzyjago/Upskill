# Make notes

A dependency-driven task runner. Still everywhere as a project's "front door".

## Anatomy of a rule
```make
target: prerequisites
	recipe        # MUST be a real TAB, not spaces
```
- A target rebuilds if any prerequisite is *newer* than the target (mtime check).

## Phony targets
- Tasks that aren't files (`build`, `test`, `clean`) must be declared, or Make
  skips them when a file of the same name exists.
```make
.PHONY: build test clean
build:
	go build -o bin/app ./...
test:
	go test ./...
clean:
	rm -rf bin
```

## Variables
- `=` recursive (expanded on use), `:=` simple (expanded once), `?=` set if unset.
- Automatic vars in a recipe: `$@` target, `$<` first prereq, `$^` all prereqs.
```make
bin/app: main.go
	go build -o $@ $<
```

## Useful patterns
- Self-documenting help:
```make
help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## ' $(MAKEFILE_LIST) | \
	  awk -F':.*?## ' '{printf "  %-12s %s\n", $$1, $$2}'
```
- `.DEFAULT_GOAL := help` to make a bare `make` print help.

## Gotchas
- Tabs, not spaces, for recipes — the #1 "missing separator" error.
- Each recipe line runs in its *own* shell; chain with `&&` or `.ONESHELL:`.
