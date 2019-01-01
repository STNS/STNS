# Simple Toml Name Service
[![Build Status](https://travis-ci.org/STNS/STNS.svg?branch=master)](https://travis-ci.org/STNS/STNS)

STNS is used by sshd to access keys and user/group resolver provided

You can see the details of STNS on [stns.jp](http://stns.jp)

```
$ ssh pyama@example.jp
$ id pyama
uid=1001(pyama) gid=1001(pyama) groups=1001(pyama)
```

# diagram
![overview](https://raw.githubusercontent.com/STNS/STNS/master/docs/images/diagram.png)

## blog
* [Linuxユーザーと公開鍵を統合管理するサーバ&クライアントを書いた](https://ten-snapon.com/archives/1228)
* [デプロイユーザーをSTNSで管理する](https://ten-snapon.com/archives/1330)
* [STNSに組織体系を管理するLinkGroup機能を追加しi386に対応しました](https://ten-snapon.com/archives/1346)
* [STNSでSudoパスワードをサポートした](https://ten-snapon.com/archives/1355)
* [パスワード暗号化について学びを得た](https://ten-snapon.com/archives/1399)

# VS
## LDAP
LDAP is used convenient and very well
However, sometimes it becomes complicated and versatile too.
STNS function is small compared with the LDAP, but it is management that much simple.
And, In many cases, meet the required functionality.

## How to Contribute
Please give me a pull request anything!

### Test
#### Server
```
$ make depsdev
$ make test
```
## author
* pyama86
