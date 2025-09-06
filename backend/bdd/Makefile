# Makefile for BDD tests

.PHONY: test test-verbose test-junit clean install help

# Default target
help:
	@echo "Available targets:"
	@echo "  test         - Run BDD tests with progress format"
	@echo "  test-verbose - Run BDD tests with pretty format"
	@echo "  test-junit   - Run BDD tests and generate JUnit XML report"
	@echo "  clean        - Clean test artifacts"
	@echo "  install      - Install godog CLI tool"
	@echo "  help         - Show this help message"

# Install godog CLI tool
install:
	go install github.com/cucumber/godog/cmd/godog@latest

# Run BDD tests with progress format
test:
	@echo "Running BDD tests..."
	cd .. && go test ./bdd -v

# Run BDD tests with pretty format
test-verbose:
	@echo "Running BDD tests with verbose output..."
	cd .. && godog --format=pretty --paths=features ./bdd

# Run BDD tests with cucumber format (for CI/CD)
test-cucumber:
	@echo "Running BDD tests with cucumber format..."
	cd .. && godog --format=cucumber --paths=features ./bdd

# Run BDD tests and generate JUnit XML report
test-junit:
	@echo "Running BDD tests with JUnit output..."
	mkdir -p reports
	cd .. && godog --format=junit --output=bdd/reports/junit.xml --paths=features ./bdd

# Run specific feature
test-feature:
	@echo "Running specific feature: $(FEATURE)"
	cd .. && godog --format=pretty --paths=features/$(FEATURE) ./bdd

# Run tests with specific tags
test-tags:
	@echo "Running tests with tags: $(TAGS)"
	cd .. && godog --format=pretty --tags="$(TAGS)" --paths=features ./bdd

# Clean test artifacts
clean:
	rm -rf reports/
	rm -f *.log
	rm -f *.xml

# Run tests in Docker (if needed)
test-docker:
	docker run --rm -v $(PWD)/..:/app -w /app golang:1.24-alpine \
		sh -c "apk add --no-cache git && go mod download && go test ./bdd -v"

# Continuous testing (watch for changes)
test-watch:
	@echo "Watching for changes... (requires 'entr' tool)"
	find ../features ../bdd -name "*.feature" -o -name "*.go" | entr -r make test

# Generate step definitions skeleton
generate-steps:
	cd .. && godog --format=pretty --paths=features --dry-run ./bdd

# Validate feature files
validate:
	@echo "Validating feature files..."
	cd .. && godog --format=pretty --paths=features --dry-run ./bdd

# Run all checks (validate + test)
check: validate test

# Debug mode (with more verbose output)
debug:
	cd .. && godog --format=pretty --paths=features --verbose ./bdd
