.PHONY: help build run dev clean test setup-simple

# Default target
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Build the application
build: ## Build the application
	@echo "Building application..."
	go build -o ./bin/main.exe .

# Run the application
run: ## Run the application
	@echo "Running application..."
	go run .

# Development with hot reload
dev: ## Run with hot reload using air
	@echo "Starting development server with hot reload..."
	air

# Install air for hot reload
install-air: ## Install air for hot reloading
	@echo "Installing air..."
	@which air > /dev/null 2>&1 || \
		(go install github.com/air-verse/air@latest && echo "Air installed successfully") || \
		(echo "Failed to install air. You may need to install it manually: go install github.com/air-verse/air@latest" && exit 1)

# Clean build artifacts
clean: ## Clean build artifacts
	@echo "Cleaning..."
	rm -rf ./bin ./tmp

# Run tests
test: ## Run tests
	@echo "Running tests..."
	go test ./...

# Tidy dependencies
tidy: ## Tidy go.mod dependencies
	@echo "Tidying dependencies..."
	go mod tidy

# Generate OpenAPI/Swagger documentation
.PHONY: docs
docs:
	swag init -g main.go --parseDependency --parseInternal --outputTypes go,json,yaml --propertyStrategy snakecase
	swagger2openapi ./docs/swagger.json --outfile ./docs/openapi.json --yaml --patch
	npx openapi-typescript ./docs/openapi.json -o ../client/src/api/api.ts

# Setup development environment
setup: install-air tidy ## Setup development environment
	@echo "Development environment setup complete!"
	@echo "Run 'make dev' to start with hot reload"

# Alternative setup (if air install fails)
setup-simple: tidy ## Simple setup without air
	@echo "Simple setup complete!"
	@echo "Run 'make run' for development without hot reload"
	@echo "To install air manually: go install github.com/air-verse/air@latest"

