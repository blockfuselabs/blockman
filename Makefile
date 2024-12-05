# Makefile for BlockMan

# Variables
GO_CMD=go
MAIN_FILE=cmd/web/main.go

# Default target
run:
	$(GO_CMD) run $(MAIN_FILE)

# Install dependencies
install:
	$(GO_CMD) mod tidy

# Build the application
build:
	$(GO_CMD) build -o blockman $(MAIN_FILE)

# Clean up build artifacts
clean:
	rm -f blockman

# Help
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  run       Run the application"
	@echo "  install   Install dependencies"
	@echo "  build     Build the application"
	@echo "  clean     Clean up build artifacts"
	@echo "  help      Display this help message"
