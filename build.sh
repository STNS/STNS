#!/bin/bash
eval $(docker-machine env dev)
rm -rf ./binary/*

docker build --no-cache --rm -t stns:stns .
docker run -v "$(pwd)"/binary:/go/src/github.com/STNS/STNS/binary -t stns:stns

while getopts brd OPT
do
    case $OPT in
        r)
          docker build --no-cache --rm -f docker/rpm -t stns:stns .
          docker run -v "$(pwd)"/binary:/go/src/github.com/STNS/STNS/binary -t stns:stns
            ;;
        d)
          docker build --no-cache --rm -f docker/deb -t stns:stns .
          docker run -v "$(pwd)"/binary:/go/src/github.com/STNS/STNS/binary -t stns:stns
            ;;
    esac
done
shift $((OPTIND - 1))
