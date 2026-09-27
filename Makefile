
deps:
	go mod download
	go mod tidy

build: deps
	mkdir -p exec
	go build -o exec/server ./cmd/server
	go build -o exec/gui ./cmd/client/gui
	go build -o exec/tui ./cmd/client/tui

server:
	./exec/server

tui:
	./exec/tui

gui:
	./exec/gui

test:
	go test ./tests/network/... ./tests/scenarios/... ./tests/testTUI/... -count=1

format:
	go fmt ./...

clean:
	rm -rf exec

re:
	make clean
	make build

lint:
	go fmt ./...
	go vet ./...

.PHONY: deps build server tui gui test format clean re lint
