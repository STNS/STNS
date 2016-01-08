eval $(docker-machine env dev)
docker build -f RhelDockerfile -t centos:stns . && docker run -v "$(pwd)"/releases:/go/src/github.com/pyama86/STNS/package/RPM/RPMS centos:stns
