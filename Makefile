BIN := mkr
VERSION := 0.40.4
CURRENT_REVISION := $(shell git rev-parse --short HEAD)
BUILD_LDFLAGS := "-w -s -X main.gitcommit=$(CURRENT_REVISION)"

export GO111MODULE=on

.PHONY: all
all: clean cross lint gofmt test rpm deb

.PHONY: test-deps
test-deps:
	cd && \
	go get golang.org/x/lint/golint

.PHONY: devel-deps
devel-deps: test-deps
	cd && \
	go get github.com/mattn/goveralls && \
	go get github.com/Songmu/goxz/cmd/goxz

.PHONY: test
test: test-deps
	go test -v ./...

.PHONY: build
build:
	go build -ldflags=$(BUILD_LDFLAGS) -o $(BIN) .

.PHONY: lint
lint: test-deps
	go vet ./...
	golint -set_exit_status ./...

GOFMT_RET = .gofmt.txt
.PHONY: gofmt
gofmt: test-deps
	rm -f $(GOFMT_RET)
	gofmt -s -d *.go | tee $(GOFMT_RET)
	test ! -s $(GOFMT_RET)

.PHONY: cross
cross: devel-deps
	goxz -d snapshot -os darwin -arch amd64 \
	  -build-ldflags=$(BUILD_LDFLAGS)
	goxz -d snapshot -os linux -arch 386,amd64,arm64,arm \
	  -build-ldflags=$(BUILD_LDFLAGS)

.PHONY: rpm
rpm: rpm-v1 rpm-v2

.PHONY: rpm-v1
rpm-v1:
	GOOS=linux GOARCH=386 make build
	rpmbuild --define "_builddir `pwd`" --define "_version ${VERSION}" --define "buildarch noarch" --target noarch -bb packaging/rpm/mkr.spec
	GOOS=linux GOARCH=amd64 make build
	rpmbuild --define "_builddir `pwd`" --define "_version ${VERSION}" --define "buildarch x86_64" --target x86_64  -bb packaging/rpm/mkr.spec

.PHONY: rpm-v2
rpm-v2: rpm-v2-x86 rpm-v2-arm64

.PHONY: rpm-v2-x86
rpm-v2-x86:
	GOOS=linux GOARCH=amd64 make build
	rpmbuild --define "_builddir `pwd`" --define "_version ${VERSION}" \
	  --define "buildarch x86_64" --target x86_64 --define "dist .el7.centos" \
	  -bb packaging/rpm/mkr-v2.spec
	rpmbuild --define "_builddir `pwd`" --define "_version ${VERSION}" \
	  --define "buildarch x86_64" --target x86_64 --define "dist .amzn2" \
	  -bb packaging/rpm/mkr-v2.spec

.PHONY: rpm-v2-arm64
rpm-v2-arm64:
	GOOS=linux GOARCH=arm64 make build
	rpmbuild --define "_builddir `pwd`" --define "_version ${VERSION}" \
	  --define "buildarch aarch64" --target aarch64 --define "dist .el7.centos" \
	  -bb packaging/rpm/mkr-v2.spec
	rpmbuild --define "_builddir `pwd`" --define "_version ${VERSION}" \
	  --define "buildarch aarch64" --target aarch64 --define "dist .amzn2" \
	  -bb packaging/rpm/mkr-v2.spec

PHONY: deb
deb: deb-v1 deb-v2

.PHONY: deb-v1
deb-v1:
	GOOS=linux GOARCH=386 make build
	cp $(BIN) packaging/deb/debian/$(BIN).bin
	cd packaging/deb && debuild --no-tgz-check -rfakeroot -uc -us

.PHONY: deb-v2
deb-v2: deb-v2-x86 deb-v2-arm64 deb-v2-arm

.PHONY: deb-v2-x86
deb-v2-x86:
	GOOS=linux GOARCH=amd64 make build
	cp $(BIN) packaging/deb-v2/debian/$(BIN).bin
	cd packaging/deb-v2 && debuild --no-tgz-check -rfakeroot -uc -us

.PHONY: deb-v2-arm64
deb-v2-arm64:
	GOOS=linux GOARCH=arm64 make build
	cp $(BIN) packaging/deb-v2/debian/$(BIN).bin
	cd packaging/deb-v2 && debuild --no-tgz-check -rfakeroot -uc -us -aarm64

.PHONY: deb-v2-arm
deb-v2-arm:
	GOOS=linux GOARCH=arm ARM=6 make build # Build ARMv6 binary for Raspbian
	cp $(BIN) packaging/deb-v2/debian/$(BIN).bin
	cd packaging/deb-v2 && debuild --no-tgz-check -rfakeroot -uc -us -aarmhf

.PHONY: check-release-deps
check-release-deps:
	@have_error=0; \
	for command in cpanm hub ghch gobump; do \
	  if ! command -v $$command > /dev/null; then \
	    have_error=1; \
	    echo "\`$$command\` command is required for releasing"; \
	  fi; \
	done; \
	test $$have_error = 0

.PHONY: release
release: check-release-deps
	(cd script && cpanm -qn --installdeps .)
	perl script/create-release-pullrequest

.PHONY: clean
clean:
	rm -fr build snapshot
	go clean

.PHONY: cover
cover: devel-deps
	goveralls
