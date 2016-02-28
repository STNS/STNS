# Simple Toml Name Service
[![Build Status](https://travis-ci.org/STNS/STNS.svg?branch=master)](https://travis-ci.org/STNS/STNS)

STNS is used by sshd to access keys and user resolver provided

client library:https://github.com/pyama86/libnss_stns

```
$ ssh pyama@example.jp
$ id pyama
uid=1001(pyama) gid=1001(pyama) groups=1001(pyama)
```

diagram

![diagram](https://cloud.githubusercontent.com/assets/8022082/13373974/250a8b16-ddba-11e5-994d-b1bbc81a6b94.png)

blog
* [Linuxユーザーと公開鍵を統合管理するサーバ&クライアントを書いた](https://ten-snapon.com/archives/1228)

## install
## redhat/centos
```
$ curl -fsSL https://repo.stns.jp/scripts/yum-repo.sh | sh
$ yum install stns
```
## debian/ubuntu
```
$ curl -fsSL https://repo.stns.jp/scripts/apt-repo.sh | sh
$ apt-get install stns
```

## config
* /etc/stns/stns.conf

```toml
port = 1104
include = "/etc/stns/conf.d/*"

# support basic auth
user = "basic_user"
password = "basic_password"

[users.example]
id = 1001
group_id = 1001
directory = "/home/example" # default:/home/:user_name
shell = "/bin/bash" # default:/bin/bash
keys = ["ssh-rsa XXXXX…"]
link_users = ["foo"]

[groups.example]
id = 1001
users = ["example"]
```

### General
|Name|Description|
|---|---|
|port|listen port|
|include|include config directory|
|user| basic authentication user|
|password| basic authentication password|

### Users
|Name|Description|
|---|---|
|id| unique user id|
|group_id|id of the group they belong|
|directory|home directory path|
|shell|default shell path|
|keys|public key list|
|link_users|merge public key from the specified user|

#### link_users
link_users params is merge public key from the specified user

```toml
[users.pyama1]
keys = ["ssh-rsa aaa"]
link_users = ["pyama2"] ←

[users.pyama2]
keys = ["ssh-rsa bbb"]
```
```
$ /user/local/bin/stns-key-wrapper pyama1
ssh-rsa aaa
ssh-rsa bbb
```

### Groups
|Name|Description|
|---|---|
|id| unique group id|
|users|user id of the members|



## author
* pyama86
