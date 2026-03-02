SHELL := /bin/bash
export PATH := /usr/local/bin:/opt/homebrew/bin:/usr/local/go/bin:$(PATH)

# Resolve docker path explicitly for environments where PATH is limited
DOCKER := $(shell which docker 2>/dev/null || echo /usr/local/bin/docker)

.PHONY: help dev dev-local dev-backend dev-frontend build build-backend build-frontend \
        test test-backend test-frontend lint lint-backend lint-frontend \
        docker-up docker-down docker-infra docker-logs \
        migrate-up migrate-down clean

# ─── Variables ───────────────────────────────────────────
BACKEND_DIR  := backend
FRONTEND_DIR := frontend

# ─── Help ────────────────────────────────────────────────
help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# ─── Development ─────────────────────────────────────────
dev: ## Start full dev stack via Docker (no local Go/Node needed)
	@if [ ! -f .env ]; then \
		echo "📋 Creating .env from .env.example..."; \
		cp .env.example .env; \
	fi
	@echo "🚀 Starting AI Learning Platform (dev mode via Docker)..."
	@echo "   Backend:  http://localhost:8080"
	@echo "   Frontend: http://localhost:3000"
	@echo ""
	$(DOCKER) compose -f docker-compose.dev.yml up --build

dev-local: ## Start dev with local Go & Node (requires Go + Node installed)
	@echo "🚀 Starting infrastructure (postgres + redis)..."
	@$(MAKE) docker-infra
	@echo "✅ Infrastructure ready. Starting backend & frontend..."
	@$(MAKE) -j2 dev-backend dev-frontend

dev-backend: ## Run backend locally in dev mode
	cd $(BACKEND_DIR) && go run ./cmd/server

dev-frontend: ## Run frontend locally in dev mode
	cd $(FRONTEND_DIR) && npm run dev

# ─── Build ───────────────────────────────────────────────
build: build-backend build-frontend ## Build all

build-backend: ## Build backend binary
	cd $(BACKEND_DIR) && go build -o bin/server ./cmd/server

build-frontend: ## Build frontend for production
	cd $(FRONTEND_DIR) && npm run build

# ─── Test ────────────────────────────────────────────────
test: test-backend test-frontend ## Run all tests

test-backend: ## Run backend tests
	cd $(BACKEND_DIR) && go test ./... -v -cover

test-frontend: ## Run frontend tests
	cd $(FRONTEND_DIR) && npm test

# ─── Lint ────────────────────────────────────────────────
lint: lint-backend lint-frontend ## Lint all

lint-backend: ## Lint backend
	cd $(BACKEND_DIR) && golangci-lint run ./...

lint-frontend: ## Lint frontend
	cd $(FRONTEND_DIR) && npm run lint

# ─── Docker ──────────────────────────────────────────────
docker-up: ## Start all Docker services (production build)
	$(DOCKER) compose up -d --build

docker-down: ## Stop all Docker services
	$(DOCKER) compose -f docker-compose.dev.yml down
	$(DOCKER) compose down

docker-infra: ## Start only infrastructure (postgres, redis)
	$(DOCKER) compose -f docker-compose.dev.yml up -d postgres redis
	@echo "⏳ Waiting for postgres & redis to be healthy..."
	@sleep 3

docker-logs: ## Tail Docker logs
	$(DOCKER) compose -f docker-compose.dev.yml logs -f

# ─── Database ────────────────────────────────────────────
migrate-up: ## Run database migrations
	cd $(BACKEND_DIR) && go run ./cmd/server migrate up

migrate-down: ## Rollback last migration
	cd $(BACKEND_DIR) && go run ./cmd/server migrate down

# ─── Clean ───────────────────────────────────────────────
clean: ## Remove build artifacts and Docker volumes
	rm -rf $(BACKEND_DIR)/bin
	rm -rf $(BACKEND_DIR)/tmp
	rm -rf $(FRONTEND_DIR)/.next
	rm -rf $(FRONTEND_DIR)/node_modules
	$(DOCKER) compose -f docker-compose.dev.yml down -v 2>/dev/null || true
