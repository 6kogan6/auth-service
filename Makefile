-include .env

MIGRATIONS_DIR=./migrations

.PHONY:	db-up db-down db-reset migrate-up migrate-down migrate-status migrate-reset run fmt tidy test

install-tools:
	go install github.com/pressly/goose/v3/cmd/goose@latest

db-up:
	docker compose up -d

db-down:
	docker compose down

db-reset:
	docker compose down -v
	docker compose up -d

migrate-up:
	goose -dir $(MIGRATIONS_DIR) postgres "$(DB_URL)" up

migrate-down:
	goose -dir $(MIGRATIONS_DIR) postgres "$(DB_URL)" down

migrate-status:
	goose -dir $(MIGRATIONS_DIR) postgres "$(DB_URL)" status

migrate-reset:
	goose -dir $(MIGRATIONS_DIR) postgres "$(DB_URL)" reset

run:
	go run ./cmd

fmt:
	go fmt ./...

tidy:
	go mod tidy

test:
	go test ./...
