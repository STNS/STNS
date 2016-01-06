eval $(docker-machine env dev)
docker build -t centos:stns . && docker run -v "$(pwd)"/releases:/go/src/github.com/pyama86/STNS/RPM/RPMS centos:stns
