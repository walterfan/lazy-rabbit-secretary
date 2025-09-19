# Define variables
PROJECT_NAME := lazy-rabbit-secretary
BINARY_NAME := $(PROJECT_NAME)
GO := go
GOFLAGS := -ldflags="-s -w"
GO_FILES := $(shell find . -name '*.go')
GOPKG := github.com/walterfan/$(PROJECT_NAME)
ENV_FILE := .env

# Default target
all: build

# Build the project
build:
	@echo "Building the project..."
	@$(GO) build $(GOFLAGS) -o $(BINARY_NAME) main.go
	@echo "Build complete."

# Run the project
run:
	@echo "Running the project..."
	@$(GO) run $(GOFLAGS) main.go

execute:
	@echo "Execute the project..."
	./lazy-rabbit-secretary serve --config ./config/walter-config.yml
# Test the project
test:
	@echo "Running tests..."
	@$(GO) test -v ./...

# Clean up build artifacts
clean:
	@echo "Cleaning up..."
	@rm -f $(BINARY_NAME)
	@echo "Cleanup complete."

# Install dependencies
deps:
	@echo "Installing dependencies..."
	@$(GO) mod download
	@echo "Dependencies installed."

# Load environment variables
load-env:
	@echo "Loading environment variables..."
	@if [ -f $(ENV_FILE) ]; then \
		export $$(cat $(ENV_FILE) | xargs); \
		echo "Environment variables loaded."; \
	else \
		echo "No .env file found."; \
	fi

# Help
help:
	@echo "Available commands:"
	@echo "  make build    - Build the project"
	@echo "  make run      - Run the project"
	@echo "  make test     - Run tests"
	@echo "  make clean    - Clean up build artifacts"
	@echo "  make deps     - Install dependencies"
	@echo "  make load-env - Load environment variables"
	@echo "  make help     - Show this help message"