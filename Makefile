BIN = mkr

VERSION = $$(git describe --tags --always --dirty)

BUILD_FLAGS = -ldflags "\
	      -X main.Version=$(VERSION) \
	      "

#all: clean cross lint test
all: clean cross

test: testdeps
	go test -v ./...

build: deps
	go build $(BUILD_FLAGS) -o $(BIN) .

LINT_RET = .golint.txt
lint: testdeps
	go vet
	rm -f $(LINT_RET)
	golint ./... | tee .golint.txt
	test ! -s $(LINT_RET)

cross: deps
	goxc -tasks='xc archive' -bc 'linux,!arm darwin' -d . -build-ldflags "-X main.Version=$(VERSION)" -resources-include='README*'
	cp -p $(PWD)/snapshot/linux_amd64/mkr $(PWD)/snapshot/mkr_linux_amd64
	cp -p $(PWD)/snapshot/linux_386/mkr $(PWD)/snapshot/mkr_linux_386
	cp -p $(PWD)/snapshot/darwin_amd64/mkr $(PWD)/snapshot/mkr_darwin_amd64
	cp -p $(PWD)/snapshot/darwin_386/mkr $(PWD)/snapshot/mkr_darwin_386

rpm:
	GOOS=linux GOARCH=amd64 make build
	rpmbuild --define "_builddir `pwd`" -ba packaging/rpm/mkr.spec

deb:
	GOOS=linux GOARCH=amd64 make build
	cp $(BIN) packaging/deb/debian/$(BIN).bin
	cd packaging/deb && debuild --no-tgz-check -rfakeroot -uc -us

deps:
	go get -d -v .

testdeps:
	go get -d -v -t .
	go get golang.org/x/tools/cmd/vet
	go get github.com/golang/lint/golint
	go get golang.org/x/tools/cmd/cover
	go get github.com/axw/gocov/gocov
	go get github.com/mattn/goveralls

release:
	script/releng

clean:
	rm -fr build
	go clean

cover: testdeps
	goveralls

.PHONY: test build cross lint deps testdeps clean deb rpm release cover
