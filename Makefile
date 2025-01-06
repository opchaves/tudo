include .env

GOOSE_CMD=goose -dir ./migrations postgres $(DATABASE_URL)

.PHONY: migrate-up migrate-down migrate-status migrate-create

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

install-tools:
	go install github.com/pressly/goose/v3/cmd/goose@latest
	go install github.com/air-verse/air@latest
