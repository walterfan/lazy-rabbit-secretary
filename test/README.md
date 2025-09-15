# API Testing

This directory contains pytest-based integration tests for the Lazy Rabbit Reminder API.

## Test Suites

- **`test_secret_api.py`** - Tests for Secret Management API endpoints
- **`test_task_api.py`** - Tests for Task Management API endpoints

## Setup

1. Install Poetry (if not already installed):
```bash
curl -sSL https://install.python-poetry.org | python3 -
# or
pip install poetry
```

2. Install test dependencies:
```bash
cd test
poetry install

# Or use the Makefile
make install

# Or install manually with pip
pip install pytest httpx python-dotenv
```

3. Start the server:
```bash
cd ../backend
go run main.go server
```

4. Set up test environment variables:
```bash
# Generate a test KEK (32 bytes, base64 encoded)
export KEK_BASE64=$(openssl rand -base64 32)
export KEK_VERSION=1

# Or use the Makefile to generate and save to .env
make setup-env
source .env

# Or set a raw 32-byte key
export KEK="your_32_byte_key_exactly_32_chars"
export KEK_VERSION=1
```

5. Create a test user with admin role (if not exists):
```bash
# This depends on your user creation process
# You may need to create via API or database directly
```

## Running Tests

### Using Poetry directly:
```bash
# Run all tests
poetry run pytest -v

# Run specific test files
poetry run pytest test_secret_api.py -v
poetry run pytest test_task_api.py -v

# Run specific test categories
poetry run pytest test_secret_api.py -k "auth" -v
poetry run pytest test_task_api.py -k "status" -v

# Run specific tests
poetry run pytest test_secret_api.py::TestSecretAPI::test_create_secret_success -v
poetry run pytest test_task_api.py::TestTaskAPI::test_task_status_transitions -v

# Run with coverage
poetry run pytest --cov=. --cov-report=html -v
```

### Using Makefile (recommended):
```bash
make test              # Run all tests
make test-verbose      # Run with verbose output
make test-coverage     # Run with coverage report
make test-auth         # Run only authentication tests
make test-debug        # Run in debug mode
make format            # Format code
make lint              # Lint code
make all               # Format, lint, and test
```

### Using the shell script:
```bash
./run_tests.sh         # Normal run
./run_tests.sh verbose # Detailed output
./run_tests.sh coverage # With coverage
./run_tests.sh format  # Format and lint code
```

### Using Poetry shell:
```bash
# Activate virtual environment
poetry shell

# Then run tests normally
pytest test_secret_api.py -v
```

## Test Structure

### `test_secret_api.py`
Main test file containing:
- **Authentication Tests**: Login/logout functionality
- **Secret Creation Tests**: Creating secrets with encryption
- **Secret Retrieval Tests**: Getting secrets by ID
- **Search Tests**: Searching and filtering secrets
- **Error Handling Tests**: Invalid inputs, unauthorized access

### `conftest.py`
Pytest configuration with:
- Test environment setup
- Shared fixtures
- Test markers (integration, auth)

### Test Data
Tests use randomly generated data to avoid conflicts:
- Secret names include random suffixes
- Unique groups and paths per test
- Base64-encoded test values

## Test User Requirements

The tests expect a user with:
- **Username**: `admin` (configurable in test)
- **Password**: `admin123` (configurable in test)  
- **Realm**: `default`
- **Role**: `admin` or `super_admin` (required for secret API access)

## Environment Variables

The test suite uses `python-dotenv` to automatically load environment variables from `.env` files in this order of precedence:
1. `./test/.env` (highest priority)
2. `./.env` (project root)  
3. `./backend/.env` (lowest priority)

Environment variables that are already set in your shell will not be overridden by .env files.

### Required Variables:
- `KEK_BASE64`: Base64-encoded 32-byte encryption key
- `KEK_VERSION`: Key version (default: 1)

### Optional Variables:
- `TEST_BASE_URL`: Server URL (default: http://localhost:9090)
- `TEST_USERNAME`: Test user (default: admin)
- `TEST_PASSWORD`: Test password (default: admin123)
- `TEST_REALM`: Test realm (default: default)
- `TEST_VERIFY_SSL`: Enable SSL verification (default: false)
- `DB_TYPE`: Database type for backend (default: sqlite)
- `DB_FILE_PATH`: Database file path (default: test_lazy_rabbit_reminder.db)

### Example .env file:
```bash
# Copy env.example to .env and customize
cp env.example .env

# Or create manually:
# Encryption
KEK_BASE64=your_base64_encoded_32_byte_key_here
KEK_VERSION=1

# Test configuration
TEST_BASE_URL=https://localhost:9090
TEST_USERNAME=admin
TEST_PASSWORD=admin123
TEST_REALM=default
TEST_VERIFY_SSL=false

# Database (for backend)
DB_TYPE=sqlite
DB_FILE_PATH=test_lazy_rabbit_reminder.db
```

## Test Coverage

The tests cover:
- ✅ User authentication (login/token handling)
- ✅ Secret creation with AES-GCM encryption
- ✅ Secret retrieval and search
- ✅ Input validation and error handling
- ✅ Authorization (admin role required)
- ✅ Base64 encoding verification
- ✅ Realm resolution

## SSL Configuration

The test suite supports both HTTP and HTTPS endpoints with configurable SSL verification:

### For Self-Signed Certificates:
```bash
# In your .env file
TEST_BASE_URL=https://localhost:9090
TEST_VERIFY_SSL=false
```

### For Valid SSL Certificates:
```bash
# In your .env file
TEST_BASE_URL=https://localhost:9090
TEST_VERIFY_SSL=true
```

**Note**: SSL verification is disabled by default (`TEST_VERIFY_SSL=false`) to support self-signed certificates commonly used in development.

## Troubleshooting

### Common Issues:

1. **SSL Certificate Verification Failed**:
   - Set `TEST_VERIFY_SSL=false` for self-signed certificates
   - Ensure `TEST_BASE_URL` uses `https://` for TLS connections

2. **Authentication Failed**:
   - Ensure test user exists with admin role
   - Check username/password in test

3. **KEK Not Configured**:
   - Set `KEK_BASE64` environment variable
   - Ensure it's a valid 32-byte base64 string

4. **Server Not Running**:
   - Start the Go server: `go run main.go server`
   - Check server is listening on expected port

5. **Database Issues**:
   - Ensure database is initialized
   - Check realm "default" exists

### Debug Mode:
```bash
# Run with detailed output and no capture
poetry run pytest test_secret_api.py -v -s --tb=long
```

### Code Quality:
```bash
# Format code with black
poetry run black .

# Lint with flake8
poetry run flake8 .

# Run tests with coverage
poetry run pytest --cov=. --cov-report=term-missing
```
