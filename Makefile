BIN = mkr

VERSION = $$(git describe --tags --always --dirty) ($(git name-rev --name-only HEAD | sed 's/^remotes\/origin\///'))

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
	goxc -tasks='xc archive' -bc 'linux,!arm windows darwin' -d . -build-ldflags "-X main.Version \"$(VERSION)\"" -resources-include='README*'

deps:
	go get -d -v .

testdeps:
	go get -d -v -t .

clean:
	rm -fr build
	go clean

.PHONY: test build cross lint deps testdeps clean
