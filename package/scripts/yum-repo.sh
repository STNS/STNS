#!/bin/sh

set -e

sudo -k

echo "This script requires superuser authority to configure stns yum repository:"

sudo sh <<'SCRIPT'
  set -x

  # import GPG key
  gpgkey_path=`mktemp`
  curl -fsS -o $gpgkey_path https://repo.stns.jp/gpg/GPG-KEY-stns
  rpm --import $gpgkey_path
  rm $gpgkey_path

  # add config for stns yum repos
  cat >/etc/yum.repos.d/stns.repo <<'EOF';
[stns]
name=stns
baseurl=https://repo.stns.jp/centos/$basearch/$releasever
gpgcheck=1
EOF
SCRIPT

echo 'done'
