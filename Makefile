# Load environment variables from .env.local
include .env.local

MIGRATIONS_PATH = $(shell pwd)/cmd/migrate/migrations
DATABASE_URL = postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST)/$(DB_NAME)?sslmode=disable

.PHONY: migrate-create
migrate-create:
	@migrate create -ext sql -dir $(MIGRATIONS_PATH) -seq $(name)
	
.PHONY: migrate-up
migrate-up:
	@migrate -database "$(DATABASE_URL)" -path $(MIGRATIONS_PATH) up

.PHONY: migrate-down
migrate-down:
	@migrate -database "$(DATABASE_URL)" -path $(MIGRATIONS_PATH) down

.PHONY: migrate-seed
migrate-seed:
	@go run ./cmd/migrate/seed_db/main.go

.PHONY: build
build: ## Build the Go application
	@echo "Building Go application..."
	@go build ./...

.PHONY: lint
lint: ## Run Go linter (golangci-lint)
	@echo "Running Go linter..."
	@# Ensure golangci-lint is installed or handle installation if needed
	@golangci-lint run ./...

.PHONY: test
test: ## Run Go backend tests
	@go test -cover ./...

.PHONY: test-fe
test-fe: ## Run frontend tests
	@echo "Running frontend tests..."
	@cd frontend && npm run test

.PHONY: test-all
test-all: test test-fe ## Run all backend and frontend tests
	@echo "All tests completed."

.PHONY: dev-be
dev-be: ## Start Go backend server (using go run)
	@echo "Starting Go backend server..."
	@go run ./cmd/main.go

# Optional: Use air for live reload if installed
# .PHONY: dev-be-air
# dev-be-air:
# 	@echo "Starting Go backend server with air..."
# 	@air

.PHONY: dev-fe
dev-fe: ## Start frontend dev server
	@echo "Starting frontend dev server..."
	@cd frontend && npm run dev

.PHONY: dev
dev: ## Start both backend (background) and frontend (foreground) dev servers
	@echo "Starting development servers..."
	@echo "Starting backend in background..."
	@go run ./cmd/main.go & \
	sleep 2 && \
	echo "Starting frontend in foreground..." && \
	cd frontend && npm run dev
	# Note: Killing the frontend process (Ctrl+C) might not automatically kill the background backend process.
	# Consider using a process manager like 'overmind' or 'foreman' for better control in the future.
