.PHONY: help dev dev-backend dev-frontend build test lint clean docker-up docker-down migrate

# ─── Variables ───────────────────────────────────────────
BACKEND_DIR  := backend
FRONTEND_DIR := frontend

# ─── Help ────────────────────────────────────────────────
help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# ─── Development ─────────────────────────────────────────
dev: ## Run backend + frontend concurrently
	@echo "Starting all services..."
	$(MAKE) docker-infra
	$(MAKE) -j2 dev-backend dev-frontend

dev-backend: ## Run backend in dev mode
	cd $(BACKEND_DIR) && go run ./cmd/server

dev-frontend: ## Run frontend in dev mode
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
docker-up: ## Start all Docker services
	docker compose up -d --build

docker-down: ## Stop all Docker services
	docker compose down

docker-infra: ## Start only infrastructure (postgres, redis)
	docker compose up -d postgres redis

docker-logs: ## Tail Docker logs
	docker compose logs -f

# ─── Database ────────────────────────────────────────────
migrate-up: ## Run database migrations
	cd $(BACKEND_DIR) && go run ./cmd/server migrate up

migrate-down: ## Rollback last migration
	cd $(BACKEND_DIR) && go run ./cmd/server migrate down

# ─── Clean ───────────────────────────────────────────────
clean: ## Remove build artifacts
	rm -rf $(BACKEND_DIR)/bin
	rm -rf $(FRONTEND_DIR)/.next
	rm -rf $(FRONTEND_DIR)/node_modules
