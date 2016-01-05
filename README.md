# Simple Toml Name Service
simple toml name service is Linux `/etc/passwd`,`/etc/group` name resolution from toml format config
> now support is x86 rhel linux server

## install
donload page <https://github.com/pyama86/STNS/releases>
```
$ wget https://github.com/pyama86/STNS/releases/download/<version>/stns-<version>.noarch.rpm
$ rpm -ivh stns-0.1-1.noarch.rpm
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
support format /etc/passwd and /etc/groups

## author
* pyama86
