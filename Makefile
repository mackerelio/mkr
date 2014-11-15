BIN = gomkr

all: clean build test

test: deps
	go test ./...

build: deps
	go build -o $(BIN) .

deps:
	go get -d .

clean:
	rm -f build/$(BIN)
	go clean

.PHONY: test build deps clean
