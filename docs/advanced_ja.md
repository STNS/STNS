# STNSアドバンスドガイド
## SSHログイン時にホームディレクトリを作成する
sshログインした際に、ユーザーのホームディレクトリは作成されません。`/etc/pam.d/sshd`に設定を行うことでホームディレクトリを自動で作成することが出来ます。

```
$ echo 'session    required     pam_mkhomedir.so skel=/etc/skel/ umask=0022' >> /etc/pam.d/sshd
```

## link_users機能を利用し、デプロイユーザーを追加する
Webサービスにおいて、アプリケーションをデプロイする際に、デプロイユーザーを作成し、デプロイする運用が多くあると思います。大体のケースにおいて、デプロイユーザーの`~/.ssh/authorized_keys`にデプロイするユーザーの公開鍵を羅列して運用することが多いでしょう。STNSの`link_users`機能を利用すると同様のことがシンプルに実現できます。

```toml
[users.deploy]
id = 1
group_id = 1
link_users = ["example1","example2"]

[users.example1]
id = 2
group_id = 1
keys = ["ssh-rsa aaa"]

[users.example2]
id = 3
group_id = 1
keys = ["ssh-rsa bbb"]
```

上記の定義を行うと、`deploy`というユーザー名で`exampl1`と`example2`がログイン可能となります。裏の動きとしては`deploy`というユーザー名でSSHログインを行う際に、example1とexample2の公開鍵である`ssh-rsa aaa`および`ssh-rsa bbb`を返却しています。

## link_groups機能を利用し、組織の階層構造を表現する
サーバのパーミッション制御を行う際に、部署単位で行いたい場合があります。例えばA課で読み取り書き込み権限を与え、A課が所属するB部では読み取り権限のみ与えると行ったケースです。STNSではこういったケースにも対応する事ができます。

```toml
[groups.department]
users = ["user1"]
link_groups = ["division"]

[groups.division]
users = ["user2"]

```

この例では`department`グループに`user1`が所属しており、`division`グループには`user2`が所属しています。また`department`には`link_groups`に`division`を定義しております。
これによって、departmentにはdivisionも所属しているという状態になり、`id`コマンドを発行すると下記のようになります。

```
$ id user2
uid=1001(user1) gid=1002(division) groups=1001(department),1002(division)
```

これによりuser2はdepartmentに所属するdivisionのユーザーということを表現できます。

## sudoパスワードを管理する
STNSでは2種類の方法でsudoのパスワードを管理することが出来ます。

1. sudo専用のアカウントを設け、共通パスワードとして管理する。
2. 従来通りのユーザーごとのパスワードを利用する。

### sudo専用のアカウントを利用する
STNSにsudo用のアカウントを設け、パスワードを管理することが出来ます。イメージとしては第2のrootパスワードです。

下記のようにサーバにsudo用の定義を行います。
* /etc/stns/stns.conf
```toml
[sudoers.example]
password = "f2ca1bb6c7e907d06dafe4687e579fce76b37e4e93b7605022da52e6ccc26fd2"
hash_type = "sha256"
```

`hash_type`には`sha256`と`sha512`が指定可能です。パスワードハッシュについては下記のように生成すると良いでしょう。

```
$ echo -n "test" | sha256sum
f2ca1bb6c7e907d06dafe4687e579fce76b37e4e93b7605022da52e6ccc26fd2
```

次にクライアントのpamの設定を行います。

* /etc/pam.d/sudo
```
#%PAM-1.0
auth       sufficient   libpam_stns.so sudo example
auth       include      system-auth
account    include      system-auth
password   include      system-auth
session    optional     pam_keyinit.so revoke
session    required     pam_limits.so
```

ポイントとしては`libpam_stns.so`の2つの引数です。ひとつ目の引数`sudo`でsudo用のアカウントを利用することを指定し、ふたつ目の引数でアカウント`example`を利用することを指定しています。
この状態で`/etc/sudoers`に正しく設定が行われていれば、sudoする際に`[sudoers.example]`で定義したパスワードを入力することにより、sudoコマンドを利用することが出来ます。

### ユーザー毎のパスワードを利用する
STNSではユーザーごとにパスワードハッシュを定義することが出来ます。

```
[users.example]
id = 1000
group_id = 1000
directory = "/home/example"
password = "f2ca1bb6c7e907d06dafe4687e579fce76b37e4e93b7605022da52e6ccc26fd2"
hash_type = "sha256"
```

このように定義した状態でクライアント側のpamを下記のように定義します。

* /etc/pam.d/sudo

```
#%PAM-1.0
auth       sufficient   libpam_stns.so
auth       include      system-auth
account    include      system-auth
password   include      system-auth
session    optional     pam_keyinit.so revoke
session    required     pam_limits.so
```

前項と異なり`libpam_stns.so`に引数を与えない場合は、Linuxからユーザー名を取得し比較を行うため、ユーザーに定義したパスワードハッシュでsudo時の認証を行うことが可能です。

また下記のように定義することによりログインなどの認証もSTNSに定義したユーザーパスワードで行うことが出来ます。

* /etc/pam.d/system-auth(rhel) or /etc/pam.d/common-auth(debian)

```
#%PAM-1.0
#%PAM-1.0
# This file is auto-generated.
# User changes will be destroyed the next time authconfig is run.
auth        required      pam_env.so
auth        sufficient    libpam_stns.so 
…
```
