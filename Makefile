.PHONY: help setup run run-api run-web run-processor test test-api test-processor test-services test-utils test-clients test-js test-js-watch test-js-coverage test-e2e test-e2e-open lint lint-js fmt fmt-js fmt-check check check-full logs logs-web logs-api logs-processor down docker-clean health shell restart restart-api restart-web restart-processor build rebuild status ps

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

# Variables for common commands
GO_TEST_CMD = GOFLAGS='-buildvcs=false' go test -v
NPM_INSTALL_CMD = cd web && npm install
NPM_TEST_CMD = cd web && npm test
NPM_LINT_CMD = cd web && npm run lint:js
NPM_LINT_FIX_CMD = cd web && npm run lint:js:fix
LOGS_TAIL_CMD = logs --tail=50
LOGS_FOLLOW_CMD = logs -f

# Tools service commands
TOOLS_CMD = $(COMPOSE_CMD) --profile tools run --rm videogrinder-tools sh -c

%:
	@:

help: ## Show this help message
	@echo '🚀 VideoGrinder Processor - Available Commands:'
	@echo ''
	@echo 'Environment Setup:'
	@echo '  make setup        # Configure environment (usage: make setup [dev|prod])'
	@echo '  make build        # Build all services (usage: make build [dev|prod])'
	@echo '  make rebuild      # Rebuild all services (usage: make rebuild [dev|prod])'
	@echo ''
	@echo 'Service Management:'
	@echo '  make run          # Run all 3 services (Web + API + Processor) (usage: make run [dev|prod])'
	@echo '  make run-api      # Run only API service (usage: make run-api [dev|prod])'
	@echo '  make run-web      # Run only web service (usage: make run-web [dev|prod])'
	@echo '  make run-processor # Run only processor service (usage: make run-processor [dev|prod])'
	@echo '  make restart      # Restart all services (usage: make restart [dev|prod])'
	@echo '  make restart-api  # Restart API service (usage: make restart-api [dev|prod])'
	@echo '  make restart-web  # Restart web service (usage: make restart-web [dev|prod])'
	@echo '  make restart-processor # Restart processor service (usage: make restart-processor [dev|prod])'
	@echo '  make down         # Stop services (usage: make down [dev|prod|all])'
	@echo ''
	@echo 'Development:'
	@echo '  make shell        # Open shell in tools container'
	@echo '  make status       # Show services status'
	@echo '  make ps           # Show running containers'
	@echo ''
	@echo 'Quality Checks:'
	@echo '  make check        # Run all quality checks (format + lint + test)'
	@echo '  make check-full   # Run quality checks + health check (like CI pipeline)'
	@echo '  make fmt          # Format code (Go + JS)'
	@echo '  make lint         # Lint code (Go + JS)'
	@echo '  make test         # Run all Go tests'
	@echo '  make test-js      # Run JS tests'
	@echo '  make health       # Check app health'
	@echo ''
	@echo 'Logs & Monitoring:'
	@echo '  make logs         # View all services logs (usage: make logs [dev|prod])'
	@echo '  make logs-tail    # Show last 50 lines of all services logs (usage: make logs-tail [dev|prod])'
	@echo '  make logs-web     # View Web service logs (usage: make logs-web [dev|prod])'
	@echo '  make logs-api     # View API service logs (usage: make logs-api [dev|prod])'
	@echo '  make logs-processor # View processor service logs (usage: make logs-processor [dev|prod])'
	@echo ''
	@echo 'Maintenance:'
	@echo '  make docker-clean # Clean Docker resources'

	@echo ''
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

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
	$(TOOLS_CMD) "$(GO_TEST_CMD) ./..."

test-api: ## Run API service unit tests
	@echo "🧪 Running API service unit tests..."
	$(TOOLS_CMD) "$(GO_TEST_CMD) ./api/internal/..."

test-processor: ## Run processor service unit tests
	@echo "🧪 Running processor service unit tests..."
	$(TOOLS_CMD) "$(GO_TEST_CMD) ./processor/internal/..."

test-services: ## Run services unit tests (API + processor)
	@echo "🧪 Running services unit tests..."
	$(TOOLS_CMD) "$(GO_TEST_CMD) ./internal/services/..."

test-utils: ## Run utils unit tests
	@echo "🧪 Running utils unit tests..."
	$(TOOLS_CMD) "$(GO_TEST_CMD) ./internal/utils/..."

test-clients: ## Run clients unit tests
	@echo "🧪 Running clients unit tests..."
	$(TOOLS_CMD) "$(GO_TEST_CMD) ./internal/clients/..."

test-js: ## Run JavaScript unit tests
	@echo "🧪 Running JavaScript unit tests..."
	$(TOOLS_CMD) "$(NPM_TEST_CMD)"

test-js-watch: ## Run unit tests in watch mode
	@echo "🧪 Running unit tests in watch mode..."
	$(TOOLS_CMD) "cd web && npm run test:watch"

test-js-coverage: ## Run unit tests with coverage report
	@echo "🧪 Running unit tests with coverage..."
	$(TOOLS_CMD) "cd web && npm run test:coverage"

test-e2e: ## Run e2e tests (requires app running)
	@echo "🎭 Running e2e tests..."
	@echo "⚠️  Make sure the app is running with 'make run' in another terminal"
	$(TOOLS_CMD) "cd web && npm install cypress --save-dev && npx cypress install"
	$(TOOLS_CMD) "cd web && npx cypress run"

test-e2e-open: ## Open Cypress interactive mode
	@echo "🎭 Opening Cypress..."
	$(TOOLS_CMD) "cd web && npm install cypress --save-dev && npx cypress install"
	$(TOOLS_CMD) "cd web && npx cypress open"

lint: ## Check code quality (Go + JS)
	@echo "🔍 Running Go linters..."
	$(TOOLS_CMD) "GOFLAGS='-buildvcs=false' golangci-lint run"
	@echo "🔍 Running JS linters..."
	$(TOOLS_CMD) "$(NPM_LINT_CMD)"

lint-js: ## Check JavaScript code quality
	@echo "🔍 Running JS linters..."
	$(TOOLS_CMD) "$(NPM_LINT_CMD)"

fmt: ## Format code (Go + JS)
	@echo "🎨 Formatting Go code..."
	$(TOOLS_CMD) "gofmt -s -w . && goimports -w ."
	@echo "🎨 Formatting JS code..."
	$(TOOLS_CMD) "$(NPM_LINT_FIX_CMD)"
	@echo "✅ Code formatted"

fmt-check: ## Check code formatting without changing files
	@echo "🎨 Checking Go code formatting..."
	$(TOOLS_CMD) "test -z \"\$$(gofmt -l .)\" || (echo 'Go files not formatted:' && gofmt -l . && exit 1)"
	$(TOOLS_CMD) "test -z \"\$$(goimports -l .)\" || (echo 'Go imports not formatted:' && goimports -l . && exit 1)"
	@echo "🎨 Checking JS code formatting..."
	$(TOOLS_CMD) "$(NPM_LINT_CMD)"
	@echo "✅ Code formatting is correct"

fmt-js: ## Format JavaScript code
	@echo "🎨 Formatting JS code..."
	$(TOOLS_CMD) "$(NPM_LINT_FIX_CMD)"
	@echo "✅ JS code formatted"

check: fmt-check lint test test-js ## Run all quality checks

check-full: check ## Run all quality checks + health check (like CI pipeline)
	@echo "🏥 Running health check (Step 5 from CI pipeline)..."
	@echo "🚀 Starting services for health check..."
	$(COMPOSE_CMD) --profile $(PROFILE) up -d --build $(WEB_SERVICE) $(API_SERVICE) $(PROCESSOR_SERVICE)
	@echo "⏳ Waiting for services to start..."
	@sleep 5
	@echo "🔍 Running health checks..."
	@if ! timeout 60s bash -c 'until make health; do sleep 2; done'; then \
		echo "❌ Health check failed. Showing service logs..."; \
		make logs-tail; \
		make down; \
		exit 1; \
	fi
	@echo "✅ Health check passed!"
	@echo "🧹 Stopping services..."
	$(COMPOSE_CMD) stop $(WEB_SERVICE) $(API_SERVICE) $(PROCESSOR_SERVICE)
	@echo "✅ Full check completed successfully!"

health: ## Check application health (usage: make health [dev|prod])
	@echo "🏥 Checking application health..."
	@echo "🌐 Checking Web Service (port 8080)..."
	@if curl -s -f http://localhost:8080/health > /dev/null 2>&1; then \
		echo "✅ Web Service: healthy"; \
	else \
		echo "❌ Web Service: failed"; \
		echo "💡 Dica: rode 'make logs-tail' para ver os logs."; \
		exit 1; \
	fi
	@echo "🔌 Checking API Service (port 8081)..."
	@if curl -s -f http://localhost:8081/health > /dev/null 2>&1; then \
		echo "✅ API Service: healthy"; \
	else \
		echo "❌ API Service: failed"; \
		echo "💡 Dica: rode 'make logs-tail' para ver os logs."; \
		exit 1; \
	fi
	@echo "⚙️  Checking Processor Service (port 8082)..."
	@if curl -s -f http://localhost:8082/health > /dev/null 2>&1; then \
		echo "✅ Processor Service: healthy"; \
	else \
		echo "❌ Processor Service: failed"; \
		echo "💡 Dica: rode 'make logs-tail' para ver os logs."; \
		exit 1; \
	fi
	@echo "✅ All services are healthy!"

logs: ## View all services logs (usage: make logs [dev|prod])
	@echo "📋 Showing all services logs..."
	$(COMPOSE_CMD) $(LOGS_FOLLOW_CMD) $(WEB_SERVICE) $(API_SERVICE) $(PROCESSOR_SERVICE)

logs-tail: ## Show last 50 lines of all services logs (usage: make logs-tail [dev|prod])
	@echo "📋 Showing last 50 lines of all services logs..."
	$(COMPOSE_CMD) $(LOGS_TAIL_CMD) $(WEB_SERVICE) $(API_SERVICE) $(PROCESSOR_SERVICE)

logs-web: ## View Web service logs (usage: make logs-web [dev|prod])
	@echo "📋 Showing Web service logs..."
	$(COMPOSE_CMD) $(LOGS_FOLLOW_CMD) $(WEB_SERVICE)

logs-web-tail: ## Show last 30 lines of Web service logs (usage: make logs-web-tail [dev|prod])
	@echo "📋 Showing last 30 lines of Web service logs..."
	$(COMPOSE_CMD) logs --tail=30 $(WEB_SERVICE)

logs-api: ## View API service logs (usage: make logs-api [dev|prod])
	@echo "📋 Showing API service logs..."
	$(COMPOSE_CMD) $(LOGS_FOLLOW_CMD) $(API_SERVICE)

logs-api-tail: ## Show last 30 lines of API service logs (usage: make logs-api-tail [dev|prod])
	@echo "📋 Showing last 30 lines of API service logs..."
	$(COMPOSE_CMD) logs --tail=30 $(API_SERVICE)

logs-processor: ## View processor service logs (usage: make logs-processor [dev|prod])
	@echo "📋 Showing processor service logs..."
	$(COMPOSE_CMD) $(LOGS_FOLLOW_CMD) $(PROCESSOR_SERVICE)

logs-processor-tail: ## Show last 30 lines of processor service logs (usage: make logs-processor-tail [dev|prod])
	@echo "📋 Showing last 30 lines of processor service logs..."
	$(COMPOSE_CMD) logs --tail=30 $(PROCESSOR_SERVICE)

down: ## Stop services (usage: make down [dev|prod|all])
	@echo "🐳 Stopping $(ENV) services..."
ifeq ($(ENV),all)
	$(COMPOSE_CMD) --profile dev down --volumes --remove-orphans
	$(COMPOSE_CMD) --profile prod down --volumes --remove-orphans
	$(COMPOSE_CMD) --profile tools down --volumes --remove-orphans
	$(COMPOSE_CMD) down --volumes --remove-orphans
else
	$(COMPOSE_CMD) --profile $(PROFILE) down --volumes --remove-orphans
endif

docker-clean: ## Clean Docker resources
	@echo "🧹 Cleaning Docker resources..."
	@echo "📦 Stopping and removing compose resources..."
	$(COMPOSE_CMD) down --volumes --rmi all || true
	@echo "🗑️  Removing project-specific volumes..."
	docker volume rm videogrinder-processor_videogrinder-uploads videogrinder-processor_videogrinder-outputs videogrinder-processor_videogrinder-temp videogrinder-processor_air-tmp 2>/dev/null || true
	@echo "🧽 Cleaning unused Docker resources..."
	docker system prune -f || true
	docker volume prune -f || true
	docker container prune -f || true
	docker network prune -f || true
	docker builder prune -f || true
	@echo "✅ Docker cleanup completed!"

shell: ## Open shell in tools container
	@echo "🐚 Opening shell in tools container..."
	$(COMPOSE_CMD) --profile tools run --rm videogrinder-tools sh

restart: ## Restart all services (usage: make restart [dev|prod])
	@echo "🔄 Restarting all services in $(ENV) mode..."
	$(COMPOSE_CMD) restart $(WEB_SERVICE) $(API_SERVICE) $(PROCESSOR_SERVICE)
	@echo "✅ Services restarted"

restart-api: ## Restart API service (usage: make restart-api [dev|prod])
	@echo "🔄 Restarting API service in $(ENV) mode..."
	$(COMPOSE_CMD) restart $(API_SERVICE)
	@echo "✅ API service restarted"

restart-web: ## Restart web service (usage: make restart-web [dev|prod])
	@echo "🔄 Restarting web service in $(ENV) mode..."
	$(COMPOSE_CMD) restart $(WEB_SERVICE)
	@echo "✅ Web service restarted"

restart-processor: ## Restart processor service (usage: make restart-processor [dev|prod])
	@echo "🔄 Restarting processor service in $(ENV) mode..."
	$(COMPOSE_CMD) restart $(PROCESSOR_SERVICE)
	@echo "✅ Processor service restarted"

build: ## Build all services (usage: make build [dev|prod])
	@echo "🔨 Building all services in $(ENV) mode..."
	$(COMPOSE_CMD) build $(WEB_SERVICE) $(API_SERVICE) $(PROCESSOR_SERVICE)
	@echo "✅ All services built"

rebuild: ## Rebuild all services (usage: make rebuild [dev|prod])
	@echo "🔨 Rebuilding all services in $(ENV) mode..."
	$(COMPOSE_CMD) build --no-cache $(WEB_SERVICE) $(API_SERVICE) $(PROCESSOR_SERVICE)
	@echo "✅ All services rebuilt"

status: ## Show services status
	@echo "📊 Services Status:"
	$(COMPOSE_CMD) ps

ps: ## Show running containers
	@echo "🐳 Running Containers:"
	docker ps --filter "name=videogrinder" --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"
