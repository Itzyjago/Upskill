# wordcount task runner — practicing the Make notes on a real project.
# `make` with no target prints help.

BINARY  := wc
IMAGE   := wordcount
PKG     := ./...
CLUSTER := upskill
COMPOSE := deploy/observability/docker-compose.yml

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

.PHONY: kind-up
kind-up: ## Create the local kind cluster (see deploy/kind.yaml)
	kind create cluster --name $(CLUSTER) --config deploy/kind.yaml

.PHONY: kind-load
kind-load: image ## Build the image and load it into the kind cluster
	kind load docker-image $(IMAGE):latest --name $(CLUSTER)

.PHONY: kind-deploy
kind-deploy: kind-load ## Load the image, then apply the manifest and wait for rollout
	kubectl apply -f deploy/k8s.yaml
	kubectl rollout status deploy/$(IMAGE)

.PHONY: kind-metrics-server
kind-metrics-server: ## Install metrics-server on kind (needed for the HPA, notes/kubernetes.md)
	kubectl apply -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml
	kubectl patch deployment metrics-server -n kube-system --type=json \
		-p '[{"op":"add","path":"/spec/template/spec/containers/0/args/-","value":"--kubelet-insecure-tls"}]'

.PHONY: kind-down
kind-down: ## Delete the local kind cluster
	kind delete cluster --name $(CLUSTER)

.PHONY: obs-up
obs-up: ## Bring up the local Prometheus + Grafana stack (Grafana on :3000)
	docker compose -f $(COMPOSE) up --build -d

.PHONY: obs-down
obs-down: ## Tear down the observability stack
	docker compose -f $(COMPOSE) down

.PHONY: clean
clean: ## Remove build artifacts
	rm -rf bin
