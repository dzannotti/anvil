# Detect OS for cross platform compatibility
ifeq ($(OS),Windows_NT)
	DETECTED_OS := Windows
	# Set Windows specific variables here
	RM := del /Q
	MKDIR := mkdir
	BINARY_NAME := anvil.exe
	PATH_SEPARATOR := \\
else
	DETECTED_OS := $(shell uname -s)
	RM := rm -f
	MKDIR := mkdir -p
	BINARY_NAME := anvil
	PATH_SEPARATOR := /
endif

# Go related variables
GOBASE := $(shell go list -m)
GOBIN := ./bin
GOSRC := ./cmd/cli/
GOGUISRC := ./cmd/gui/

# Simplified run that doesn't rebuild unless needed
.PHONY: run-fast
run-fast:
	@cd $(GOSRC) && go run .

# Main targets
.PHONY: all build clean test lint run help deps

all: clean lint test build

build:
	@echo Building for $(DETECTED_OS)...
	@go build -o $(GOBIN)/$(BINARY_NAME) $(GOSRC)/main.go

clean:
	@echo Cleaning...
ifeq ($(DETECTED_OS),Windows)
	@if exist $(GOBIN) rmdir /s /q $(GOBIN)
	@mkdir $(GOBIN)
else
	@rm -rf $(GOBIN)
	@mkdir -p $(GOBIN)
endif

release:
	@echo Building for $(DETECTED_OS)...
	@go build -trimpath -ldflags="-w -s" -o $(GOBIN)/$(BINARY_NAME) $(GOSRC)/main.go

test:
	@echo Running tests...
	@go test -v ./...

lint:
	@echo Running linter...
	@golangci-lint run --fix ./...

fmt:
	@echo Running formatter...
	@go fmt ./...
	@goimports -w .

run:
	@echo Running application...
	@go run $(GOSRC)/main.go

cli:
	@make run

gui:
	@echo Running gui application...
	@go run $(GOGUISRC)

deps:
	@echo Installing dependencies...
	@go mod tidy
	@go mod download
	@go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.1.6
	@go install github.com/mickamy/gotcha@latest

loc:
	@echo Counting LOC...
	@npx sloc --keys source --format cli-table --format-option no-head --exclude ".*_test.go" .

tdd:
	@echo Running tests...
	@gotcha watch

help:
	@echo Make targets:"
	@echo   all    - Clean, lint, test, and build
	@echo   build  - Build the application
	@echo   clean  - Clean build artifacts
	@echo   test   - Run tests
	@echo   tdd   - Run tests tdd
	@echo   lint   - Run linter
	@echo   run    - Run the application
	@echo   gui    - Run the gui application
	@echo   run-fast - Run without full make overhead
	@echo   deps   - Install dependencies
	@echo   loc    - Count lines of code
	@echo   fmt    - Format code
	@echo   help   - Show this help message
