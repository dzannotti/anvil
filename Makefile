# Cross-platform compatibility
ifeq ($(OS),Windows_NT)
	BINARY_EXT := .exe
	MKDIR_CMD := if not exist bin mkdir bin
	RM_CMD := if exist bin rmdir /s /q bin & if exist coverage.out del coverage.out & if exist coverage.html del coverage.html
else
	BINARY_EXT :=
	MKDIR_CMD := mkdir -p bin/
	RM_CMD := rm -rf bin/ coverage.out coverage.html
endif

.PHONY: build test lint clean clean-cache install-tools cli gui fmt tdd help hooks setup run run-cli run-gui

# Build targets
build: cli gui

cli:
	@$(MKDIR_CMD)
	go build -o bin/cli$(BINARY_EXT) ./cmd/cli

gui:
	@$(MKDIR_CMD)
	go build -o bin/gui$(BINARY_EXT) ./cmd/gui

# Run targets
run: run-cli

run-cli:
	go run ./cmd/cli

run-gui:
	go run ./cmd/gui

# Test targets
test:
	gotestsum --format testname ./internal/...

tdd:
	@echo "Starting TDD mode - watching ./internal/ for changes..."
	gotestsum --watch --format dots-v2 --format-icons hivis ./internal/...

test-watch:
	gotestsum --watch --format testname ./internal/...

test-coverage:
	gotestsum --format testname ./internal/... -- -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html

# Format targets
fmt:
	go fmt ./...

fmt-check:
ifeq ($(OS),Windows_NT)
	@for /f %%i in ('gofmt -l .') do @echo Code needs formatting: %%i && exit 1
else
	@test -z "$$(gofmt -l .)" || (echo "Code needs formatting. Run 'make fmt'" && exit 1)
endif

# Lint targets
lint:
	$(HOME)/go/bin/golangci-lint run ./...

lint-fix:
	$(HOME)/go/bin/golangci-lint run --fix ./...

# Development targets
deps:
	go install gotest.tools/gotestsum@latest
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest
	pip install pre-commit

clean:
	@$(RM_CMD)

clean-cache:
	@echo "Cleaning Go cache..."
	go clean -cache
	go clean -testcache
	go clean -modcache

hooks:
	@echo "Installing pre-commit hooks..."
	pre-commit install

setup:
	@echo "Running full development environment setup..."
	./scripts/setup.sh

# CI target
ci: fmt-check test lint build

# Help target
help:
	@echo "Spellbinder - Game Development Commands"
	@echo ""
	@echo "Build Commands:"
	@echo "  build          Build both CLI and GUI binaries"
	@echo "  cli            Build CLI binary only"
	@echo "  gui            Build GUI binary only"
	@echo ""
	@echo "Run Commands:"
	@echo "  run            Run CLI (alias for run-cli)"
	@echo "  run-cli        Run CLI directly"
	@echo "  run-gui        Run GUI directly"
	@echo ""
	@echo "Test Commands:"
	@echo "  test           Run all tests"
	@echo "  tdd            Run tests in watch mode with concise output"
	@echo "  test-coverage  Run tests with coverage report"
	@echo ""
	@echo "Code Quality:"
	@echo "  fmt            Format all Go code"
	@echo "  fmt-check      Check if code is formatted (CI)"
	@echo "  lint           Run linter on all code"
	@echo "  lint-fix       Run linter and auto-fix issues"
	@echo ""
	@echo "Utilities:"
	@echo "  clean          Remove build artifacts"
	@echo "  clean-cache    Clean Go cache and module cache"
	@echo "  deps           Install development tools"
	@echo "  hooks          Install pre-commit hooks"
	@echo "  setup          Full development environment setup"
	@echo "  ci             Run full CI pipeline locally"
	@echo "  help           Show this help message"
