---
inclusion: auto
name: makefile
description: Guide and best practices for writing Makefile. Apply this skill when writing code in the backend.
---

# Makefile Best Practices for Go Projects

This document defines best practices for creating and maintaining a **Makefile** in Go projects. It standardizes the structure, naming conventions, and behavior of targets to ensure clarity and maintainability.

---

## 1. General Principles

1. **Always Include a `help` Target**

   * The `help` target must display a list of available commands, grouped logically.
   * Each target must have a clear, concise description.

2. **Environment‑First Grouped Naming**

   * Use `/` as the separator. The **first segment is always the environment** (e.g., `dev`, `ci`, `staging`, `prod`, `local`).
   * Patterns:

     ```
     <env>/<command>
     <env>/<domain>/<action>
     ```
   * Examples:

     * `dev/run` – Run the application in development mode.
     * `dev/build` – Build for local development.
     * `ci/test/unit` – Run unit tests in CI.
     * `ci/lint/go` – Run linters in CI.
     * `staging/deploy` – Deploy to staging.
     * `prod/deploy` – Deploy to production.

3. **Use Colors in Outputs**

   * Use ANSI escape codes to colorize messages for better readability:

     * **Green** for success messages.
     * **Yellow** for informational messages.
     * **Red** for errors.

4. **Be Explicit and Safe**

   * Use `.PHONY` for all non-file targets.
   * Prefer explicit dependencies instead of implicit ones.
   * Fail fast: stop execution on the first error (`set -e`).

---

## 2. Environment‑First Naming

* Allowed environments: `dev`, `ci`, `staging`, `prod`, `local`.
* Use lowercase with `-` or `_` if needed; avoid spaces.
* Keep segments short and descriptive (ideally 1–3 segments).
* **Do not** place the environment after the action (avoid `deploy/prod`); **always** use `prod/deploy`.
* Prefer idempotent targets where possible.

---

## 3. Recommended Environment & Domain Map

| Environment | Purpose                              | Example Commands                           |
| ----------- | ------------------------------------ | ------------------------------------------ |
| `dev/`      | Local development workflow           | `dev/run`, `dev/build`                     |
| `ci/`       | Continuous Integration (automation)  | `ci/test/unit`, `ci/lint/go`               |
| `staging/`  | Pre‑production verification          | `staging/deploy`, `staging/smoke`          |
| `prod/`     | Production operations                | `prod/deploy`, `prod/migrate`              |
| `local/`    | Environment‑agnostic local utilities | `local/tools/install`, `local/clean/cache` |

---

## 4. Example Makefile

Below is an example `Makefile` following these best practices:

```makefile
# Make defaults
SHELL := /bin/bash
.DEFAULT_GOAL := help

# Binaries and paths
GO ?= go
APP_NAME ?= app
BIN_DIR ?= bin

# Colors
COLOR_RESET=[0m
COLOR_GREEN=[32m
COLOR_YELLOW=[33m
COLOR_RED=[31m

# Help
.PHONY: help
help: ## Show this help
	@awk 'BEGIN {FS = ":.*##"; printf "Available commands:
"} /^[a-zA-Z0-9_\/-]+:.*##/ {printf "  [33m%-24s[0m %s
", $$1, $$2}' $(MAKEFILE_LIST) | sort

# Development
.PHONY: dev/run
dev/run: ## Run the application in development mode
	@echo "${COLOR_YELLOW}Starting application in dev mode...${COLOR_RESET}"
	$(GO) run ./cmd/$(APP_NAME)

.PHONY: dev/build
dev/build: ## Build the application for local development
	@echo "${COLOR_YELLOW}Building application...${COLOR_RESET}"
	mkdir -p $(BIN_DIR)
	$(GO) build -o $(BIN_DIR)/$(APP_NAME) ./cmd/$(APP_NAME)

# CI tasks
.PHONY: ci/test/unit
ci/test/unit: ## Run unit tests in CI
	@echo "${COLOR_YELLOW}Running unit tests...${COLOR_RESET}"
	$(GO) test ./... -count=1

.PHONY: ci/lint/go
ci/lint/go: ## Run Go linters in CI
	@echo "${COLOR_YELLOW}Linting code...${COLOR_RESET}"
	golangci-lint run

# Deployments
.PHONY: staging/deploy
staging/deploy: ## Deploy to staging
	@echo "${COLOR_YELLOW}Deploying to staging...${COLOR_RESET}"
	# TODO: add staging deploy steps

.PHONY: prod/deploy
prod/deploy: ## Deploy to production
	@echo "${COLOR_RED}Deploying to production...${COLOR_RESET}"
	# TODO: add production deploy steps
```

---

| Environment.   | Purpose                                         | Example Commands                |
|----------------|-------------------------------------------------|---------------------------------|
| `dev/`         | Local development tasks                         | `dev/run`, `dev/build`          |
| `test/`        | Running tests (unit, integration, etc.)         | `test/unit`, `test/integration` |
| `lint/`        | Code analysis and linting                       | `lint/go`, `lint/security`      |
| `build/`       | Production-ready builds                         | `build/release`                 |
| `staging/`     | Deployment to staging environment               | `staging/deploy`                |
| `prod/`        | Deployment to production environment            | `prod/deploy`                   |
| `tools/`       | Tool installation and management                | `tools/install`                 |

---

## 3. Example Makefile

Below is an example `Makefile` following these best practices:

```makefile
# Colors
COLOR_RESET=\033[0m
COLOR_GREEN=\033[32m
COLOR_YELLOW=\033[33m
COLOR_RED=\033[31m

# Help
.PHONY: help
help:
	@echo "Available commands:"
	@grep -E '^([a-zA-Z_-]+/[^:]+):.*##' Makefile | sort | awk 'BEGIN {FS = ":.*##"}; {printf "  \033[33m%-20s\033[0m %s\n", $$1, $$2}'

# Development group
.PHONY: dev/run
## Run the application in development mode
dev/run:
	@echo "${COLOR_YELLOW}Starting application in dev mode...${COLOR_RESET}"
	go run ./cmd/app

.PHONY: dev/build
## Build the application for development
dev/build:
	@echo "${COLOR_YELLOW}Building application...${COLOR_RESET}"
	go build -o bin/app ./cmd/app

# Testing group
.PHONY: test/unit
## Run unit tests
test/unit:
	@echo "${COLOR_YELLOW}Running unit tests...${COLOR_RESET}"
	go test ./...

# Lint group
.PHONY: lint/go
## Run Go linter
lint/go:
	@echo "${COLOR_YELLOW}Linting code...${COLOR_RESET}"
	golangci-lint run

# Production group
.PHONY: prod/deploy
## Deploy application to production
prod/deploy:
	@echo "${COLOR_YELLOW}Deploying to production...${COLOR_RESET}"
	./scripts/deploy_prod.sh
```

---

## 5. Additional Recommendations

1. **Consistency**: Maintain the same group naming structure across multiple projects.
2. **Documentation**: Update the `help` target whenever new commands are added.
3. **Graceful Error Handling**: Use clear, colorized error messages to indicate failure causes.
4. **Version Control**: Always include the `Makefile` in version control for team collaboration.

---

By following these rules, Go projects will have **organized, maintainable, and self-documented** Makefiles that improve developer experience and workflow.