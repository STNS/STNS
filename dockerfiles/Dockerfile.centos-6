FROM centos:6
MAINTAINER pyama86 <www.kazu.com@gmail.com>

ARG GO_VERSION

RUN yum install -y epel-release rpmdevtools make glibc gcc
RUN yum install -y http://opensource.wandisco.com/centos/6/git/x86_64/wandisco-git-release-6-1.noarch.rpm
RUN yum install -y git clang

ENV FILE go$GO_VERSION.linux-amd64.tar.gz
ENV URL https://storage.googleapis.com/golang/$FILE

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

RUN set -eux &&\
  yum -y clean all &&\
  curl -OL $URL &&\
	tar -C /usr/local -xzf $FILE &&\
	rm $FILE &&\
  mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"

ADD . /go/src/github.com/STNS/STNS
WORKDIR /go/src/github.com/STNS/STNS

RUN mkdir -p /root/rpmbuild/{BUILD,RPMS,SOURCES,SPECS,SRPMS}
RUN sed -i "s;%_build_name_fmt.*;%_build_name_fmt\t%%{ARCH}/%%{NAME}-%%{VERSION}-%%{RELEASE}.%%{ARCH}.el6.rpm;" /usr/lib/rpm/macros
