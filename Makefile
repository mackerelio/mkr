BIN = mkr
VERSION = 0.16.0
CURRENT_REVISION = $(shell git rev-parse --short HEAD)

all: clean cross lint gofmt test

test: testdeps
	go test -v ./...

build: deps
	go build -ldflags "-X main.gitcommit=$(CURRENT_REVISION)" -o $(BIN) .

lint: testdeps
	go vet ./...
	golint -set_exit_status ./...

GOFMT_RET = .gofmt.txt
gofmt: testdeps
	rm -f $(GOFMT_RET)
	gofmt -s -d *.go | tee $(GOFMT_RET)
	test ! -s $(GOFMT_RET)

cross: deps
	goxc -tasks='xc archive' -bc 'linux,!arm darwin' -d . -build-ldflags "-X main.gitcommit=$(CURRENT_REVISION)" -resources-include='README*'
	cp -p $(PWD)/snapshot/linux_amd64/mkr $(PWD)/snapshot/mkr_linux_amd64
	cp -p $(PWD)/snapshot/linux_386/mkr $(PWD)/snapshot/mkr_linux_386
	cp -p $(PWD)/snapshot/darwin_amd64/mkr $(PWD)/snapshot/mkr_darwin_amd64
	cp -p $(PWD)/snapshot/darwin_386/mkr $(PWD)/snapshot/mkr_darwin_386

rpm:
	GOOS=linux GOARCH=386 make build
	rpmbuild --define "_builddir `pwd`" --define "_version ${VERSION}" --define "buildarch noarch" -bb packaging/rpm/mkr.spec
	GOOS=linux GOARCH=amd64 make build
	rpmbuild --define "_builddir `pwd`" --define "_version ${VERSION}" --define "buildarch x86_64" -bb packaging/rpm/mkr.spec

deb:
	GOOS=linux GOARCH=386 make build
	cp $(BIN) packaging/deb/debian/$(BIN).bin
	cd packaging/deb && debuild --no-tgz-check -rfakeroot -uc -us

deps:
	go get -d -v .

testdeps:
	go get -d -v -t .
	go get github.com/golang/lint/golint
	go get golang.org/x/tools/cmd/cover
	go get github.com/axw/gocov/gocov
	go get github.com/mattn/goveralls

check-release-deps:
	@have_error=0; \
	for command in cpanm hub ghch gobump; do \
	  if ! command -v $$command > /dev/null; then \
	    have_error=1; \
	    echo "\`$$command\` command is required for releasing"; \
	  fi; \
	done; \
	test $$have_error = 0

release: check-release-deps
	(cd script && cpanm -qn --installdeps .)
	perl script/create-release-pullrequest

clean:
	rm -fr build
	go clean

cover: testdeps
	goveralls

.PHONY: test build cross lint gofmt deps testdeps clean deb rpm release cover
