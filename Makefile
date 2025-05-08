BIN := mkr
# This VERSION variable indicates the latest tag.
VERSION := $(subst v,,$(shell git describe --abbrev=0 --tags))
BUILD_LDFLAGS := "-w -s"
export CGO_ENABLED := 0

.PHONY: all
all: clean cross test rpm deb

.PHONY: test
test:
	go test -v ./...

.PHONY: build
build:
	go build -ldflags=$(BUILD_LDFLAGS) -o $(BIN) .

.PHONY: cross
cross:
	go tool github.com/Songmu/goxz/cmd/goxz -d snapshot -os darwin -arch amd64,arm64 \
	  -build-ldflags=$(BUILD_LDFLAGS)
	go tool github.com/Songmu/goxz/cmd/goxz -d snapshot -os linux -arch 386,amd64,arm64,arm \
	  -build-ldflags=$(BUILD_LDFLAGS)

.PHONY: rpm
rpm: rpm-v2

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

.PHONY: deb
deb: deb-v2-x86 deb-v2-arm64 deb-v2-arm

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

.PHONY: clean
clean:
	rm -fr build snapshot
	go clean
