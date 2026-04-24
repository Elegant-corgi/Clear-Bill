.PHONY: all build clean dev dist install lint preview

all: build

install:
	pnpm install

dev:
	pnpm dev

lint:
	pnpm lint

build:
	pnpm build

dist: build

preview:
	pnpm preview

clean:
	rm -rf dist
	mkdir -p dist
