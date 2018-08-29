#!/bin/sh

set -e

sudo -k

echo "This script requires superuser authority to configure stns apt repository:"

sudo sh <<'SCRIPT'
  set -x
  echo "deb https://repo.stns.jp/debian/ stns main" > /etc/apt/sources.list.d/stns.list
  curl -fsS https://repo.stns.jp/gpg/GPG-KEY-stns| apt-key add -
  apt-get update -qq
SCRIPT

echo 'done'
