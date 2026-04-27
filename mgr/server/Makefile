BINARY_NAME := clearbill$(if $(filter Windows_NT,$(OS)),.exe,)
DIST_DIR := dist
CMD_DIR := ./cmd/clearbill

.PHONY: all build clean dist run start test mod lint wire swagger

all: swagger wire test build

build:
	go build -o $(DIST_DIR)/$(BINARY_NAME) $(CMD_DIR)

dist: build

run:
	go run $(CMD_DIR)

start: run

test:
	go test ./...

mod:
	go mod tidy

lint:
	go vet ./...

wire:
	go run github.com/google/wire/cmd/wire@v0.7.0 gen ./internal/app

swagger:
	go run github.com/swaggo/swag/cmd/swag@v1.16.4 init --parseInternal --parseDependency --parseDepth 2 -g cmd/clearbill/main.go -o api/docs

clean:
	rm -rf $(DIST_DIR)
	mkdir -p $(DIST_DIR)
