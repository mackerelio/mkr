#! /bin/sh

set -ex

DOCKER=${DOCKER:-"docker"}

PACKAGE=${PACKAGE:-"mkr"}

IMAGE_SRC="hatena/mkr_src"
IMAGE_GO="hatena/mkr_golang"
IMAGE_DEB="hatena/mkr_deb"
IMAGE_RPM="hatena/mkr_rpm"

NAME="mkr"
REPOSITORY_DIR="/go/src/github.com/mackerelio/mkr"

# UNIXTIMEつけて Data volume コンテナを使いまわさないようにする
UNIXTIME=$(date '+%s')
DATA_VOLUME_CONTAINER=${DATA_VOLUME_CONTAINER:-"$PACKAGE-src-$UNIXTIME"}

function build_all_image() {
  (cd packaging/docker/src    && $DOCKER build -t $IMAGE_SRC .)
  (cd packaging/docker/golang && $DOCKER build -t $IMAGE_GO  .)
  (cd packaging/docker/deb    && $DOCKER build -t $IMAGE_DEB .)
  (cd packaging/docker/rpm    && $DOCKER build -t $IMAGE_RPM .)
}

function build_container() {
  $DOCKER run --name $DATA_VOLUME_CONTAINER -v /Users/Sixeight/local/go/src/github.com/mackerelio/mkr:/go/src/github.com/mackerelio/mkr $IMAGE_SRC
}

function clean_containers() {
  $DOCKER rm $DATA_VOLUME_CONTAINER
}

function src_run() {
  $DOCKER run --rm --volumes-from $DATA_VOLUME_CONTAINER debian:jessie $@
}

function go_run() {
  $DOCKER run --rm --volumes-from $DATA_VOLUME_CONTAINER $IMAGE_GO $@
}

function debuild_run() {
  $DOCKER run --rm --volumes-from $DATA_VOLUME_CONTAINER $IMAGE_DEB $@
}

function rpmbuild_run() {
  $DOCKER run --rm --volumes-from $DATA_VOLUME_CONTAINER $IMAGE_RPM $@
}

function debian_run_info() {
  echo "$DOCKER run -it --rm --volumes-from $DATA_VOLUME_CONTAINER debian:jessie /bin/bash"
}

function centos_run_info() {
  echo "$DOCKER run -it --rm --volumes-from $DATA_VOLUME_CONTAINER centos:centos7 /bin/bash"
}

function prepare() {
  build_all_image
  build_container
}

function build() {
  go_run make clean build
}

function deb() {
  src_run cp -r $REPOSITORY_DIR/packaging/deb/debian/ /deb/build/debian
  src_run cp $REPOSITORY_DIR/$NAME /deb/build/debian/$NAME.bin
  debuild_run --no-tgz-check -uc -us
}

function rpm() {
  src_run cp -r $REPOSITORY_DIR/packaging/rpm/$NAME.spec /rpm/SPECS/$NAME.spec
  src_run cp $REPOSITORY_DIR/$NAME /rpm/BUILD/$NAME
  rpmbuild_run -ba /rpm/SPECS/$NAME.spec
}

function main() {
  prepare
  build
  deb
  rpm
  debian_run_info
  centos_run_info
}

main
