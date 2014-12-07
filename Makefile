BIN = mkr

VERSION = $$(git describe --tags --always --dirty)

BUILD_FLAGS = -ldflags "\
	      -X main.Version \"$(VERSION)\" \
	      "

all: clean cross test

test: testdeps
	go test -v ./...

build: deps
	go build $(BUILD_FLAGS) -o $(BIN) .

lint: deps testdeps
	go vet
	golint

cross: deps
	mkdir -p build
	gox $(BUILD_FLAGS) -osarch="linux/amd64" -output build/linux/amd64/mkr
	gox $(BUILD_FLAGS) -osarch="darwin/amd64" -output build/darwin/amd64/mkr

deps:
	go get -d -v .

testdeps:
	go get -d -v -t .

clean:
	rm -fr build
	go clean

.PHONY: test build cross lint deps testdeps clean
