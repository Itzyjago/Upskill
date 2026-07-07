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

## Phony targets chain too (wordcount's real Makefile)
The examples above show one prerequisite; the real Makefile has a three-deep
phony chain that's worth tracing because it breaks the file-target mental
model in a specific way:
```make
.PHONY: image
image: ## Build the container image
	docker build -t $(IMAGE) .

.PHONY: kind-load
kind-load: image ## Build the image and load it into the kind cluster
	kind load docker-image $(IMAGE):latest --name $(CLUSTER)

.PHONY: kind-deploy
kind-deploy: kind-load ## Load the image, then apply the manifest and wait for rollout
	kubectl apply -f deploy/k8s.yaml
	kubectl rollout status deploy/$(IMAGE)
```
- `make kind-deploy` runs `image`, then `kind-load`, then `kind-deploy`'s own
  recipe — Make walks the prerequisite chain depth-first, same as a file
  target's dependency graph.
- **The part that's different from a file target**: the "rebuild only if a
  prerequisite is newer" mtime check at the top of this file doesn't apply
  here at all. `.PHONY` targets have no file, so they're never "up to date"
  — every target in the chain reruns on every invocation, every time. The
  dependency line (`kind-load: image`) is establishing *order*
  ("build the image before loading it"), not *staleness* ("only rebuild the
  image if something changed"). Read a phony chain as a pipeline, not a
  build-avoidance graph — that's the file-target feature this pattern
  deliberately doesn't get.
- `$(MAKEFILE_LIST)` (used in `help`'s `grep`) is a built-in variable Make
  populates with every makefile it read to build the current run — parsing
  it instead of hardcoding the filename is what lets `help` keep working if
  the file gets split (e.g. an `include`d `Makefile.common`) without editing
  the `help` target itself.
