# Simple Toml Name Service
simple toml name service is Linux `/etc/passwd`,`/etc/group`,`/etc/shadow` name resolution from toml format config
> now support is x86 rhel linux server

client library:https://github.com/pyama86/libnss_stns

## install
donload page <https://github.com/pyama86/STNS/releases>
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
group_id = 1
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
