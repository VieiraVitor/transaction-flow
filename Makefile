DB_URL="postgres://postgres:postgres@localhost:5433/transaction_flow?sslmode=disable"
MIGRATIONS_PATH=./migrations

# Commands
.PHONY: all build migrate up down run swag lint test

## Run all the commands
all: format lint test run

## Build docker
build:
	docker build --target build -t transaction-flow .

## Create Database
create-db:
	@echo "Creating database..."
	@PGPASSWORD=postgres psql -h localhost -U postgres -p 5433 -tc "SELECT 1 FROM pg_database WHERE datname = 'transaction_flow'" | grep -q 1 || psql -h localhost -U postgres -p 5433 -c "CREATE DATABASE transaction_flow"
	@echo "Database created successfully (if it didn't exist)"

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
