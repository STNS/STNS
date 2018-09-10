FROM golang:1.11rc2
RUN apt-get update -qqy --fix-missing
RUN apt-get install -qqy build-essential \
    git \
    curl \
    libcurl4-openssl-dev \
    gdb \
    sudo \
    rsyslog \
    clang \
    lsof
ADD . /go/src/github.com/STNS/STNS
WORKDIR /go/src/github.com/STNS/STNS
