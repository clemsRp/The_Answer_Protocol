
MAKEFLAGS += --no-print-directory

# Cleaning
MAPS_DIR = ./client/gui/maps
TILESETS_DIR = ./client/gui/tilesets
ASSETS_DIR = ./client/gui/assets

# Source files
COMMON_FILES := $(shell find engine protocol -type f -name '*.go' 2>/dev/null)
SERVER_FILES := $(shell find client/controller client/network client/state -type f -name '*.go' 2>/dev/null)
CLIENT_FILES := $(shell find engine protocol -type f -name '*.go' 2>/dev/null)
GUI_FILES    := $(shell find client/gui cmd/client/gui -type f -name '*.go' 2>/dev/null) $(COMMON_FILES) $(CLIENT_FILES)
TUI_FILES    := $(shell find client/tui cmd/client/tui -type f -name '*.go' 2>/dev/null) $(COMMON_FILES) $(CLIENT_FILES)

# Build
install:
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy

build: install exec/server exec/gui exec/tui
	@echo "Build completed."

# Execution
run-server:
	@echo "Starting server..."
	go run ./cmd/server/main.go

run-client:
	@echo "Starting TUI client..."
	go run ./cmd/client/tui/main.go

run-client-gui:
	@echo "Starting GUI client..."
	go run ./cmd/client/gui/main.go

# Tests
test:
	@echo "Running tests..."
	@go test ./tests/network/... ./tests/scenarios/... ./tests/leaks/... -count=1

# Manage project
format:
	@echo "Formatting code..."
	@go fmt ./...

re:
	@$(MAKE) clean
	@$(MAKE) build

lint:
	@echo "Linting code..."
	@go fmt ./...
	@go vet ./...

check:
	@echo "Running full project check..."
	@$(MAKE) re
	@$(MAKE) lint
	@$(MAKE) test
	@$(MAKE) clean
	@echo "All checks passed successfully."

# Clean
clean:
	@echo "Cleaning executables..."
	@rm -rf exec


.PHONY: install build  run-server run-client run-client-gui test format clean_tsx clean_img clean clean_strict re lint check
