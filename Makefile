# ==============================================================================
# Yubi App (OSS) - Makefile
# ==============================================================================
# All commands run inside Docker containers.
# Run `make up` first to start the services.

BE_EXEC := docker compose exec backend
FE_EXEC := docker compose exec frontend

# ==============================================================================
# Help
# ==============================================================================
.PHONY: help
help: ## Show this help message
	@echo ""
	@echo "Yubi App (OSS)"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' Makefile | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'
	@echo ""

# ==============================================================================
# Docker Compose
# ==============================================================================
PLATFORM ?= amd64

.PHONY: up
up: ## Start all services (PLATFORM=amd64|arm64)
	DOCKER_PLATFORM=linux/$(PLATFORM) docker compose up -d --build

.PHONY: down
down: ## Stop all services
	docker compose down

.PHONY: reset
reset: ## Stop all services and delete volumes (DB data)
	docker compose down -v

.PHONY: ps
ps: ## Show status of all containers
	docker compose ps

.PHONY: logs
logs: ## Show logs from all services
	docker compose logs -f

# ==============================================================================
# Database
# ==============================================================================
.PHONY: migrate
MIGRATE_DIR := file://internal/infra/database/migrate
DB_URL_EXPR := postgres://$$DB_USER:$$DB_PASSWORD@$$DB_HOST:$$DB_PORT/$$DB_NAME?sslmode=disable

migrate: ## Apply database migrations
	$(BE_EXEC) sh -c 'atlas migrate apply --dir "$(MIGRATE_DIR)" --url "$(DB_URL_EXPR)"'

.PHONY: migrate-status
migrate-status: ## Show database migration status
	$(BE_EXEC) sh -c 'atlas migrate status --dir "$(MIGRATE_DIR)" --url "$(DB_URL_EXPR)"'

.PHONY: seed
seed: ## Seed database with initial data
	cat backend/internal/infra/database/seeder/initial_data.sql | docker compose exec -T postgres psql -U $${DB_USER:-postgres} -d $${DB_NAME:-airoa}

# ==============================================================================
# Backend
# ==============================================================================
.PHONY: be-test
be-test: ## Run backend tests
	$(BE_EXEC) go test -v ./...

.PHONY: be-lint
be-lint: ## Run backend linter
	$(BE_EXEC) sh -c 'go list ./... | grep -v /gen/ | xargs staticcheck -f stylish'

.PHONY: be-fmt
be-fmt: ## Format backend code
	$(BE_EXEC) go fmt ./...

.PHONY: be-tidy
be-tidy: ## Tidy Go modules
	$(BE_EXEC) go mod tidy

.PHONY: be-schema-gen
be-schema-gen: ## Generate schema.up.sql from Go entities
	$(BE_EXEC) go run ./cmd/create-db-schema/main.go

.PHONY: be-migrate-diff
be-migrate-diff: be-schema-gen ## Generate migration diff (usage: make be-migrate-diff NAME=xxx)
	$(BE_EXEC) atlas migrate diff $(NAME) --env dev

.PHONY: be-generate-api
be-generate-api: ## Generate Go server code from OpenAPI spec
	$(BE_EXEC) sh -c 'mkdir -p /app/internal/gen/openapi && oapi-codegen -config /app/openapi.yaml /openapi/openapi.yaml'

# ==============================================================================
# Dashboard / Batch
# ==============================================================================
PERIOD ?= hourly

.PHONY: be-aggregate
be-aggregate: ## Aggregate dashboard stats for the previous period (PERIOD=hourly|daily|monthly)
	$(BE_EXEC) go run ./cmd/aggregate-episode-stats/main.go --period $(PERIOD)

.PHONY: be-aggregate-backfill
be-aggregate-backfill: ## Backfill dashboard stats (PERIOD= FROM= TO=, e.g. FROM=2025-11-01 TO=2026-06-01)
	$(BE_EXEC) go run ./cmd/aggregate-episode-stats/main.go --period $(PERIOD) --backfill --from $(FROM) --to $(TO)

.PHONY: be-uptime-writer
be-uptime-writer: ## Run the robot uptime metrics writer (long-running daemon; Ctrl-C to stop)
	$(BE_EXEC) go run ./cmd/write-robot-status-metrics/main.go

# ==============================================================================
# Frontend
# ==============================================================================
.PHONY: fe-install
fe-install: ## Install frontend dependencies
	$(FE_EXEC) npm ci

.PHONY: fe-lint
fe-lint: ## Run frontend linter
	$(FE_EXEC) npm run lint

.PHONY: fe-fmt
fe-fmt: ## Format frontend code
	$(FE_EXEC) npm run format

.PHONY: fe-typecheck
fe-typecheck: ## Run TypeScript typecheck
	$(FE_EXEC) npx tsc --noEmit

.PHONY: fe-ci
fe-ci: ## Run all frontend CI checks (lint, import boundaries, format, typecheck, build)
	$(FE_EXEC) sh -c 'npm run lint && npm run lint:boundaries && npm run format:check && npx tsc --noEmit && npm run build'

.PHONY: fe-generate-api
fe-generate-api: ## Generate frontend API client from OpenAPI spec
	$(FE_EXEC) npm run generate:api
