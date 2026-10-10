
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
TILESETS_DIR := $(GUI_DIR)/tilesets
ASSETS_DIR := $(GUI_DIR)/assets

# Source files
COMMON_FILES := $(shell find $(ENGINE_DIR) $(PROTOCOL_DIR) -type f -name '*.go' 2>/dev/null)
SERVER_FILES := $(shell find $(SERVER_DIR) $(SERVER_CMD_DIR) -type f -name '*.go' 2>/dev/null) $(COMMON_FILES)
CLIENT_FILES := $(shell find $(NETWORK_DIR) $(CONTROLLER_DIR) $(STATE_DIR) -type f -name '*.go' 2>/dev/null)
GUI_FILES := $(shell find $(GUI_DIR) $(GUI_CMD_DIR) -type f -name '*.go' 2>/dev/null) $(COMMON_FILES) $(CLIENT_FILES)
TUI_FILES := $(shell find $(TUI_DIR) $(TUI_CMD_DIR) -type f -name '*.go' 2>/dev/null) $(COMMON_FILES) $(CLIENT_FILES)

# Build
install:
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy

$(EXEC_DIR)/server: $(SERVER_FILES) go.mod go.sum
	@mkdir -p $(EXEC_DIR)
	@echo "Building server..."
	@go build -o $@ $(SERVER_CMD_DIR)

$(EXEC_DIR)/gui: $(GUI_FILES) go.mod go.sum
	@mkdir -p $(EXEC_DIR)
	@echo "Building GUI client..."
	@go build -o $@ $(GUI_CMD_DIR)

$(EXEC_DIR)/tui: $(TUI_FILES) go.mod go.sum
	@mkdir -p $(EXEC_DIR)
	@echo "Building TUI client..."
	@go build -o $@ $(TUI_CMD_DIR)

build: install $(EXEC_DIR)/server $(EXEC_DIR)/gui $(EXEC_DIR)/tui
	@echo "Build completed."

# Execution
run-server: $(EXEC_DIR)/server
	@echo "Starting server..."
	@./$(EXEC_DIR)/server

run-client: $(EXEC_DIR)/tui
	@echo "Starting TUI client..."
	@./$(EXEC_DIR)/tui

run-client-gui: $(EXEC_DIR)/gui
	@echo "Starting GUI client..."
	@./$(EXEC_DIR)/gui

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
	@rm -rf $(EXEC_DIR)


.PHONY: install build run-server run-client run-client-gui test format clean re lint check
