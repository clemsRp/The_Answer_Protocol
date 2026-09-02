
server:
	go run ./cmd/server

client:
	go run ./cmd/client/tui

gui:
	go run ./cmd/client/gui

debug_project:
	(tree -I 'node_modules|venv|.git|__pycache__' && echo -e "\n=== CONTENU DES FICHIERS ===\n" && find . -type f ! -path '*/.*' ! -path '*/node_modules/*' ! -path '*/venv/*' ! -name 'Makefile' ! -name '*.ans' ! -name '*.excalidraw' ! -name '*.html' ! -name '*.png' ! -name '*.jpg' ! -name '*.jpeg' ! -name '*.gif' ! -name '*.svg' ! -name '*.webp' -exec sh -c 'for f; do echo "\n--- FILE: $$f ---"; cat "$$f"; done' _ {} +) > project.txt


.PHONY: client server debug_project

