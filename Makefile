include .env

GOOSE_CMD=goose -dir ./migrations postgres $(DATABASE_URL)

.PHONY: migrate-up
migrate-up:
	@$(GOOSE_CMD) up

.PHONY: migrate-down
migrate-down:
	@$(GOOSE_CMD) down

.PHONY: migrate-status
migrate-status:
	@$(GOOSE_CMD) status

.PHONY: migrate-create
migrate-create:
	@read -p "Enter migration name: " name; \
	$(GOOSE_CMD) create $$name sql

.PHONY: psql-dev
psql-dev:
	psql $(DATABASE_URL)

.PHONY: psql-test
psql-test:
	psql $(TEST_DATABASE_URL)

.PHONY: db-dump
db-dump:
	@docker compose exec db pg_dump -U dev --schema-only -d devdb > internal/db/schema.sql

.PHONY: db-restore
install-tools:
	@go install github.com/pressly/goose/v3/cmd/goose@latest
	@go install github.com/air-verse/air@latest
	@go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

.PHONY: sqlc
sqlc:
	@sqlc generate

.PHONY: test-api
test-api:
	@APP_ENV=test go test -v ./internal/handlers

.PHONY: start
start:
	@go build -o tmp/tudo .
	@APP_ENV=production ./tmp/tudo

.PHONY: dc-up
dc-up:
	@docker compose up -d --build

.PHONY: dc-stop
dc-stop:
	@docker compose stop

.PHONY: dc-down
dc-down:
	@docker compose down --remove-orphans --volumes
