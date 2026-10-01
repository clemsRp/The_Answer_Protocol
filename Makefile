
MAKEFLAGS += --no-print-directory
PROD := false

deps:
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy

build: deps
	@echo "Building binaries..."
	@mkdir -p exec
	@go build -o exec/server ./cmd/server
	@go build -o exec/gui ./cmd/client/gui
	@go build -o exec/tui ./cmd/client/tui
	@echo "Build completed."

server:
	@echo "Starting server..."
ifeq ($(PROD),true)
	./exec/server
else
	go run ./cmd/server
endif

tui:
	@echo "Starting TUI client..."
ifeq ($(PROD),true)
	./exec/tui
else
	go run ./cmd/client/tui
endif

gui:
	@echo "Starting GUI client..."
ifeq ($(PROD),true)
	./exec/client/gui
else
	go run ./cmd/client/gui
endif

test:
	@echo "Running tests..."
	@go test ./tests/network/... ./tests/scenarios/... ./tests/leaks/... -count=1

format:
	@echo "Formatting code..."
	@go fmt ./...

clean:
	@echo "Cleaning executables..."
	@rm -rf exec
	@echo "Executables removed."

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

debug_project:
	@echo "Generating debug file..."
	(tree -I 'node_modules|venv|.git|__pycache__' && echo -e "\n=== FILE CONTENTS ===\n" && find . -type f ! -path '*/.*' ! -path '*/node_modules/*' ! -path '*/venv/*' ! -name 'Makefile' ! -name '*.ans' ! -name '*.excalidraw' ! -name '*.html' ! -name '*.png' ! -name '*.jpg' ! -name '*.jpeg' ! -name '*.gif' ! -name '*.svg' ! -name '*.webp' -exec sh -c 'for f; do echo "\n--- FILE: $$f ---"; cat "$$f"; done' _ {} +) > project.txt
	@echo "project.txt generated."

.PHONY: deps build server tui gui test format clean re lint check debug_project