# Define variables
PROJECT_NAME := lazy-rabbit-secretary
BINARY_NAME := $(PROJECT_NAME)
GO := go
GOFLAGS := -ldflags="-s -w"
GO_FILES := $(shell find . -name '*.go')
GOPKG := github.com/walterfan/$(PROJECT_NAME)
ENV_FILE := .env

# Cross-compilation targets
OS := linux darwin windows
ARCH := amd64 arm64
BUILD_DIR := build

# Default target
all: build-all

# Build for current platform
build:
	@echo "Building the project for current platform..."
	@$(GO) build $(GOFLAGS) -o $(BINARY_NAME) main.go
	@echo "Build complete."

# Build for all platforms
build-all:
	@echo "Building for all platforms..."
	@mkdir -p $(BUILD_DIR)
	@echo "Building for linux/amd64..."
	@GOOS=linux GOARCH=amd64 CGO_ENABLED=0 $(GO) build $(GOFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 main.go
	@echo "Building for linux/arm64..."
	@GOOS=linux GOARCH=arm64 CGO_ENABLED=0 $(GO) build $(GOFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 main.go
	@echo "Building for darwin/amd64..."
	@GOOS=darwin GOARCH=amd64 CGO_ENABLED=1 $(GO) build $(GOFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 main.go
	@echo "Building for darwin/arm64..."
	@GOOS=darwin GOARCH=arm64 CGO_ENABLED=1 $(GO) build $(GOFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 main.go
	@echo "Building for windows/amd64..."
	@GOOS=windows GOARCH=amd64 CGO_ENABLED=0 $(GO) build $(GOFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe main.go
	@echo "Building for windows/arm64..."
	@GOOS=windows GOARCH=arm64 CGO_ENABLED=0 $(GO) build $(GOFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-arm64.exe main.go
	@echo "All builds complete."

# Build for specific platform
build-linux:
	@echo "Building for Linux..."
	@mkdir -p $(BUILD_DIR)
	@echo "Building for linux/amd64..."
	@GOOS=linux GOARCH=amd64 CGO_ENABLED=0 $(GO) build $(GOFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 main.go
	@echo "Building for linux/arm64..."
	@GOOS=linux GOARCH=arm64 CGO_ENABLED=0 $(GO) build $(GOFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 main.go
	@echo "Linux builds complete."

build-darwin:
	@echo "Building for macOS..."
	@mkdir -p $(BUILD_DIR)
	@echo "Building for darwin/amd64..."
	@GOOS=darwin GOARCH=amd64 CGO_ENABLED=1 $(GO) build $(GOFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 main.go
	@echo "Building for darwin/arm64..."
	@GOOS=darwin GOARCH=arm64 CGO_ENABLED=1 $(GO) build $(GOFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 main.go
	@echo "macOS builds complete."

build-windows:
	@echo "Building for Windows..."
	@mkdir -p $(BUILD_DIR)
	@echo "Building for windows/amd64..."
	@GOOS=windows GOARCH=amd64 CGO_ENABLED=0 $(GO) build $(GOFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe main.go
	@echo "Building for windows/arm64..."
	@GOOS=windows GOARCH=arm64 CGO_ENABLED=0 $(GO) build $(GOFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-arm64.exe main.go
	@echo "Windows builds complete."

# Build static binaries (for Docker)
build-static:
	@echo "Building static binaries..."
	@mkdir -p $(BUILD_DIR)
	@echo "Building static binary for linux/amd64..."
	@GOOS=linux GOARCH=amd64 CGO_ENABLED=0 $(GO) build -a -ldflags '-s -w' -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64-static main.go
	@echo "Building static binary for linux/arm64..."
	@GOOS=linux GOARCH=arm64 CGO_ENABLED=0 $(GO) build -a -ldflags '-s -w' -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64-static main.go
	@echo "Static builds complete."

# Run the project
run:
	@echo "Running the project..."
	@$(GO) run $(GOFLAGS) main.go

execute:
	@echo "Execute the project..."
	./lazy-rabbit-secretary serve --config ./config/config.yml
# Test the project
test:
	@echo "Running tests..."
	@$(GO) test -v ./...

# Clean up build artifacts
clean:
	@echo "Cleaning up..."
	@rm -f $(BINARY_NAME)
	@rm -rf $(BUILD_DIR)
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

# List all built binaries
list-binaries:
	@echo "Built binaries:"
	@ls -la $(BUILD_DIR)/ 2>/dev/null || echo "No binaries found. Run 'make build-all' first."

# Help
help:
	@echo "Available commands:"
	@echo "  make build         - Build for current platform"
	@echo "  make build-all     - Build for all platforms (linux, darwin, windows)"
	@echo "  make build-linux   - Build for Linux (amd64, arm64)"
	@echo "  make build-darwin  - Build for macOS (amd64, arm64)"
	@echo "  make build-windows - Build for Windows (amd64, arm64)"
	@echo "  make build-static  - Build static binaries for Docker"
	@echo "  make run           - Run the project"
	@echo "  make execute       - Execute with config"
	@echo "  make test          - Run tests"
	@echo "  make clean         - Clean up build artifacts"
	@echo "  make deps          - Install dependencies"
	@echo "  make load-env      - Load environment variables"
	@echo "  make list-binaries - List all built binaries"
	@echo "  make help          - Show this help message"