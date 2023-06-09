all: tidy build test

tidy:
	go mod tidy -v

build:
	go build .
	go build -o returnstyles ./cmd/returnstyles

test:
ifneq ($(shell which gotest),)
	gotest -v .
else
	go test -v .
endif
