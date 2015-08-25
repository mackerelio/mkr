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
	goxc -tasks='xc archive' -bc 'linux,!arm darwin' -d . -build-ldflags "-X main.Version \"$(VERSION)\"" -resources-include='README*'
	cp -p $(PWD)/snapshot/linux_amd64/mkr $(PWD)/snapshot/mkr_linux_amd64
	cp -p $(PWD)/snapshot/linux_386/mkr $(PWD)/snapshot/mkr_linux_386
	cp -p $(PWD)/snapshot/darwin_amd64/mkr $(PWD)/snapshot/mkr_darwin_amd64
	cp -p $(PWD)/snapshot/darwin_386/mkr $(PWD)/snapshot/mkr_darwin_386

rpm:
	GOOS=linux GOARCH=386 make build
	rpmbuild --define "_builddir `pwd`" -ba packaging/rpm/mkr.spec

deb:
	GOOS=linux GOARCH=386 make build
	cp $(BIN) packaging/deb/debian/$(BIN).bin
	cd packaging/deb && debuild --no-tgz-check -rfakeroot -uc -us

deps:
	go get -d -v .

testdeps:
	go get -d -v -t .

release:
	script/releng

clean:
	rm -fr build
	go clean

.PHONY: test build cross lint deps testdeps clean deb rpm release
