# 1. Load environment variables from .env file for all commands
ifneq (,$(wildcard ./.env))
    include .env
    export
endif

# Default target when you just run 'make'
.DEFAULT_GOAL := help

.PHONY: run sqlc build tidy help docker-up docker-down

## docker-up: Start the local PostgreSQL container in the background
docker-up:
	@echo "🐘 Starting PostgreSQL container..."
	@docker compose up -d
	@echo "⏳ Waiting for database to be ready..."
	@until docker exec go_dev_db pg_isready -U $(DB_USER) -d $(DB_NAME) > /dev/null 2>&1; do \
		sleep 1; \
	done
	@echo "✅ Database is fully ready!"

## docker-down: Stop and remove the PostgreSQL container
docker-down:
	@echo "🛑 Stopping PostgreSQL container..."
	@docker compose down

## run: Load env vars and start the Go application
## run: Ensure docker container is running, generate code, and start Go application
run: sqlc tidy
	@echo "🚀 Starting Go application..."
	@go run .
	#docker-up <- add to target later

## sqlc: Run sqlc code generation
sqlc:
	@echo "⚡ Generating sqlc code..."
	@../tools/sqlc generate

## build: Compile the binary
build: sqlc tidy
	@echo "📦 Building application binary..."
	@go build -o bin/app .

## tidy: Run go mod tidy to clean up dependencies
tidy:
	@echo "🧹 Cleaning up Go modules..."
	@go mod tidy

## help: Show this help message with descriptions
help:
	@echo "Available commands:"
	@grep -F -h "##" $(MAKEFILE_LIST) | grep -F -v grep | sed -e 's/\\$$//' | sed -e 's/##//'