all: tidy build test

tidy:
	go mod tidy -v

build:
	go build .
	go build -o returnstyles ./cmd/returnstyles

test:
	go test -v .

clean:
	rm -f returnstyles
	go clean -cache -modcache
