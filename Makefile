BIN = mkr
VERSION = 0.35.1
CURRENT_REVISION = $(shell git rev-parse --short HEAD)

all: clean cross lint gofmt test rpm deb

deps:
	go get -d -v ./...

test-deps:
	go get -d -v -t ./...
	go get golang.org/x/lint/golint

devel-deps: test-deps
	go get github.com/mattn/goveralls
	go get github.com/Songmu/goxz/cmd/goxz

test: test-deps
	go test -v ./...

build: deps
	go build -ldflags "-w -s -X main.gitcommit=$(CURRENT_REVISION)" -o $(BIN) .

lint: test-deps
	go vet ./...
	golint -set_exit_status ./...

GOFMT_RET = .gofmt.txt
gofmt: test-deps
	rm -f $(GOFMT_RET)
	gofmt -s -d *.go | tee $(GOFMT_RET)
	test ! -s $(GOFMT_RET)

cross: devel-deps
	goxz -d snapshot -os darwin,linux -arch 386,amd64 \
	  -build-ldflags "-X main.gitcommit=$(CURRENT_REVISION)"

rpm: rpm-v1 rpm-v2

rpm-v1:
	GOOS=linux GOARCH=386 make build
	rpmbuild --define "_builddir `pwd`" --define "_version ${VERSION}" --define "buildarch noarch" -bb packaging/rpm/mkr.spec
	GOOS=linux GOARCH=amd64 make build
	rpmbuild --define "_builddir `pwd`" --define "_version ${VERSION}" --define "buildarch x86_64" -bb packaging/rpm/mkr.spec

rpm-v2:
	GOOS=linux GOARCH=amd64 make build
	rpmbuild --define "_builddir `pwd`" --define "_version ${VERSION}" \
	  --define "buildarch x86_64" --define "dist .el7.centos" \
	  -bb packaging/rpm/mkr-v2.spec
	rpmbuild --define "_builddir `pwd`" --define "_version ${VERSION}" \
	  --define "buildarch x86_64" --define "dist .amzn2" \
	  -bb packaging/rpm/mkr-v2.spec

deb: deb-v1 deb-v2

deb-v1:
	GOOS=linux GOARCH=386 make build
	cp $(BIN) packaging/deb/debian/$(BIN).bin
	cd packaging/deb && debuild --no-tgz-check -rfakeroot -uc -us

deb-v2:
	GOOS=linux GOARCH=amd64 make build
	cp $(BIN) packaging/deb-v2/debian/$(BIN).bin
	cd packaging/deb-v2 && debuild --no-tgz-check -rfakeroot -uc -us

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

cover: devel-deps
	goveralls

.PHONY: test build cross lint gofmt deps test-deps devel-deps clean deb deb-v1 deb-v2 rpm rpm-v1 rpm-v2 release cover
