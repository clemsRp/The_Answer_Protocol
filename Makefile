
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
deps:
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy

exec/server: $(SERVER_FILES)
	@mkdir -p exec
	@echo "Building server..."
	@go build -o exec/server ./cmd/server

exec/gui: $(GUI_FILES)
	@mkdir -p exec
	@echo "Building GUI client..."
	@go build -o exec/gui ./cmd/client/gui

exec/tui: $(TUI_FILES)
	@mkdir -p exec
	@echo "Building TUI client..."
	@go build -o exec/tui ./cmd/client/tui

build: deps exec/server exec/gui exec/tui
	@echo "Build completed."

# Execution
server: exec/server
	@echo "Starting server..."
	./exec/server

tui: exec/tui
	@echo "Starting TUI client..."
	./exec/tui

gui: exec/gui
	@echo "Starting GUI client..."
	./exec/gui

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
clean_tsx:
	@echo "Cleaning tilesets..."
	@./clean_tsx.sh $(MAPS_DIR) $(TILESETS_DIR) ./

clean_img:
	@echo "Cleaning images..."
	@./clean_img.sh $(MAPS_DIR) $(ASSETS_DIR)

clean:
	@echo "Cleaning executables..."
	@rm -rf exec

clean_strict:
	@$(MAKE) clean
	@$(MAKE) clean_tsx
	@$(MAKE) clean_img

.PHONY: deps build server tui gui test format clean_tsx clean_img clean clean_strict re lint check
