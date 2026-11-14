SHELL := /bin/bash

GO       ?= go
PODMAN   ?= podman
BIN_DIR  := bin
MODULE   := github.com/aminbenmansour/terrakube-go

VERSION  ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
COMMIT   ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo none)
DATE     ?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ)

LDFLAGS := -s -w \
	-X $(MODULE)/internal/version.Version=$(VERSION) \
	-X $(MODULE)/internal/version.Commit=$(COMMIT) \
	-X $(MODULE)/internal/version.Date=$(DATE)

.PHONY: help all build api executor registry migrate test test-race vet fmt lint tidy clean check \
	services-up services-down services-status services-wipe run-api run-executor run-registry \
	db-migrate docker-build tools generate smoke

help: ## Show this help
	@grep -hE '^[a-zA-Z0-9_-]+:.*?## ' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-18s\033[0m %s\n", $$1, $$2}'

all: build ## Build every binary

build: api executor registry migrate ## Build every binary

api: ## Build the API server
	@mkdir -p $(BIN_DIR)
	$(GO) build -trimpath -ldflags "$(LDFLAGS)" -o $(BIN_DIR)/terrakube-api ./cmd/api

executor: ## Build the executor
	@mkdir -p $(BIN_DIR)
	$(GO) build -trimpath -ldflags "$(LDFLAGS)" -o $(BIN_DIR)/terrakube-executor ./cmd/executor

registry: ## Build the registry
	@mkdir -p $(BIN_DIR)
	$(GO) build -trimpath -ldflags "$(LDFLAGS)" -o $(BIN_DIR)/terrakube-registry ./cmd/registry

migrate: ## Build the migration command
	@mkdir -p $(BIN_DIR)
	$(GO) build -trimpath -ldflags "$(LDFLAGS)" -o $(BIN_DIR)/terrakube-migrate ./cmd/migrate

test: ## Run the unit tests
	$(GO) test ./...

test-race: ## Run the unit tests with the race detector
	$(GO) test -race ./...

vet: ## Run go vet
	$(GO) vet ./...

fmt: ## Format the source tree
	gofmt -l -w .

lint: ## Run golangci-lint when installed
	@command -v golangci-lint >/dev/null || { echo "golangci-lint is not installed"; exit 1; }
	golangci-lint run

tidy: ## Tidy go.mod/go.sum
	$(GO) mod tidy

check: fmt vet test ## Format, vet and test

clean: ## Remove build artifacts and local state
	rm -rf $(BIN_DIR) .data

services-up: ## Start Postgres, Redis and MinIO with podman
	./scripts/dev-services.sh up

services-down: ## Stop the development services (keeps volumes)
	./scripts/dev-services.sh down

services-status: ## Show development service status
	./scripts/dev-services.sh status

services-wipe: ## Stop the development services and delete their volumes
	./scripts/dev-services.sh wipe

run-api: ## Run the API server from source
	$(GO) run ./cmd/api

run-executor: ## Run the executor from source
	$(GO) run ./cmd/executor

run-registry: ## Run the registry from source
	$(GO) run ./cmd/registry

db-migrate: ## Apply database migrations
	$(GO) run ./cmd/migrate

docker-build: ## Build the container image (SERVICE=api|executor|registry)
	$(PODMAN) build --build-arg SERVICE=$(or $(SERVICE),api) -t terrakube-$(or $(SERVICE),api):$(VERSION) -f deploy/docker/Dockerfile .

tools: ## Install code generation tools
	$(GO) install github.com/sqlc-dev/sqlc/cmd/sqlc@v1.31.1

generate: ## Regenerate sqlc code
	sqlc generate

smoke: ## Run the management API smoke test against a running server
	./scripts/smoke.sh
