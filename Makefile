.PHONY: test test-verbose test-cover test-integration test-all test-bench clean db-setup db-test help

# Default environment variables
DB_HOST ?= localhost
DB_PORT ?= 5432
DB_USER ?= benedictmwanga
DB_PASSWORD ?= 
DB_NAME ?= validation_test

# Export environment variables for tests
export DB_HOST
export DB_PORT
export DB_USER
export DB_PASSWORD
export DB_NAME

# Default target
all: test

# Run all tests
test:
	go test -v ./tests

# Run specific tests
test-unique:
	go test -v ./tests -run TestUnique

test-exists:
	go test -v ./tests -run TestExists

# Run all tests with verbose output
test-verbose:
	go test -v ./tests

# Run tests with coverage
test-cover:
	go test -v ./tests -coverprofile=coverage.out
	go tool cover -html=coverage.out

# Run integration tests (all database-related tests)
test-integration: test-unique test-exists

# Run benchmark tests
test-bench:
	go test -v ./tests -bench=. -run=Benchmark

# Run specific benchmark
test-bench-unique:
	go test -v ./tests -bench=BenchmarkUniqueValidation

test-bench-exists:
	go test -v ./tests -bench=BenchmarkExistsValidation

# Clean test artifacts
clean:
	rm -f coverage.out
	go clean -testcache

# Setup database
db-setup:
	mysql -h $(DB_HOST) -u $(DB_USER) -p$(DB_PASSWORD) < scripts/setup_db.sql

# Run tests with database setup
test-all: db-setup test-verbose

# Run tests without cache (fresh test)
test-fresh:
	go clean -testcache
	go test -v ./tests

# Show help
help:
	@echo "Available targets:"
	@echo ""
	@echo "  Testing:"
	@echo "  make test              - Run all tests"
	@echo "  make test-unique       - Run only unique validation tests"
	@echo "  make test-exists       - Run only exists validation tests"
	@echo "  make test-verbose      - Run all tests with verbose output"
	@echo "  make test-fresh        - Clear cache and run all tests"
	@echo "  make test-cover        - Run tests with coverage report"
	@echo "  make test-integration  - Run all integration tests"
	@echo "  make test-bench        - Run all benchmark tests"
	@echo "  make test-bench-unique - Run unique validation benchmark"
	@echo "  make test-bench-exists - Run exists validation benchmark"
	@echo ""
	@echo "  Database:"
	@echo "  make db-setup          - Setup test database"
	@echo "  make test-all          - Setup DB and run all tests"
	@echo ""
	@echo "  Maintenance:"
	@echo "  make clean             - Clean test artifacts"
	@echo "  make help              - Show this help"
	@echo ""
	@echo "Current database settings:"
	@echo "  DB_HOST: $(DB_HOST)"
	@echo "  DB_PORT: $(DB_PORT)"
	@echo "  DB_USER: $(DB_USER)"
	@echo "  DB_NAME: $(DB_NAME)"