#!/bin/bash
eval $(docker-machine env dev)
rm -rf ./binary/*
rm -rf ./releases/*

docker build --no-cache --rm -t stns:stns .
docker run -v "$(pwd)"/binary:/go/src/github.com/STNS/STNS/binary -t stns:stns

func_libnss_rpm_build()
{
  cd ../libnss_stns
  ./build.sh -r
  cd ../STNS
}

func_libnss_deb_build()
{
  cd ../libnss_stns
  ./build.sh -d
  cd ../STNS
}

func_package_build()
{
  docker build --no-cache --rm -f docker/$1 -t stns:stns . && \
  docker run -it -v "$(pwd)"/binary:/go/src/github.com/STNS/STNS/binary -t stns:stns
}

func_repo_build()
{
  docker build --no-cache --rm -f docker/$1 -t stns:stns . && \
  docker run -it -v "$(pwd)"/releases:/go/src/github.com/STNS/STNS/releases -t stns:stns
}

while getopts brya OPT
do
    case $OPT in
        r)
          func_package_build "rpm"
            ;;
        d)
          func_package_build "deb"
            ;;
        y)
          func_libnss_rpm_build && \
          func_package_build "rpm" && \
          cp -pr ../libnss_stns/binary/*.rpm binary && \
          func_repo_build "yum_repo"
            ;;
        a)
          func_libnss_deb_build && \
          func_package_build "deb" && \
          cp -pr ../libnss_stns/binary/*.deb binary && \
          func_repo_build "apt_repo"
            ;;
    esac
done
shift $((OPTIND - 1))
