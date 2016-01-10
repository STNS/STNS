eval $(docker-machine env dev)
rm -rf releases/stns*deb
docker build -f docker/Ubuntu -t ubuntu:stns . && docker run -v "$(pwd)"/releases:/go/src/github.com/pyama86/STNS/releases ubuntu:stns
