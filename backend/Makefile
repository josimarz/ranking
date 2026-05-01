# Make defaults
SHELL := /bin/bash
.DEFAULT_GOAL := help

# Binaries and paths
GO ?= go
APP_NAME ?= api
BIN_DIR ?= bin
GOPATH_BIN := $(shell $(GO) env GOPATH)/bin

# Colors
COLOR_RESET  = \033[0m
COLOR_GREEN  = \033[32m
COLOR_YELLOW = \033[33m
COLOR_RED    = \033[31m

# Help
.PHONY: help
help: ## Show this help
	@echo ""
	@echo "Available commands:"
	@echo ""
	@awk 'BEGIN {FS = ":.*##"} /^[a-zA-Z0-9_\/-]+:.*##/ {printf "  \033[33m%-24s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST) | sort

# Development
.PHONY: dev/run
dev/run: ## Run the application in development mode
	@echo -e "$(COLOR_YELLOW)Starting application in dev mode...$(COLOR_RESET)"
	$(GO) run ./cmd/$(APP_NAME)

.PHONY: dev/build
dev/build: ## Build the application for local development
	@echo -e "$(COLOR_YELLOW)Building application...$(COLOR_RESET)"
	@mkdir -p $(BIN_DIR)
	$(GO) build -o $(BIN_DIR)/$(APP_NAME) ./cmd/$(APP_NAME)
	@echo -e "$(COLOR_GREEN)Build complete: $(BIN_DIR)/$(APP_NAME)$(COLOR_RESET)"

# Testing
.PHONY: test/unit
test/unit: ## Run unit tests
	@echo -e "$(COLOR_YELLOW)Running unit tests...$(COLOR_RESET)"
	$(GO) test ./... -count=1 -race
	@echo -e "$(COLOR_GREEN)All tests passed.$(COLOR_RESET)"

# Linting
.PHONY: lint/go
lint/go: ## Run Go linter
	@echo -e "$(COLOR_YELLOW)Linting code...$(COLOR_RESET)"
	$(GOPATH_BIN)/golangci-lint run ./...
	@echo -e "$(COLOR_GREEN)Linting passed.$(COLOR_RESET)"

.PHONY: test/integration
test/integration: ## Run integration tests (requires Docker)
	@echo -e "$(COLOR_YELLOW)Running integration tests...$(COLOR_RESET)"
	$(GO) test ./test/integration/... -tags=integration -count=1 -v -timeout=300s
	@echo -e "$(COLOR_GREEN)Integration tests passed.$(COLOR_RESET)"
