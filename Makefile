# 🚗 Expense Tracker Bot Makefile

# Variables
BINARY_NAME=expense-tracker-bot
GO=go
GOFMT=gofmt
GOLINT=golangci-lint
GOTEST=go test
GOCOVER=go tool cover

# Emojis for better visibility
EMOJI_BUILD=🔨
EMOJI_RUN=🚀
EMOJI_TEST=🧪
EMOJI_CLEAN=🧹
EMOJI_LINT=🔍
EMOJI_DEPS=📦
EMOJI_FMT=✨
EMOJI_HELP=❓
EMOJI_SETUP=⚙️

# Default target
.PHONY: all
all: fmt lint test build

# Setup targets
.PHONY: setup
setup:
	@echo "$(EMOJI_SETUP) Running complete setup..."
	@chmod +x scripts/*.sh
	@./scripts/setup.sh

.PHONY: setup-quick
setup-quick:
	@echo "$(EMOJI_SETUP) Running quick setup..."
	@chmod +x scripts/*.sh
	@./scripts/quick-start.sh

.PHONY: setup-db
setup-db:
	@echo "$(EMOJI_SETUP) Setting up database..."
	@chmod +x scripts/*.sh
	@./scripts/setup_database.sh

# Build the application
.PHONY: build
build:
	@echo "$(EMOJI_BUILD) Building..."
	$(GO) build -o $(BINARY_NAME) ./cmd

# Run the application
.PHONY: run
run:
	@echo "$(EMOJI_RUN) Running..."
	$(GO) run ./cmd

# Run tests
.PHONY: test
test:
	@echo "$(EMOJI_TEST) Running tests with race detection..."
	$(GOTEST) -v -race ./...

# Run tests silently (for production/CI environments)
.PHONY: test-silent
test-silent:
	@echo "$(EMOJI_TEST) Running tests silently with race detection..."
	$(GOTEST) -race ./...

# Run tests with coverage
.PHONY: test-coverage
test-coverage:
	@echo "$(EMOJI_TEST) Running tests with coverage..."
	$(GOTEST) -coverprofile=coverage.out ./...
	$(GOCOVER) -html=coverage.out

# Clean build files
.PHONY: clean
clean:
	@echo "$(EMOJI_CLEAN) Cleaning..."
	rm -f $(BINARY_NAME)
	rm -f coverage.out

# Run linter
.PHONY: lint
lint:
	@echo "$(EMOJI_LINT) Running linter..."
	$(GOLINT) run

# Format code
.PHONY: fmt
fmt:
	@echo "$(EMOJI_FMT) Formatting code..."
	$(GOFMT) -w .

# Install dependencies
.PHONY: deps
deps:
	@echo "$(EMOJI_DEPS) Installing dependencies..."
	$(GO) mod download
	$(GO) mod tidy

# Show help
.PHONY: help
help:
	@echo "$(EMOJI_HELP) Available commands:"
	@echo ""
	@echo "$(EMOJI_SETUP) Setup commands:"
	@echo "  make setup         - Complete automated setup"
	@echo "  make setup-quick   - Quick setup with confirmation"
	@echo "  make setup-db      - Database setup only"
	@echo ""
	@echo "$(EMOJI_BUILD) Build commands:"
	@echo "  make build         - Build the application"
	@echo "  make run          - Run the application"
	@echo "  make deps         - Install dependencies"
	@echo ""
	@echo "$(EMOJI_TEST) Test commands:"
	@echo "  make test         - Run tests"
	@echo "  make test-silent  - Run tests silently (for production/CI)"
	@echo "  make test-coverage - Run tests with coverage"
	@echo ""
	@echo "$(EMOJI_LINT) Quality commands:"
	@echo "  make lint         - Run linter"
	@echo "  make fmt          - Format code"
	@echo "  make clean        - Clean build files"
	@echo ""
	@echo "$(EMOJI_HELP) Help:"
	@echo "  make help         - Show this help message" 