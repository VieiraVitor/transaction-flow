DB_URL="postgres://postgres:postgres@localhost:5433/transactions?sslmode=disable"
MIGRATIONS_PATH=./migrations

# Commands
.PHONY: all build migrate up down run swag lint test

## Run all the commands
all: format lint test run

## Build docker
build:
	docker build --target build -t transaction-flow .

## Run Migrations
migrate-up:
	migrate -database $(DB_URL) -path $(MIGRATIONS_PATH) up

## Down Migrations
migrate-down:
	migrate -database $(DB_URL) -path $(MIGRATIONS_PATH) down

## Start projetct locally without docker compose
run:
	go run ./cmd/transaction-flow/main.go

## Update swagger
swag:
	swag init -g cmd/transaction-flow/main.go --output ./docs

## Run Linter
lint:
	golangci-lint run

## Run tests
test:
	go test -v ./...
