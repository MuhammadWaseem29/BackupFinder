# BackupFinder Makefile
# Automated build and setup commands

# Variables
BINARY_NAME := backupfinder
CMD_DIR := ./cmd/backupfinder
BUILD_FLAGS := -ldflags="-w -s"
GO_VERSION := 1.18

# Default target
.PHONY: all
all: check build test

# Check if Go is installed and version is sufficient
.PHONY: check
check:
	@echo "🔍 Checking Go installation..."
	@which go > /dev/null || (echo "❌ Go is not installed. Run 'make install-go' or './install.sh'" && exit 1)
	@echo "✅ Go is installed: $$(go version)"

# Install Go (Linux only)
.PHONY: install-go
install-go:
	@echo "🔧 Installing Go..."
	@chmod +x install.sh
	@./install.sh

# Build the application
.PHONY: build
build: check
	@echo "🏗️  Building BackupFinder..."
	@go build $(BUILD_FLAGS) -o $(BINARY_NAME) $(CMD_DIR)
	@chmod +x $(BINARY_NAME)
	@echo "✅ Build complete: $(BINARY_NAME)"

# Build for multiple platforms
.PHONY: build-all
build-all: check
	@echo "🏗️  Building for multiple platforms..."
	@GOOS=linux GOARCH=amd64 go build $(BUILD_FLAGS) -o $(BINARY_NAME)-linux-amd64 $(CMD_DIR)
	@GOOS=linux GOARCH=arm64 go build $(BUILD_FLAGS) -o $(BINARY_NAME)-linux-arm64 $(CMD_DIR)
	@GOOS=darwin GOARCH=amd64 go build $(BUILD_FLAGS) -o $(BINARY_NAME)-darwin-amd64 $(CMD_DIR)
	@GOOS=darwin GOARCH=arm64 go build $(BUILD_FLAGS) -o $(BINARY_NAME)-darwin-arm64 $(CMD_DIR)
	@GOOS=windows GOARCH=amd64 go build $(BUILD_FLAGS) -o $(BINARY_NAME)-windows-amd64.exe $(CMD_DIR)
	@echo "✅ Multi-platform build complete"

# Run tests
.PHONY: test
test: build
	@echo "🧪 Testing BackupFinder..."
	@./$(BINARY_NAME) health-check
	@./$(BINARY_NAME) version
	@echo "✅ Tests passed"

# Run system diagnostics
.PHONY: system-check
system-check: build
	@echo "🔍 Running system diagnostics..."
	@./$(BINARY_NAME) system-check

# Setup (interactive)
.PHONY: setup
setup: build
	@echo "🔧 Running interactive setup..."
	@./$(BINARY_NAME) setup

# Quick start guide
.PHONY: quick-start
quick-start: build
	@echo "🚀 Showing quick start guide..."
	@./$(BINARY_NAME) quick-start

# Clean build artifacts
.PHONY: clean
clean:
	@echo "🧹 Cleaning build artifacts..."
	@rm -f $(BINARY_NAME) $(BINARY_NAME)-*
	@echo "✅ Clean complete"

# Development setup
.PHONY: dev-setup
dev-setup: install-go build test
	@echo "🎉 Development environment ready!"

# Run with example (requires target)
.PHONY: example
example: build
	@echo "📖 Running example scan..."
	@./$(BINARY_NAME) -u httpbin.org -w --silent | head -10
	@echo "..."
	@echo "ℹ️  Use 'make run TARGET=your-domain.com' for custom target"

# Run with custom target
.PHONY: run
run: build
ifndef TARGET
	@echo "❌ Please specify TARGET: make run TARGET=example.com"
	@exit 1
endif
	@echo "🎯 Scanning $(TARGET)..."
	@./$(BINARY_NAME) -u $(TARGET) -w

# Install as system command (Linux/macOS)
.PHONY: install
install: build
	@echo "📦 Installing BackupFinder system-wide..."
	@sudo cp $(BINARY_NAME) /usr/local/bin/
	@echo "✅ Installed to /usr/local/bin/$(BINARY_NAME)"
	@echo "   You can now run 'backupfinder' from anywhere"

# Uninstall system command
.PHONY: uninstall
uninstall:
	@echo "🗑️  Uninstalling BackupFinder..."
	@sudo rm -f /usr/local/bin/$(BINARY_NAME)
	@echo "✅ Uninstalled"

# Show help
.PHONY: help
help:
	@echo "BackupFinder Makefile Commands:"
	@echo ""
	@echo "Setup Commands:"
	@echo "  make install-go     Install Go (Linux only)"
	@echo "  make dev-setup      Complete development setup"
	@echo "  make setup          Interactive setup guide"
	@echo ""
	@echo "Build Commands:"
	@echo "  make build          Build BackupFinder"
	@echo "  make build-all      Build for multiple platforms"
	@echo "  make clean          Clean build artifacts"
	@echo ""
	@echo "Test Commands:"
	@echo "  make test           Run health checks"
	@echo "  make system-check   Run system diagnostics"
	@echo "  make example        Run example scan"
	@echo ""
	@echo "Usage Commands:"
	@echo "  make run TARGET=domain.com    Scan specific target"
	@echo "  make quick-start             Show usage guide"
	@echo ""
	@echo "Install Commands:"
	@echo "  make install        Install system-wide"
	@echo "  make uninstall      Remove system installation"
	@echo ""
	@echo "Utility Commands:"
	@echo "  make check          Check Go installation"
	@echo "  make help           Show this help"

# Default help when no target specified
.DEFAULT_GOAL := help
