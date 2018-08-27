FROM debian:stretch
MAINTAINER pyama86 <www.kazu.com@gmail.com>

RUN apt-get -qq update && \
    apt-get install -qq gcc make debhelper dh-make clang git curl devscripts
ARG GO_VERSION

ENV FILE go$GO_VERSION.linux-amd64.tar.gz
ENV URL https://storage.googleapis.com/golang/$FILE

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

RUN set -eux &&\
  curl -OL $URL &&\
	tar -C /usr/local -xzf $FILE &&\
	rm $FILE &&\
  mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"

ENV USER root
ADD . /go/src/github.com/STNS/STNS
WORKDIR /go/src/github.com/STNS/STNS
