FROM debian:jessie
ENV DEBIAN_FRONTEND noninteractive

RUN apt-get -y update && apt-get -y install git gnupg

RUN mkdir -p /go/src/github.com/mackerelio/mkr

RUN mkdir -p /rpm /deb/build
RUN mkdir -p /rpm/BUILD /rpm/RPMS /rpm/SOURCES /rpm/SPECS /rpm/SRPMS

VOLUME ["/go/src/github.com/mackerelio/mkr", "/rpm", "/deb"]
