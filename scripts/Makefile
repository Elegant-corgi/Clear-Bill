.PHONY: all clean dist

all: dist

dist:
	mkdir -p dist
	cp -f *.sh dist/ 2>/dev/null || true
	cp -f README.md dist/README.md

clean:
	rm -rf dist
	mkdir -p dist
