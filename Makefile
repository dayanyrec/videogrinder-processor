.PHONY: help setup run test test-js test-js-watch test-js-coverage test-e2e test-e2e-open lint lint-js fmt fmt-js fmt-ci lint-ci test-ci check logs down docker-clean

DOCKER_IMAGE=videogrinder-processor
ENV ?= $(word 2,$(MAKECMDGOALS))
ENV := $(if $(ENV),$(ENV),dev)
PROFILE = $(if $(filter prod,$(ENV)),prod,dev)
SERVICE = $(if $(filter prod,$(ENV)),videogrinder-prod,videogrinder-dev)

# Detect which Docker Compose command is available
DOCKER_COMPOSE := $(shell command -v docker-compose 2> /dev/null)
ifdef DOCKER_COMPOSE
    COMPOSE_CMD = docker-compose
else
    COMPOSE_CMD = docker compose
endif

%:
	@:

help: ## Show available commands
	@echo 'VideoGrinder - Essential Commands:'
	@echo ''
	@echo 'Usage: make <command> [environment]'
	@echo 'Environment: dev (default) | prod'
	@echo 'Examples:'
	@echo '  make run          # Run in dev mode with hot reload'
	@echo '  make run prod     # Run in production mode'
	@echo '  make logs prod    # View production logs'
	@echo '  make down dev     # Stop dev services'
	@echo ''
	@echo 'Development Commands (use the dev container):'
	@echo '  make fmt          # Format code (Go + JS)'
	@echo '  make lint         # Lint code (Go + JS)'
	@echo '  make test         # Run Go tests'
	@echo '  make test-js      # Run JavaScript tests'
	@echo ''
	@echo 'CI/CD Commands (work inside Docker containers):'
	@echo '  make fmt-ci       # Format code (CI-friendly)'
	@echo '  make lint-ci      # Lint code (CI-friendly)'
	@echo '  make test-ci      # Run tests (CI-friendly)'
	@echo ''
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

setup: ## Configure environment (usage: make setup [dev|prod])
	@echo "🔧 Setting up $(ENV) environment..."
	$(COMPOSE_CMD) build $(SERVICE)
	@echo "✅ $(ENV) environment ready"

run: ## Run application with auto-build (usage: make run [dev|prod])
	@echo "🚀 Starting application in $(ENV) mode..."
	$(COMPOSE_CMD) --profile $(PROFILE) up --build $(SERVICE)

test: ## Run Go unit tests
	@echo "🧪 Running Go unit tests..."
	$(COMPOSE_CMD) run --rm videogrinder-dev go test -v ./...

test-js: ## Run JavaScript unit tests
	@echo "🧪 Running JavaScript unit tests..."
	npm install
	npm test

test-js-watch: ## Run unit tests in watch mode
	@echo "🧪 Running unit tests in watch mode..."
	npm install
	npm run test:watch

test-js-coverage: ## Run unit tests with coverage report
	@echo "🧪 Running unit tests with coverage..."
	npm install
	npm run test:coverage

test-e2e: ## Run e2e tests (requires app running)
	@echo "🎭 Running e2e tests..."
	@echo "⚠️  Make sure the app is running with 'make run' in another terminal"
	npm install cypress --save-dev
	npx cypress run

test-e2e-open: ## Open Cypress interactive mode
	@echo "🎭 Opening Cypress..."
	npm install cypress --save-dev
	npx cypress open

lint: ## Check code quality (Go + JS)
	@echo "🔍 Running Go linters..."
	$(COMPOSE_CMD) run --rm videogrinder-dev sh -c "GOFLAGS='-buildvcs=false' golangci-lint run"
	@echo "🔍 Running JS linters..."
	npm install
	npx eslint . --ext .js

lint-js: ## Check JavaScript code quality
	@echo "🔍 Running JS linters..."
	npm install
	npx eslint . --ext .js

fmt: ## Format code (Go + JS)
	@echo "🎨 Formatting Go code..."
	$(COMPOSE_CMD) run --rm videogrinder-dev sh -c "gofmt -s -w . && goimports -w ."
	@echo "🎨 Formatting JS code..."
	npm install
	npx eslint . --ext .js --fix
	@echo "✅ Code formatted"

fmt-ci: ## Format code for CI (without Docker Compose)
	@echo "🎨 Formatting Go code..."
	gofmt -s -w .
	goimports -w .
	@echo "🎨 Formatting JS code..."
	npm install
	npx eslint . --ext .js --fix
	@echo "✅ Code formatted"

fmt-js: ## Format JavaScript code
	@echo "🎨 Formatting JS code..."
	npm install
	npx eslint . --ext .js --fix
	@echo "✅ JS code formatted"

lint-ci: ## Check code quality for CI (without Docker Compose)
	@echo "🔍 Running Go linters..."
	GOFLAGS='-buildvcs=false' golangci-lint run
	@echo "🔍 Running JS linters..."
	npm install
	npx eslint . --ext .js

test-ci: ## Run Go unit tests for CI (without Docker Compose)
	@echo "🧪 Running Go unit tests..."
	GOFLAGS='-buildvcs=false' go test -v ./...

check: fmt lint test test-js ## Run all quality checks

logs: ## View application logs (usage: make logs [dev|prod])
	@echo "📋 Showing $(ENV) logs..."
	$(COMPOSE_CMD) logs -f $(SERVICE)

down: ## Stop services (usage: make down [dev|prod|all])
	@echo "🐳 Stopping $(ENV) services..."
ifeq ($(ENV),all)
	$(COMPOSE_CMD) down
else
	$(COMPOSE_CMD) stop $(SERVICE)
endif

docker-clean: ## Clean Docker resources
	@echo "🧹 Cleaning Docker resources..."
	$(COMPOSE_CMD) down --volumes --rmi all || true
	docker system prune -f || true
