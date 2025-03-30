# Makefile
BIN_NAME = clothing-recommendation

.PHONY: help build run migrate docker-up docker-down clean

help:  ## Show this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n\nTargets:\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

build: ## Build the application
	@echo "Building..."
	@CGO_ENABLED=0 go build -o bin/$(BIN_NAME) ./cmd/clothing-recommendation

run: build ## Run the application
	@echo "Starting server..."
	@CONFIG_PATH=./config/local.yaml ./bin/$(BIN_NAME)

migrate: ## Run database migrations
	@echo "Running migrations..."
	@go run migrations/migration.go

docker-up: ## Start all services
	@docker compose up --build -d

docker-down: ## Stop all services
	@docker compose down

clean: ## Clean build artifacts
	@rm -rf bin/
	@docker compose rm -f