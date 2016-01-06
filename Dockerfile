FROM centos:latest

# install go
RUN yum install -y wget rpmdevtools make git
RUN wget https://storage.googleapis.com/golang/go1.5.2.linux-amd64.tar.gz -P /tmp
RUN tar zxvf /tmp/go1.5.2.linux-amd64.tar.gz -C /usr/local 
ENV GOPATH /go
ENV PATH $PATH:/usr/local/go/bin:$GOPATH/bin

# rpm config
ADD ./ /go/src/github.com/pyama86/STNS
RUN chown root:root -R /go/src/github.com/pyama86/STNS/RPM
RUN echo '%_topdir /go/src/github.com/pyama86/STNS/RPM' > ~/.rpmmacros

# rpm build
WORKDIR /go/src/github.com/pyama86/STNS
RUN go get github.com/tools/godep
RUN godep restore
RUN go build -o RPM/BUILD/stns
CMD rpmbuild -ba RPM/SPECS/stns.spec
