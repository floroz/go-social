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
migrate-down:
	@go run ./cmd/migrate/seed_db/main.go

.PHONE: test
test:
	@go test -cover ./...