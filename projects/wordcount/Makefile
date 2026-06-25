# wordcount task runner — practicing the Make notes on a real project.
# `make` with no target prints help.

BINARY := wc
IMAGE  := wordcount
PKG    := ./...

.DEFAULT_GOAL := help

.PHONY: help
help: ## List the available targets
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
		| awk 'BEGIN {FS = ":.*?## "} {printf "  \033[36m%-12s\033[0m %s\n", $$1, $$2}'

.PHONY: build
build: ## Compile the CLI into bin/wc
	go build -o bin/$(BINARY) .

.PHONY: test
test: ## Run the test suite with the race detector
	go test -race -cover $(PKG)

.PHONY: lint
lint: ## Run golangci-lint (see .golangci.yml)
	golangci-lint run

.PHONY: vet
vet: ## Run go vet
	go vet $(PKG)

.PHONY: run
run: build ## Build, then count this README
	./bin/$(BINARY) README.md

.PHONY: image
image: ## Build the container image
	docker build -t $(IMAGE) .

.PHONY: clean
clean: ## Remove build artifacts
	rm -rf bin
