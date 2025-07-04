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
NPM_TEST_CMD = cd web && npm run test
NPM_TEST_WATCH_CMD = cd web && npm run test:watch
NPM_TEST_COVERAGE_CMD = cd web && npm run test:coverage
NPM_TEST_E2E_CMD = cd web && npx cypress install && npm run test:e2e
NPM_TEST_E2E_DEV_CMD = cd web && npx cypress install && npm run test:e2e:dev
NPM_LINT_CMD = cd web && npm run lint:js
NPM_LINT_FIX_CMD = cd web && npm run lint:js:fix
LOGS_TAIL_CMD = logs --tail=50
LOGS_FOLLOW_CMD = logs -f

# Performance Optimization Variables
TOOLS_CONTAINER_NAME = videogrinder-tools-persistent
TOOLS_IMAGE = videogrinder-processor-videogrinder-tools

# Check if tools container is running
TOOLS_RUNNING = $(shell docker ps -q -f name=$(TOOLS_CONTAINER_NAME) 2>/dev/null)

# Check if we have native Go tools installed
HAS_NATIVE_GO = $(shell command -v go 2> /dev/null)
HAS_NATIVE_NODE = $(shell command -v node 2> /dev/null)

# Performance-optimized command execution
define EXEC_TOOLS_CMD
	$(if $(TOOLS_RUNNING), \
		docker exec $(TOOLS_CONTAINER_NAME) sh -c $(1), \
		$(if $(HAS_NATIVE_GO)$(HAS_NATIVE_NODE), \
			$(if $(findstring go,$(1)), \
				$(if $(HAS_NATIVE_GO), \
					$(1), \
					$(COMPOSE_CMD) --profile tools run --rm videogrinder-tools sh -c $(1) \
				), \
				$(if $(findstring npm,$(1)), \
					$(if $(HAS_NATIVE_NODE), \
						$(1), \
						$(COMPOSE_CMD) --profile tools run --rm videogrinder-tools sh -c $(1) \
					), \
					$(COMPOSE_CMD) --profile tools run --rm videogrinder-tools sh -c $(1) \
				) \
			), \
			$(COMPOSE_CMD) --profile tools run --rm videogrinder-tools sh -c $(1) \
		) \
	)
endef

# Tools service commands (optimized)
TOOLS_CMD = $(COMPOSE_CMD) --profile tools run --rm videogrinder-tools sh -c
TOOLS_EXEC_CMD = $(call EXEC_TOOLS_CMD,"$(1)")

# Batch command for multiple operations
define BATCH_TOOLS_CMD
	$(if $(TOOLS_RUNNING), \
		docker exec $(TOOLS_CONTAINER_NAME) sh -c $(1), \
		$(COMPOSE_CMD) --profile tools run --rm videogrinder-tools sh -c $(1) \
	)
endef

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
	@echo '  make tools-start  # Start persistent tools container (faster commands)'
	@echo '  make tools-stop   # Stop persistent tools container'
	@echo '  make tools-status # Show tools container status'
	@echo '  make status       # Show services status'
	@echo '  make ps           # Show running containers'
	@echo ''
	@echo 'Quality Checks:'
	@echo '  make check        # Run all quality checks (format + lint + test)'
	@echo '  make check-fast   # Run all quality checks (optimized single batch - fastest)'
	@echo '  make check-full   # Run quality checks + health check (like CI pipeline)'
	@echo '  make fmt          # Format code (Go + JS)'
	@echo '  make fmt-fast     # Format code (optimized single batch - fastest)'
	@echo '  make lint         # Lint code (Go + JS)'
	@echo '  make lint-fast    # Lint code (optimized single batch - fastest)'
	@echo '  make test         # Run all Go tests'
	@echo '  make test-fast    # Run all tests (optimized single batch - fastest)'
	@echo '  make test-js      # Run JS tests'
	@echo '  make health       # Check app health'
	@echo '  make health-deps  # Check service dependencies status'
	@echo '  make health-json  # Get health status in JSON format'
	@echo ''
	@echo 'Logs & Monitoring:'
	@echo '  make logs         # View all services logs (usage: make logs [dev|prod])'
	@echo '  make logs-tail    # Show last 50 lines of all services logs (usage: make logs-tail [dev|prod])'
	@echo '  make logs-web     # View Web service logs (usage: make logs-web [dev|prod])'
	@echo '  make logs-api     # View API service logs (usage: make logs-api [dev|prod])'
	@echo '  make logs-processor # View processor service logs (usage: make logs-processor [dev|prod])'
	@echo ''
	@echo 'LocalStack (AWS Development):'
	@echo '  make localstack-start # Start LocalStack services'
	@echo '  make localstack-stop  # Stop LocalStack services'
	@echo '  make localstack-init  # Initialize LocalStack resources'
	@echo '  make localstack-status # Check LocalStack status'
	@echo '  make localstack-logs  # View LocalStack logs'
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
	$(call BATCH_TOOLS_CMD,"$(GO_TEST_CMD) ./...")

test-api: ## Run API service unit tests
	@echo "🧪 Running API service unit tests..."
	$(call BATCH_TOOLS_CMD,"$(GO_TEST_CMD) ./api/internal/...")

test-processor: ## Run processor service unit tests
	@echo "🧪 Running processor service unit tests..."
	$(call BATCH_TOOLS_CMD,"$(GO_TEST_CMD) ./processor/internal/...")

test-services: ## Run services unit tests (API + processor)
	@echo "🧪 Running services unit tests..."
	$(call BATCH_TOOLS_CMD,"$(GO_TEST_CMD) ./internal/services/...")

test-utils: ## Run utils unit tests
	@echo "🧪 Running utils unit tests..."
	$(call BATCH_TOOLS_CMD,"$(GO_TEST_CMD) ./internal/utils/...")

test-clients: ## Run clients unit tests
	@echo "🧪 Running clients unit tests..."
	$(call BATCH_TOOLS_CMD,"$(GO_TEST_CMD) ./internal/clients/...")

test-js: ## Run JavaScript unit tests
	@echo "🧪 Running JavaScript unit tests..."
	$(call BATCH_TOOLS_CMD,"$(NPM_TEST_CMD)")

test-js-watch: ## Run unit tests in watch mode
	@echo "🧪 Running unit tests in watch mode..."
	$(call BATCH_TOOLS_CMD,"$(NPM_TEST_WATCH_CMD)")

test-js-coverage: ## Run unit tests with coverage report
	@echo "🧪 Running unit tests with coverage..."
	$(call BATCH_TOOLS_CMD,"$(NPM_TEST_COVERAGE_CMD)")

test-fast: ## Run all tests optimized (Go + JS in single batch)
	@echo "🚀 Running all tests (optimized batch)..."
	$(call BATCH_TOOLS_CMD,"$(GO_TEST_CMD) ./... && $(NPM_TEST_CMD)")

test-e2e: ## Run e2e tests (requires app running)
	@echo "🎭 Running e2e tests..."
	@echo "⚠️  Make sure the app is running with 'make run' in another terminal"
	$(TOOLS_CMD) "$(NPM_TEST_E2E_CMD)"

test-e2e-open: ## Open Cypress interactive mode
	@echo "🎭 Opening Cypress..."
	$(TOOLS_CMD) "$(NPM_TEST_E2E_DEV_CMD)"

test-e2e-local: ## Run e2e tests locally (outside container)
	@echo "🎭 Running e2e tests locally..."
	@echo "⚠️  Make sure the app is running with 'make run' in another terminal"
	@echo "⚠️  Make sure you have Node.js and npm installed locally"
	cd web && npm run test:e2e

test-e2e-open-local: ## Open Cypress interactive mode locally (outside container)
	@echo "🎭 Opening Cypress locally..."
	@echo "⚠️  Make sure you have Node.js and npm installed locally"
	cd web && npm run test:e2e:dev

lint: ## Check code quality (Go + JS)
	@echo "🔍 Running Go + JS linters (optimized batch)..."
	$(call BATCH_TOOLS_CMD,"GOFLAGS='-buildvcs=false' golangci-lint run && $(NPM_LINT_CMD)")

lint-fast: ## Check code quality (Go + JS in single batch - fastest)
	@echo "🚀 Running all linters (optimized)..."
	$(call BATCH_TOOLS_CMD,"GOFLAGS='-buildvcs=false' golangci-lint run && $(NPM_LINT_CMD)")

lint-js: ## Check JavaScript code quality
	@echo "🔍 Running JS linters..."
	$(call BATCH_TOOLS_CMD,"$(NPM_LINT_CMD)")

fmt: ## Format code (Go + JS)
	@echo "🎨 Formatting all code (optimized batch)..."
	$(call BATCH_TOOLS_CMD,"gofmt -s -w . && goimports -w . && $(NPM_LINT_FIX_CMD)")
	@echo "✅ Code formatted"

fmt-fast: ## Format code (Go + JS in single batch - fastest)
	@echo "🚀 Formatting all code (optimized)..."
	$(call BATCH_TOOLS_CMD,"gofmt -s -w . && goimports -w . && $(NPM_LINT_FIX_CMD)")
	@echo "✅ Code formatted"

fmt-check: ## Check code formatting without changing files
	@echo "🎨 Checking code formatting (optimized batch)..."
	$(call BATCH_TOOLS_CMD,"test -z \"\$$(gofmt -l .)\" || (echo 'Go files not formatted:' && gofmt -l . && exit 1) && test -z \"\$$(goimports -l .)\" || (echo 'Go imports not formatted:' && goimports -l . && exit 1) && $(NPM_LINT_CMD)")
	@echo "✅ Code formatting is correct"

fmt-js: ## Format JavaScript code
	@echo "🎨 Formatting JS code..."
	$(call BATCH_TOOLS_CMD,"$(NPM_LINT_FIX_CMD)")
	@echo "✅ JS code formatted"

check: fmt-check lint test test-js ## Run all quality checks

check-fast: ## Run all quality checks (optimized single batch - fastest)
	@echo "🚀 Running all quality checks (optimized)..."
	$(call BATCH_TOOLS_CMD,"test -z \"\$$(gofmt -l .)\" || (echo 'Go files not formatted:' && gofmt -l . && exit 1) && test -z \"\$$(goimports -l .)\" || (echo 'Go imports not formatted:' && goimports -l . && exit 1) && $(NPM_LINT_CMD) && GOFLAGS='-buildvcs=false' golangci-lint run && $(GO_TEST_CMD) ./... && $(NPM_TEST_CMD)")
	@echo "✅ All quality checks passed!"

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
	@echo "=== 📊 VIDEOGRINDER HEALTH STATUS ==="
	@echo ""
	@echo "🌐 Web Service:"
	@if curl -s -f http://localhost:8080/health > /dev/null 2>&1; then \
		curl -s http://localhost:8080/health | jq -r '"  Status: " + .status' 2>/dev/null || echo "  Status: healthy"; \
	else \
		echo "  Status: ❌ FAILED"; \
		echo "💡 Dica: rode 'make logs-tail' para ver os logs."; \
		exit 1; \
	fi
	@echo ""
	@echo "🔗 API Service:"
	@if curl -s -f http://localhost:8081/health > /dev/null 2>&1; then \
		curl -s http://localhost:8081/health | jq -r '"  Overall: " + .status' 2>/dev/null || echo "  Overall: healthy"; \
		curl -s http://localhost:8081/health | jq -r '.checks | to_entries[] | "  " + (if .value.status == "healthy" then "✅" else "❌" end) + " " + .key + ": " + .value.status' 2>/dev/null || echo "  Dependencies: OK"; \
	else \
		echo "  Overall: ❌ FAILED"; \
		echo "💡 Dica: rode 'make logs-tail' para ver os logs."; \
		exit 1; \
	fi
	@echo ""
	@echo "⚙️ Processor Service:"
	@if curl -s -f http://localhost:8082/health > /dev/null 2>&1; then \
		curl -s http://localhost:8082/health | jq -r '"  Overall: " + .status' 2>/dev/null || echo "  Overall: healthy"; \
		curl -s http://localhost:8082/health | jq -r '.checks | to_entries[] | "  " + (if .value.status == "healthy" then "✅" else "❌" end) + " " + .key + ": " + .value.status' 2>/dev/null || echo "  Dependencies: OK"; \
	else \
		echo "  Overall: ❌ FAILED"; \
		echo "💡 Dica: rode 'make logs-tail' para ver os logs."; \
		exit 1; \
	fi
	@echo ""
	@echo "✅ All services are healthy!"

health-deps: ## Check service dependencies status (usage: make health-deps [dev|prod])
	@echo "=== 🔍 DEPENDENCIES STATUS ==="
	@echo ""
	@echo "🔗 API Dependencies:"
	@curl -s http://localhost:8081/health | jq -r '.checks | to_entries[] | "  " + (if .value.status == "healthy" then "✅" else "❌" end) + " " + .key + ": " + .value.status' 2>/dev/null || echo "  ❌ API Service not available"
	@echo ""
	@echo "⚙️ Processor Dependencies:"
	@curl -s http://localhost:8082/health | jq -r '.checks | to_entries[] | "  " + (if .value.status == "healthy" then "✅" else "❌" end) + " " + .key + ": " + .value.status' 2>/dev/null || echo "  ❌ Processor Service not available"
	@echo ""

health-json: ## Get health status in JSON format (usage: make health-json [dev|prod])
	@echo "{"
	@echo "  \"web\": $(shell curl -s http://localhost:8080/health 2>/dev/null || echo 'null'),"
	@echo "  \"api\": $(shell curl -s http://localhost:8081/health 2>/dev/null || echo 'null'),"
	@echo "  \"processor\": $(shell curl -s http://localhost:8082/health 2>/dev/null || echo 'null')"
	@echo "}"

help-health: ## Show detailed help for health commands
	@echo "=== 🏥 HEALTH CHECK COMMANDS ==="
	@echo ""
	@echo "📊 make health"
	@echo "   Complete health check of all services with detailed dependency status"
	@echo "   Shows: Web, API, and Processor services with their dependencies"
	@echo "   Dependencies checked:"
	@echo "     • API: directories, processor connectivity, S3 connectivity"
	@echo "     • Processor: directories, FFmpeg availability, S3 connectivity"
	@echo ""
	@echo "🔍 make health-deps"
	@echo "   Quick check of service dependencies only"
	@echo "   Shows only the status of each dependency without service details"
	@echo ""
	@echo "📄 make health-json"
	@echo "   Get complete health status in JSON format"
	@echo "   Useful for monitoring tools and automated scripts"
	@echo ""
	@echo "💡 Tips:"
	@echo "   • Use 'make health' for human-readable status"
	@echo "   • Use 'make health-deps' for quick dependency overview"
	@echo "   • Use 'make health-json' for programmatic access"
	@echo "   • All commands work with [dev|prod] environments"
	@echo ""

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
	@echo "🛑 Stopping persistent tools container if running..."
	@if [ -n "$(TOOLS_RUNNING)" ]; then \
		docker stop $(TOOLS_CONTAINER_NAME) && docker rm $(TOOLS_CONTAINER_NAME) || true; \
	fi
else
	$(COMPOSE_CMD) --profile $(PROFILE) down --volumes --remove-orphans
endif

docker-clean: ## Clean Docker resources
	@echo "🧹 Cleaning Docker resources..."
	@echo "📦 Stopping and removing compose resources..."
	$(COMPOSE_CMD) down --volumes --rmi all || true
	@echo "🗑️  Removing project-specific volumes..."
	docker volume rm videogrinder-processor_videogrinder-uploads videogrinder-processor_videogrinder-outputs videogrinder-processor_videogrinder-temp videogrinder-processor_air-tmp videogrinder-processor_localstack-data 2>/dev/null || true
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

tools-start: ## Start persistent tools container (faster commands)
	@echo "🚀 Starting persistent tools container..."
	@if [ -z "$(TOOLS_RUNNING)" ]; then \
		echo "⚡ Starting $(TOOLS_CONTAINER_NAME) for faster development..."; \
		docker run -d --name $(TOOLS_CONTAINER_NAME) \
			-v "$(PWD):/workspace" \
			-w /workspace \
			$(TOOLS_IMAGE) tail -f /dev/null; \
		echo "✅ Tools container started! Commands will now be much faster."; \
	else \
		echo "⚡ Tools container already running!"; \
	fi

tools-stop: ## Stop persistent tools container
	@echo "🛑 Stopping persistent tools container..."
	@if [ -n "$(TOOLS_RUNNING)" ]; then \
		docker stop $(TOOLS_CONTAINER_NAME) && docker rm $(TOOLS_CONTAINER_NAME); \
		echo "✅ Tools container stopped."; \
	else \
		echo "ℹ️  Tools container not running."; \
	fi

tools-status: ## Show tools container status
	@echo "📊 Tools Container Status:"
	@if [ -n "$(TOOLS_RUNNING)" ]; then \
		echo "✅ Persistent tools container is RUNNING"; \
		echo "   Container: $(TOOLS_CONTAINER_NAME)"; \
		echo "   Commands will execute via docker exec (fast)"; \
	else \
		echo "❌ Persistent tools container is STOPPED"; \
		echo "   Commands will use docker run (slower)"; \
		echo "💡 Run 'make tools-start' for faster development"; \
	fi
	@echo ""
	@echo "🛠️  Native Tools Available:"
	@if [ -n "$(HAS_NATIVE_GO)" ]; then \
		echo "   ✅ Go: $(shell go version 2>/dev/null || echo 'not available')"; \
	else \
		echo "   ❌ Go: not available"; \
	fi
	@if [ -n "$(HAS_NATIVE_NODE)" ]; then \
		echo "   ✅ Node: $(shell node --version 2>/dev/null || echo 'not available')"; \
	else \
		echo "   ❌ Node: not available"; \
	fi

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

# LocalStack commands
localstack-start: ## Start LocalStack services
	@echo "🚀 Starting LocalStack..."
	$(COMPOSE_CMD) --profile localstack up -d localstack
	@echo "⏳ Waiting for LocalStack to be ready..."
	@sleep 10
	@echo "✅ LocalStack started! Available at http://127.0.0.1:4566"

localstack-stop: ## Stop LocalStack services
	@echo "🛑 Stopping LocalStack..."
	$(COMPOSE_CMD) stop localstack
	@echo "✅ LocalStack stopped"

localstack-init: ## Initialize LocalStack resources (S3, DynamoDB, SQS)
	@echo "🔧 Initializing LocalStack resources..."
	@if ! docker ps --filter "name=localstack" --format "table {{.Names}}" | grep -q localstack; then \
		echo "❌ LocalStack is not running. Starting it first..."; \
		make localstack-start; \
	fi
	@echo "📦 Running initialization script..."
	./localstack-init.sh
	@echo "✅ LocalStack resources initialized!"

localstack-status: ## Check LocalStack status and resources
	@echo "📊 LocalStack Status:"
	@if docker ps --filter "name=localstack" --format "table {{.Names}}" | grep -q localstack; then \
		echo "✅ LocalStack container: running"; \
		echo "🔗 Health check:"; \
		curl -s http://127.0.0.1:4566/health | jq . 2>/dev/null || curl -s http://127.0.0.1:4566/health; \
		echo ""; \
		echo "📦 S3 Buckets:"; \
		AWS_ACCESS_KEY_ID=test AWS_SECRET_ACCESS_KEY=test AWS_DEFAULT_REGION=us-east-1 aws s3 ls --endpoint-url=http://127.0.0.1:4566 2>/dev/null || echo "  No buckets found"; \
		echo "🗃️ DynamoDB Tables:"; \
		AWS_ACCESS_KEY_ID=test AWS_SECRET_ACCESS_KEY=test AWS_DEFAULT_REGION=us-east-1 aws dynamodb list-tables --endpoint-url=http://127.0.0.1:4566 --output text --query 'TableNames[*]' 2>/dev/null | tr '\t' '\n' | sed 's/^/  /' || echo "  No tables found"; \
		echo "📬 SQS Queues:"; \
		AWS_ACCESS_KEY_ID=test AWS_SECRET_ACCESS_KEY=test AWS_DEFAULT_REGION=us-east-1 aws sqs list-queues --endpoint-url=http://127.0.0.1:4566 --output text --query 'QueueUrls[*]' 2>/dev/null | sed 's|.*/||' | sed 's/^/  /' || echo "  No queues found"; \
	else \
		echo "❌ LocalStack container: not running"; \
		echo "💡 Run 'make localstack-start' to start LocalStack"; \
	fi

localstack-logs: ## View LocalStack logs
	@echo "📋 LocalStack Logs:"
	$(COMPOSE_CMD) logs localstack

localstack-reset: ## Reset LocalStack (stop, remove data, start fresh)
	@echo "🔄 Resetting LocalStack..."
	$(COMPOSE_CMD) stop localstack
	$(COMPOSE_CMD) rm -f localstack
	docker volume rm videogrinder-processor_localstack-data 2>/dev/null || true
	@echo "✅ LocalStack reset completed. Run 'make localstack-start' to start fresh."
