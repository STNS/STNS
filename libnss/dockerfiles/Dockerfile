FROM centos:latest
MAINTAINER pyama86 <www.kazu.com@gmail.com>

RUN yum install -y glibc gcc make libcurl-devel wget bzip2 git vim epel-release && \
    yum install -y clang
RUN mkdir /stns
WORKDIR /stns
