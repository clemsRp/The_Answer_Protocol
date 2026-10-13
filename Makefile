
MAKEFLAGS += --no-print-directory

# Directories
ROOT_DIR := .
SRC_DIR := $(ROOT_DIR)/src
EXEC_DIR := $(ROOT_DIR)/exec

# Main packages
# Server/Common
SERVER_DIR := $(SRC_DIR)/server
ENGINE_DIR := $(SRC_DIR)/engine
PROTOCOL_DIR := $(SRC_DIR)/protocol

# Client
CLIENT_DIR := $(SRC_DIR)/client
GUI_DIR := $(CLIENT_DIR)/gui
TUI_DIR := $(CLIENT_DIR)/tui

NETWORK_DIR := $(CLIENT_DIR)/network
CONTROLLER_DIR := $(CLIENT_DIR)/controller
STATE_DIR := $(CLIENT_DIR)/state

CMD_DIR := $(SRC_DIR)/cmd
TESTS_DIR := $(SRC_DIR)/tests

# Entrypoints
CMD_DIR := $(SRC_DIR)/cmd
GUI_CMD_DIR := $(CMD_DIR)/client/gui
TUI_CMD_DIR := $(CMD_DIR)/client/tui
SERVER_CMD_DIR := $(CMD_DIR)/server

# Utils
MAPS_DIR := $(GUI_DIR)/maps
GUI_ASSETS_DIR := $(GUI_DIR)/assets
TUI_ASSETS_DIR := $(TUI_DIR)/assets
SHADER_DIR := $(GUI_DIR)/src/drawer/shader
MAP_FILE := $(SRC_DIR)/world.json

GUI_UTILS := $(MAPS_DIR) $(GUI_ASSETS_DIR) $(SHADER_DIR)

# Source files
COMMON_FILES := $(shell find $(ENGINE_DIR) $(PROTOCOL_DIR) -type f -name '*.go' 2>/dev/null) $(MAP_FILE) go.mod go.sum
SERVER_FILES := $(shell find $(SERVER_DIR) $(SERVER_CMD_DIR) -type f -name '*.go' 2>/dev/null) $(COMMON_FILES)
CLIENT_FILES := $(shell find $(NETWORK_DIR) $(CONTROLLER_DIR) $(STATE_DIR) -type f -name '*.go' 2>/dev/null)
GUI_FILES := $(shell find $(GUI_DIR) $(GUI_CMD_DIR) -type f -name '*.go' 2>/dev/null) $(COMMON_FILES) $(CLIENT_FILES) $(GUI_UTILS)
TUI_FILES := $(shell find $(TUI_DIR) $(TUI_CMD_DIR) -type f -name '*.go' 2>/dev/null) $(COMMON_FILES) $(CLIENT_FILES) $(TUI_ASSETS_DIR)

# Build
install:
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy

$(EXEC_DIR)/server: $(SERVER_FILES)
	@mkdir -p $(EXEC_DIR)
	@echo "Building server..."
	@go build -o $@ $(SERVER_CMD_DIR)

$(EXEC_DIR)/gui: $(GUI_FILES)
	@mkdir -p $(EXEC_DIR)
	@echo "Building GUI client..."
	@go build -o $@ $(GUI_CMD_DIR)

$(EXEC_DIR)/tui: $(TUI_FILES)
	@mkdir -p $(EXEC_DIR)
	@echo "Building TUI client..."
	@go build -o $@ $(TUI_CMD_DIR)

build: install $(EXEC_DIR)/server $(EXEC_DIR)/gui $(EXEC_DIR)/tui
	@echo "Build completed."

# Execution
ARGS ?=

run-server: $(EXEC_DIR)/server
	@echo "Starting server..."
	@./$(EXEC_DIR)/server $(ARGS)

run-client: $(EXEC_DIR)/tui
	@echo "Starting TUI client..."
	@./$(EXEC_DIR)/tui $(ARGS)

run-client-gui: $(EXEC_DIR)/gui
	@echo "Starting GUI client..."
	@./$(EXEC_DIR)/gui $(ARGS)

# Tests
test:
	@echo "Running tests..."
	@go test $(TESTS_DIR)/network/... \
			 $(TESTS_DIR)/scenarios/... \
			 $(TESTS_DIR)/leaks/... \
			 -count=1

test-race:
	@echo "Running tests..."
	@go test -race $(TESTS_DIR)/network/... \
			 $(TESTS_DIR)/scenarios/... \
			 $(TESTS_DIR)/leaks/... \
			 -count=1

test-strict:
	@echo "Test eveything..."
	@$(MAKE) test
	@$(MAKE) test-race
	@echo "All tests passed successfully...."

# Manage project
re:
	@$(MAKE) clean
	@$(MAKE) build

format:
	@echo "Format code..."
	@go fmt ./...
	@echo "Linting code..."
	@go vet ./...

lint:
	@echo "Static check code..."
	@PATH="$$PATH:$$(go env GOPATH)/bin" which staticcheck >/dev/null 2>&1 && PATH="$$PATH:$$(go env GOPATH)/bin" staticcheck ./... || true
	@echo "Check lint..."
	@test -z "$(gofmt -l src)" || (echo "Error: Incorrect lint !" && exit 1)

check:
	@echo "Running full project check..."
	@$(MAKE) format
	@$(MAKE) lint
	@$(MAKE) re
	@$(MAKE) test-strict
	@$(MAKE) clean
	@echo "All checks passed successfully."

stats:
	@echo "Go files: $$(find . -name "*.go" -type f | wc -l)"
	@echo "Total lines: $$(find . -name "*.go" -type f -exec cat {} + | wc -l)"

# Clean
clean:
	@echo "Cleaning executables..."
	@rm -rf $(EXEC_DIR)


.PHONY: install build run-server run-client run-client-gui test test-race stats clean re format lint check
