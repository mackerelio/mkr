BIN = mkr

all: clean cross test

test: testdeps
	go test -v ./...

build: deps
	go build -o $(BIN) .

cross: deps
	mkdir -p build
	gox -osarch="linux/amd64" -output build/linux/amd64/mkr
	gox -osarch="darwin/amd64" -output build/darwin/amd64/mkr

deps:
	go get -d -v .

testdeps:
	go get -d -v -t .

clean:
	rm -fr build
	go clean

.PHONY: test build cross deps testdeps clean
