eval $(docker-machine env dev)
rm -rf releases/stns*rpm
docker build --rm -f docker/Rhel -t centos:stns . && docker run -v "$(pwd)"/releases:/go/src/github.com/pyama86/STNS/releases centos:stns
