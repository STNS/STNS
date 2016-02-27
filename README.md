# Simple Toml Name Service
[![Build Status](https://travis-ci.org/pyama86/STNS.svg?branch=master)](https://travis-ci.org/pyama86/STNS)

STNS is used by sshd to access keys and user resolver provided
client library:https://github.com/pyama86/libnss_stns

```
$ ssh pyama@example.jp                                                                                                                                                   $ id pyama
uid=1001(pyama) gid=1001(pyama) groups=1001(pyama)
```

diagram
![diagram](https://cloud.githubusercontent.com/assets/8022082/13373739/362ca2c8-ddb3-11e5-97e2-13ea1269c26e.png)



## install
download page <https://github.com/pyama86/STNS/releases>
```
$ wget https://github.com/pyama86/STNS/releases/download/<version>/stns-<version>.noarch.rpm
$ rpm -ivh stns-<version>.noarch.rpm
$ service stns start
```

## config
* /etc/stns/stns.conf
```
port = 1104
[users.example]
id = 1001
group_id = 1001
directory = "/home/example"
shell = "/bin/bash"
keys = ["ssh-rsa XXXXX…"]
[groups.example]
id = 1001
users = ["example"]
```
support format /etc/passwd,/etc/groups,/etc/shadow

## author
* pyama86
