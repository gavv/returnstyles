all: tidy build test

tidy:
	go mod tidy -v

build:
	go build .
	go build -o returnstyles ./cmd/returnstyles

test: build
	go test -v .
	./returnstyles -debug fpsv .

clean:
	rm -f returnstyles
	go clean -cache -modcache

md:
	md-authors -a -f modern AUTHORS.md
