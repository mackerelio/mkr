BIN = gomkr

all: clean cross test

test: deps
	go test ./...

build: deps
	go build -o $(BIN) .

cross: deps
	gox -osarch="linux/amd64" -output build/linux/amd64/gomkr
	gox -osarch="darwin/amd64" -output build/darwin/amd64/gomkr

deps:
	go get -d .

clean:
	rm -f build/$(BIN)
	go clean

.PHONY: test build cross deps clean
