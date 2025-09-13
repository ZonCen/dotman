.PHONY: test test-verbose test-coverage test-unit test-integration build clean

# Default target
all: test build

# Run all tests
test:
	go test ./...

# Run tests with verbose output
test-verbose:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Run only unit tests
test-unit:
	go test ./internal/...

# Run only integration tests
test-integration:
	go test ./cmd/...

# Build the application
build:
	go build -o dotman

# Clean build artifacts
clean:
	rm -f dotman coverage.out coverage.html

# Install dependencies
deps:
	go mod tidy
	go mod download

# Run linter
lint:
	golangci-lint run --config .golangci.yml

# Format code
fmt:
	go fmt ./...

# Run tests in watch mode (requires entr)
test-watch:
	find . -name "*.go" | entr -c make test

# Help
help:
	@echo "Available targets:"
	@echo "  test          - Run all tests"
	@echo "  test-verbose  - Run tests with verbose output"
	@echo "  test-coverage - Run tests with coverage report"
	@echo "  test-unit     - Run only unit tests"
	@echo "  test-integration - Run only integration tests"
	@echo "  build         - Build the application"
	@echo "  clean         - Clean build artifacts"
	@echo "  deps          - Install dependencies"
	@echo "  lint          - Run linter"
	@echo "  fmt           - Format code"
	@echo "  test-watch    - Run tests in watch mode (requires entr)"
	@echo "  help          - Show this help"
