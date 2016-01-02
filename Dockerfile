FROM centos:latest

# install go
RUN yum install -y wget rpmdevtools make git
RUN rm -rf /usr/local/go
RUN wget https://storage.googleapis.com/golang/go1.5.2.linux-amd64.tar.gz -P /tmp
RUN tar zxvf /tmp/go1.5.2.linux-amd64.tar.gz
RUN mv /go /usr/local
ENV GOPATH /go
ENV PATH /usr/local/go/bin:$PATH

# rpm config
RUN mkdir -p /rpmbuild
ADD ./ /rpmbuild
RUN chown root:root -R /rpmbuild/RPM
RUN echo '%_topdir /rpmbuild/RPM' > ~/.rpmmacros

# rpm build
WORKDIR /rpmbuild
RUN go get github.com/tools/godep
CMD godep restore
CMD go build -o /rpmbuild/RPM/BUILD/tns
CMD rpmbuild -ba RPM/SPECS/tns.spec
