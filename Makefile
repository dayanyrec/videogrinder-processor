.PHONY: help setup run test test-e2e test-e2e-open lint lint-js fmt fmt-js check logs down docker-clean ci-validate ci-build ci-test-local

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
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

setup: ## Configure environment (usage: make setup [dev|prod])
	@echo "🔧 Setting up $(ENV) environment..."
	$(COMPOSE_CMD) build $(SERVICE)
	@echo "✅ $(ENV) environment ready"

run: ## Run application with auto-build (usage: make run [dev|prod])
	@echo "🚀 Starting application in $(ENV) mode..."
	$(COMPOSE_CMD) --profile $(PROFILE) up --build $(SERVICE)

test: ## Run unit tests
	@echo "🧪 Running unit tests..."
	$(COMPOSE_CMD) run --rm videogrinder-dev go test -v ./...

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
	$(COMPOSE_CMD) --profile tools run --rm videogrinder-devtools
	@echo "🔍 Running JS linters..."
	npm install
	npx eslint . --ext .js

lint-js: ## Check JavaScript code quality
	@echo "🔍 Running JS linters..."
	npm install
	npx eslint . --ext .js

fmt: ## Format code (Go + JS)
	@echo "🎨 Formatting Go code..."
	$(COMPOSE_CMD) --profile tools run --rm videogrinder-devtools sh -c "gofmt -s -w . && goimports -w ."
	@echo "🎨 Formatting JS code..."
	npm install
	npx eslint . --ext .js --fix
	@echo "✅ Code formatted"

fmt-js: ## Format JavaScript code
	@echo "🎨 Formatting JS code..."
	npm install
	npx eslint . --ext .js --fix
	@echo "✅ JS code formatted"

check: fmt lint test ## Run all quality checks

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

ci-validate: ## Run CI validation locally (equivalent to PR validation)
	@echo "🔍 Running CI validation locally..."
	@echo "📁 Creating directories..."
	@mkdir -p uploads outputs temp tmp
	@echo "🎨 Running formatting..."
	@make fmt
	@echo "🔍 Running linting..."
	@make lint
	@echo "🧪 Running unit tests..."
	@make test
	@echo "✅ CI validation completed successfully!"

ci-build: ## Build production image (like CI)
	@echo "🏗️ Building production image for CI validation..."
	@make setup prod
	@echo "✅ Production image built successfully!"

ci-test-local: ## Run complete CI test suite locally
	@echo "🚀 Running complete CI test suite locally..."
	@make ci-validate
	@make ci-build
	@echo "🎭 Running E2E tests..."
	@make run dev &
	@sleep 10
	@make test-e2e || (make down && exit 1)
	@make down
	@echo "🎉 Complete CI test suite passed!"
