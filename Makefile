.PHONY: help setup run run-api run-web run-processor test test-api test-processor test-services test-utils test-clients test-js test-js-watch test-js-coverage test-e2e test-e2e-open lint lint-js fmt fmt-js fmt-ci lint-ci test-ci test-api-ci test-processor-ci test-js-ci check check-ci logs logs-web logs-api logs-processor down docker-clean health health-ci

DOCKER_IMAGE=videogrinder-processor
ENV ?= $(word 2,$(MAKECMDGOALS))
ENV := $(if $(ENV),$(ENV),dev)
PROFILE = $(if $(filter prod,$(ENV)),prod,dev)
SERVICE = $(if $(filter prod,$(ENV)),videogrinder-prod,videogrinder-dev)
API_SERVICE = $(if $(filter prod,$(ENV)),videogrinder-api-prod,videogrinder-api-dev)
WEB_SERVICE = $(if $(filter prod,$(ENV)),videogrinder-web-prod,videogrinder-web-dev)
PROCESSOR_SERVICE = $(if $(filter prod,$(ENV)),videogrinder-processor-prod,videogrinder-processor-dev)

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
	@echo ''
	@echo 'Multi-Service Architecture:'
	@echo '  make run          # Run all 3 services (Web + API + Processor)'
	@echo '  make run-web      # Run only Web service (static files)'
	@echo '  make run-api      # Run only API service'
	@echo '  make run-processor # Run only processor service'
	@echo ''
	@echo 'Port Configuration:'
	@echo '  Web Service:      http://localhost:8080 (static files)'
	@echo '  API Service:      http://localhost:8081 (REST API)'
	@echo '  Processor Service: http://localhost:8082 (video processing)'
	@echo ''
	@echo 'Testing:'
	@echo '  make test         # Run all Go tests (API + processor)'
	@echo '  make test-api     # Run only API service tests'
	@echo '  make test-processor # Run only processor service tests'
	@echo '  make test-js      # Run JavaScript tests'
	@echo '  make test-e2e     # Run end-to-end tests'
	@echo ''
	@echo 'Examples:'
	@echo '  make run dev      # Run all 3 services in dev mode'
	@echo '  make run prod     # Run all 3 services in production mode'
	@echo '  make logs prod    # View production logs'
	@echo '  make down dev     # Stop dev services'
	@echo ''
	@echo 'Quality Checks:'
	@echo '  make check        # Run all quality checks (format + lint + test)'
	@echo '  make fmt          # Format code (Go + JS)'
	@echo '  make lint         # Lint code (Go + JS)'
	@echo '  make test         # Run all Go tests'
	@echo '  make test-js      # Run JS tests'
	@echo '  make health       # Check app health'
	@echo ''
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

setup: ## Configure environment (usage: make setup [dev|prod])
	@echo "🔧 Setting up $(ENV) environment..."
	$(COMPOSE_CMD) build $(API_SERVICE) $(WEB_SERVICE) $(PROCESSOR_SERVICE)
	@echo "✅ $(ENV) environment ready"

run: ## Run all 3 services (Web + API + Processor) (usage: make run [dev|prod])
	@echo "🚀 Starting all 3 services (Web + API + Processor) in $(ENV) mode..."
	$(COMPOSE_CMD) --profile $(PROFILE) up -d --build $(WEB_SERVICE) $(API_SERVICE) $(PROCESSOR_SERVICE)

run-api: ## Run only API service (usage: make run-api [dev|prod])
	@echo "🎬 Starting API service in $(ENV) mode..."
	$(COMPOSE_CMD) --profile $(PROFILE) up --build $(API_SERVICE)

run-web: ## Run only web service (usage: make run-web [dev|prod])
	@echo "🎬 Starting web service in $(ENV) mode..."
	$(COMPOSE_CMD) --profile $(PROFILE) up --build $(WEB_SERVICE)

run-processor: ## Run only processor service (usage: make run-processor [dev|prod])
	@echo "🔧 Starting processor service in $(ENV) mode..."
	$(COMPOSE_CMD) --profile $(PROFILE) up --build $(PROCESSOR_SERVICE)

test: ## Run all Go unit tests (API + processor)
	@echo "🧪 Running all Go unit tests..."
	$(COMPOSE_CMD) --profile tools run --rm videogrinder-devtools sh -c "GOFLAGS='-buildvcs=false' go test -v ./..."

test-api: ## Run API service unit tests
	@echo "🧪 Running API service unit tests..."
	$(COMPOSE_CMD) --profile tools run --rm videogrinder-devtools sh -c "GOFLAGS='-buildvcs=false' go test -v ./api/internal/..."

test-processor: ## Run processor service unit tests
	@echo "🧪 Running processor service unit tests..."
	$(COMPOSE_CMD) --profile tools run --rm videogrinder-devtools sh -c "GOFLAGS='-buildvcs=false' go test -v ./processor/internal/..."

test-services: ## Run services unit tests (API + processor)
	@echo "🧪 Running services unit tests..."
	$(COMPOSE_CMD) --profile tools run --rm videogrinder-devtools sh -c "GOFLAGS='-buildvcs=false' go test -v ./internal/services/..."

test-utils: ## Run utils unit tests
	@echo "🧪 Running utils unit tests..."
	$(COMPOSE_CMD) --profile tools run --rm videogrinder-devtools sh -c "GOFLAGS='-buildvcs=false' go test -v ./internal/utils/..."

test-clients: ## Run clients unit tests
	@echo "🧪 Running clients unit tests..."
	$(COMPOSE_CMD) --profile tools run --rm videogrinder-devtools sh -c "GOFLAGS='-buildvcs=false' go test -v ./internal/clients/..."

test-js: ## Run JavaScript unit tests
	@echo "🧪 Running JavaScript unit tests..."
	$(COMPOSE_CMD) --profile tools run --rm videogrinder-devtools sh -c "cd web && npm install"
	$(COMPOSE_CMD) --profile tools run --rm videogrinder-devtools sh -c "cd web && npm test"

test-js-watch: ## Run unit tests in watch mode
	@echo "🧪 Running unit tests in watch mode..."
	$(COMPOSE_CMD) --profile tools run --rm videogrinder-devtools sh -c "cd web && npm install"
	$(COMPOSE_CMD) --profile tools run --rm videogrinder-devtools sh -c "cd web && npm run test:watch"

test-js-coverage: ## Run unit tests with coverage report
	@echo "🧪 Running unit tests with coverage..."
	$(COMPOSE_CMD) --profile tools run --rm videogrinder-devtools sh -c "cd web && npm install"
	$(COMPOSE_CMD) --profile tools run --rm videogrinder-devtools sh -c "cd web && npm run test:coverage"

test-e2e: ## Run e2e tests (requires app running)
	@echo "🎭 Running e2e tests..."
	@echo "⚠️  Make sure the app is running with 'make run' in another terminal"
	$(COMPOSE_CMD) --profile tools run --rm videogrinder-devtools sh -c "cd web && npm install cypress --save-dev && npx cypress install"
	$(COMPOSE_CMD) --profile tools run --rm videogrinder-devtools sh -c "cd web && npx cypress run"

test-e2e-open: ## Open Cypress interactive mode
	@echo "🎭 Opening Cypress..."
	$(COMPOSE_CMD) --profile tools run --rm videogrinder-devtools sh -c "cd web && npm install cypress --save-dev && npx cypress install"
	$(COMPOSE_CMD) --profile tools run --rm videogrinder-devtools sh -c "cd web && npx cypress open"

lint: ## Check code quality (Go + JS)
	@echo "🔍 Running Go linters..."
	$(COMPOSE_CMD) --profile tools run --rm videogrinder-devtools
	@echo "🔍 Running JS linters..."
	$(COMPOSE_CMD) --profile tools run --rm videogrinder-devtools sh -c "cd web && npm install"
	$(COMPOSE_CMD) --profile tools run --rm videogrinder-devtools sh -c "cd web && npx eslint . --ext .js"

lint-js: ## Check JavaScript code quality
	@echo "🔍 Running JS linters..."
	$(COMPOSE_CMD) --profile tools run --rm videogrinder-devtools sh -c "cd web && npm install"
	$(COMPOSE_CMD) --profile tools run --rm videogrinder-devtools sh -c "cd web && npx eslint . --ext .js"

fmt: ## Format code (Go + JS)
	@echo "🎨 Formatting Go code..."
	$(COMPOSE_CMD) --profile tools run --rm videogrinder-devtools sh -c "gofmt -s -w . && goimports -w ."
	@echo "🎨 Formatting JS code..."
	$(COMPOSE_CMD) --profile tools run --rm videogrinder-devtools sh -c "cd web && npm install"
	$(COMPOSE_CMD) --profile tools run --rm videogrinder-devtools sh -c "cd web && npx eslint . --ext .js --fix"
	@echo "✅ Code formatted"

fmt-check: ## Check code formatting without changing files
	@echo "🎨 Checking Go code formatting..."
	$(COMPOSE_CMD) --profile tools run --rm videogrinder-devtools sh -c "test -z \"\$$(gofmt -l .)\" || (echo 'Go files not formatted:' && gofmt -l . && exit 1)"
	$(COMPOSE_CMD) --profile tools run --rm videogrinder-devtools sh -c "test -z \"\$$(goimports -l .)\" || (echo 'Go imports not formatted:' && goimports -l . && exit 1)"
	@echo "🎨 Checking JS code formatting..."
	$(COMPOSE_CMD) --profile tools run --rm videogrinder-devtools sh -c "cd web && npm install"
	$(COMPOSE_CMD) --profile tools run --rm videogrinder-devtools sh -c "cd web && npx eslint . --ext .js"
	@echo "✅ Code formatting is correct"

fmt-js: ## Format JavaScript code
	@echo "🎨 Formatting JS code..."
	$(COMPOSE_CMD) --profile tools run --rm videogrinder-devtools sh -c "cd web && npm install"
	$(COMPOSE_CMD) --profile tools run --rm videogrinder-devtools sh -c "cd web && npx eslint . --ext .js --fix"
	@echo "✅ JS code formatted"

check: fmt lint test test-js ## Run all quality checks

health: ## Check application health (usage: make health [dev|prod])
	@echo "🏥 Checking application health..."
	@echo "🌐 Checking Web Service (port 8080)..."
	$(COMPOSE_CMD) --profile tools run --rm videogrinder-devtools sh -c "curl -f http://host.docker.internal:8080/health || echo '❌ Web Service failed'"
	@echo "🔌 Checking API Service (port 8081)..."
	$(COMPOSE_CMD) --profile tools run --rm videogrinder-devtools sh -c "curl -f http://host.docker.internal:8081/health || echo '❌ API Service failed'"
	@echo "⚙️  Checking Processor Service (port 8082)..."
	$(COMPOSE_CMD) --profile tools run --rm videogrinder-devtools sh -c "curl -f http://host.docker.internal:8082/health || echo '❌ Processor Service failed'"
	@echo "✅ Health check completed!"

logs: ## View all services logs (usage: make logs [dev|prod])
	@echo "📋 Showing all services logs..."
	$(COMPOSE_CMD) logs -f $(WEB_SERVICE) $(API_SERVICE) $(PROCESSOR_SERVICE)

logs-web: ## View Web service logs (usage: make logs-web [dev|prod])
	@echo "📋 Showing Web service logs..."
	$(COMPOSE_CMD) logs -f $(WEB_SERVICE)

logs-api: ## View API service logs (usage: make logs-api [dev|prod])
	@echo "📋 Showing API service logs..."
	$(COMPOSE_CMD) logs -f $(API_SERVICE)

logs-processor: ## View processor service logs (usage: make logs-processor [dev|prod])
	@echo "📋 Showing processor service logs..."
	$(COMPOSE_CMD) logs -f $(PROCESSOR_SERVICE)

down: ## Stop services (usage: make down [dev|prod|all])
	@echo "🐳 Stopping $(ENV) services..."
ifeq ($(ENV),all)
	$(COMPOSE_CMD) down
else
	$(COMPOSE_CMD) stop $(WEB_SERVICE) $(API_SERVICE) $(PROCESSOR_SERVICE)
endif

docker-clean: ## Clean Docker resources
	@echo "🧹 Cleaning Docker resources..."
	$(COMPOSE_CMD) down --volumes --rmi all || true
	docker system prune -f || true
