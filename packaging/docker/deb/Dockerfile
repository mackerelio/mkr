FROM debian:wheezy
ENV DEBIAN_FRONTEND noninteractive

RUN apt-get -y update
RUN apt-get -y install build-essential devscripts debhelper fakeroot --no-install-recommends

WORKDIR /deb/build
ENTRYPOINT ["debuild"]

