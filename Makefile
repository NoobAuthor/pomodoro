# Pomodoro CLI Makefile

.PHONY: build clean test lint run install help

# Variables
BINARY_NAME=pomodoro-cli
BUILD_DIR=bin
MAIN_PATH=.
VERSION?=dev
LDFLAGS=-ldflags="-s -w -X main.version=${VERSION}"

# Default target
all: clean lint test build

# Build the application
build:
	@echo "Building ${BINARY_NAME}..."
	@mkdir -p ${BUILD_DIR}
	@go build ${LDFLAGS} -o ${BUILD_DIR}/${BINARY_NAME} ${MAIN_PATH}

# Build for multiple platforms
build-all: clean
	@echo "Building for all platforms..."
	@mkdir -p ${BUILD_DIR}
	@GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o ${BUILD_DIR}/${BINARY_NAME}-linux-amd64 ${MAIN_PATH}
	@GOOS=linux GOARCH=arm64 go build ${LDFLAGS} -o ${BUILD_DIR}/${BINARY_NAME}-linux-arm64 ${MAIN_PATH}
	@GOOS=darwin GOARCH=amd64 go build ${LDFLAGS} -o ${BUILD_DIR}/${BINARY_NAME}-darwin-amd64 ${MAIN_PATH}
	@GOOS=darwin GOARCH=arm64 go build ${LDFLAGS} -o ${BUILD_DIR}/${BINARY_NAME}-darwin-arm64 ${MAIN_PATH}
	@GOOS=windows GOARCH=amd64 go build ${LDFLAGS} -o ${BUILD_DIR}/${BINARY_NAME}-windows-amd64.exe ${MAIN_PATH}

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf ${BUILD_DIR}
	@go clean

# Run tests
test:
	@echo "Running tests..."
	@go test -v -race -coverprofile=coverage.out ./...

# Run tests with coverage report
test-coverage: test
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run linter
lint:
	@echo "Running linter..."
	@golangci-lint run

# Run the application in development mode
run:
	@go run ${MAIN_PATH}

# Install dependencies
deps:
	@echo "Installing dependencies..."
	@go mod download
	@go mod verify

# Install the binary to GOPATH/bin
install: build
	@echo "Installing ${BINARY_NAME}..."
	@cp ${BUILD_DIR}/${BINARY_NAME} ${GOPATH}/bin/

# Generate module dependencies
mod-tidy:
	@go mod tidy

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# Vet code
vet:
	@echo "Vetting code..."
	@go vet ./...

# Quick development workflow
dev: fmt vet test run

# Release workflow
release: clean lint test build-all

# Help
help:
	@echo "Available targets:"
	@echo "  build       - Build the application"
	@echo "  build-all   - Build for all platforms"
	@echo "  clean       - Clean build artifacts"
	@echo "  test        - Run tests"
	@echo "  test-coverage - Run tests with coverage report"
	@echo "  lint        - Run linter"
	@echo "  run         - Run the application"
	@echo "  deps        - Install dependencies"
	@echo "  install     - Install binary to GOPATH/bin"
	@echo "  mod-tidy    - Tidy module dependencies"
	@echo "  fmt         - Format code"
	@echo "  vet         - Vet code"
	@echo "  dev         - Quick development workflow"
	@echo "  release     - Release workflow"
	@echo "  help        - Show this help message"
