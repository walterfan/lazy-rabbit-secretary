#!/bin/bash
set -e

# Test runner script for Secret API
echo "🧪 Starting Secret API Tests"

# Check if server is running
echo "📡 Checking if server is running..."
if ! curl -s http://localhost:9090/ping > /dev/null; then
    echo "❌ Server is not running on localhost:9090"
    echo "Please start the server first:"
    echo "  cd ../backend && go run main.go server"
    exit 1
fi
echo "✅ Server is running"

# Check if Poetry is installed
if ! command -v poetry &> /dev/null; then
    echo "❌ Poetry is not installed. Please install it first:"
    echo "  curl -sSL https://install.python-poetry.org | python3 -"
    echo "  # or"
    echo "  pip install poetry"
    echo ""
    echo "Alternatively, use pip to install dependencies:"
    echo "  pip install pytest requests pytest-cov pytest-mock black flake8"
    exit 1
fi

# Set up test KEK if not present
if [ -z "$KEK_BASE64" ] && [ -z "$KEK" ]; then
    echo "🔐 Generating test KEK..."
    export KEK_BASE64=$(openssl rand -base64 32)
    export KEK_VERSION=1
    echo "Generated KEK_BASE64: $KEK_BASE64"
fi

# Install dependencies if needed
if [ ! -f "poetry.lock" ] || [ ! -d ".venv" ]; then
    echo "📦 Installing test dependencies with Poetry..."
    poetry install
fi

# Run tests
echo "🚀 Running tests..."
echo "=================================="

# Run with different verbosity levels based on argument
case "${1:-normal}" in
    "verbose"|"-v")
        poetry run pytest test_secret_api.py -v -s
        ;;
    "quiet"|"-q")
        poetry run pytest test_secret_api.py -q
        ;;
    "debug")
        poetry run pytest test_secret_api.py -v -s --tb=long --capture=no
        ;;
    "auth")
        poetry run pytest test_secret_api.py -k "auth" -v
        ;;
    "integration")
        poetry run pytest test_secret_api.py -m integration -v
        ;;
    "coverage")
        poetry run pytest test_secret_api.py --cov=. --cov-report=html --cov-report=term-missing -v
        ;;
    "format")
        echo "🎨 Formatting code..."
        poetry run black .
        echo "🔍 Linting code..."
        poetry run flake8 .
        ;;
    *)
        poetry run pytest test_secret_api.py
        ;;
esac

echo "=================================="
echo "✅ Tests completed!"
