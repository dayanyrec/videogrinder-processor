.PHONY: help setup run test test-e2e test-e2e-open lint fmt check logs down docker-clean

DOCKER_IMAGE=videogrinder-processor
ENV ?= $(word 2,$(MAKECMDGOALS))
ENV := $(if $(ENV),$(ENV),dev)
PROFILE = $(if $(filter prod,$(ENV)),prod,dev)
SERVICE = $(if $(filter prod,$(ENV)),videogrinder-prod,videogrinder-dev)

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
	docker-compose build $(SERVICE)
	@echo "✅ $(ENV) environment ready"

run: ## Run application with auto-build (usage: make run [dev|prod])
	@echo "🚀 Starting application in $(ENV) mode..."
	docker-compose --profile $(PROFILE) up --build $(SERVICE)

test: ## Run unit tests
	@echo "🧪 Running unit tests..."
	docker-compose run --rm videogrinder-dev go test -v ./...

test-e2e: ## Run e2e tests (requires app running)
	@echo "🎭 Running e2e tests..."
	@echo "⚠️  Make sure the app is running with 'make run' in another terminal"
	npm install cypress --save-dev
	npx cypress run

test-e2e-open: ## Open Cypress interactive mode
	@echo "🎭 Opening Cypress..."
	npm install cypress --save-dev
	npx cypress open

lint: ## Check code quality
	@echo "🔍 Running linters..."
	docker-compose --profile tools run --rm videogrinder-devtools

fmt: ## Format code
	@echo "🎨 Formatting code..."
	docker-compose --profile tools run --rm videogrinder-devtools sh -c "gofmt -s -w . && goimports -w ."
	@echo "✅ Code formatted"

check: fmt lint test ## Run all quality checks

logs: ## View application logs (usage: make logs [dev|prod])
	@echo "📋 Showing $(ENV) logs..."
	docker-compose logs -f $(SERVICE)

down: ## Stop services (usage: make down [dev|prod|all])
	@echo "🐳 Stopping $(ENV) services..."
ifeq ($(ENV),all)
	docker-compose down
else
	docker-compose stop $(SERVICE)
endif

docker-clean: ## Clean Docker resources
	@echo "🧹 Cleaning Docker resources..."
	docker-compose down --volumes --rmi all || true
	docker system prune -f || true
