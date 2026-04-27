BINARY_NAME := clearbill$(if $(filter Windows_NT,$(OS)),.exe,)
DIST_DIR := dist
CMD_DIR := ./cmd/clearbill

.PHONY: all build clean dist run start test

all: test build

build:
	go build -o $(DIST_DIR)/$(BINARY_NAME) $(CMD_DIR)

dist: build

run:
	go run $(CMD_DIR)

start: run

test:
	go test ./...

clean:
	rm -rf $(DIST_DIR)
	mkdir -p $(DIST_DIR)
