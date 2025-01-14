include .env

GOOSE_CMD=goose -dir ./migrations postgres $(DATABASE_URL)

.PHONY: migrate-up migrate-down migrate-status migrate-create psql-dev install-tools sqlc db-dump

migrate-up:
	$(GOOSE_CMD) up

migrate-down:
	$(GOOSE_CMD) down

migrate-status:
	$(GOOSE_CMD) status

migrate-create:
	@read -p "Enter migration name: " name; \
	$(GOOSE_CMD) create $$name sql

psql-dev:
	psql $(DATABASE_URL)

psql-test:
	psql $(TEST_DATABASE_URL)

db-dump:
	@docker compose exec db pg_dump -U dev --schema-only -d devdb > internal/db/schema.sql

install-tools:
	@go install github.com/pressly/goose/v3/cmd/goose@latest
	@go install github.com/air-verse/air@latest
	@go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

sqlc:
	@sqlc generate

test-api:
	@APP_ENV=test go test -v ./internal/handlers
