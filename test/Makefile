.PHONY: help install test test-verbose test-coverage format lint clean all setup-env show-env

# Default target
help:
	@echo "Available targets:"
	@echo "  install       - Install dependencies with Poetry"
	@echo "  test          - Run all tests"
	@echo "  test-verbose  - Run tests with verbose output"
	@echo "  test-coverage - Run tests with coverage report"
	@echo "  test-auth     - Run only authentication tests"
	@echo "  test-debug    - Run tests in debug mode"
	@echo "  format        - Format code with Black"
	@echo "  lint          - Lint code with Flake8"
	@echo "  clean         - Clean up generated files"
	@echo "  setup-env     - Create .env file with default settings"
	@echo "  show-env      - Show current environment variables"
	@echo "  all           - Install, format, lint, and test"

# Install dependencies
install:
	@echo "📦 Installing dependencies with Poetry..."
	poetry install

# Run tests
test:
	@echo "🧪 Running tests..."
	poetry run pytest test_secret_api.py

# Run tests with verbose output
test-verbose:
	@echo "🧪 Running tests (verbose)..."
	poetry run pytest test_secret_api.py -v -s

# Run tests with coverage
test-coverage:
	@echo "🧪 Running tests with coverage..."
	poetry run pytest test_secret_api.py --cov=. --cov-report=html --cov-report=term-missing -v

# Run only authentication tests
test-auth:
	@echo "🔐 Running authentication tests..."
	poetry run pytest test_secret_api.py -k "auth" -v

# Run tests in debug mode
test-debug:
	@echo "🐛 Running tests in debug mode..."
	poetry run pytest test_secret_api.py -v -s --tb=long --capture=no

# Format code
format:
	@echo "🎨 Formatting code with Black..."
	poetry run black .

# Lint code
lint:
	@echo "🔍 Linting code with Flake8..."
	poetry run flake8 .

# Clean up
clean:
	@echo "🧹 Cleaning up..."
	rm -rf __pycache__/
	rm -rf .pytest_cache/
	rm -rf htmlcov/
	rm -rf .coverage
	rm -f test_lazy_rabbit_reminder.db

# All-in-one target
all: install format lint test

# Setup environment variables
setup-env:
	@echo "🔐 Setting up test environment..."
	@if [ ! -f ".env" ]; then \
		if [ -f "env.example" ]; then \
			cp env.example .env; \
			sed -i.bak "s/your_base64_encoded_32_byte_key_here/$$(openssl rand -base64 32)/" .env; \
			rm -f .env.bak; \
			echo "✅ Created .env from env.example with generated KEK"; \
		else \
			echo "# Test environment variables" > .env; \
			echo "KEK_BASE64=$$(openssl rand -base64 32)" >> .env; \
			echo "KEK_VERSION=1" >> .env; \
			echo "" >> .env; \
			echo "# Test configuration" >> .env; \
			echo "TEST_BASE_URL=http://localhost:9090" >> .env; \
			echo "TEST_USERNAME=admin" >> .env; \
			echo "TEST_PASSWORD=admin123" >> .env; \
			echo "TEST_REALM=default" >> .env; \
			echo "" >> .env; \
			echo "# Database configuration" >> .env; \
			echo "DB_TYPE=sqlite" >> .env; \
			echo "DB_FILE_PATH=test_lazy_rabbit_reminder.db" >> .env; \
			echo "✅ Created .env file with test configuration"; \
		fi; \
		echo "📝 Edit .env to customize settings"; \
	else \
		echo "⚠️  .env file already exists"; \
	fi

# Show current environment
show-env:
	@echo "🔍 Current environment variables:"
	@echo "  KEK_BASE64=$${KEK_BASE64:-<not set>}" | sed 's/\(KEK_BASE64=.\{8\}\).*/\1.../'
	@echo "  KEK_VERSION=$${KEK_VERSION:-<not set>}"
	@echo "  TEST_BASE_URL=$${TEST_BASE_URL:-<not set>}"
	@echo "  TEST_USERNAME=$${TEST_USERNAME:-<not set>}"
	@echo "  TEST_REALM=$${TEST_REALM:-<not set>}"
