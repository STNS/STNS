# STNSインストールガイド

## はじめに
STNSはTOML形式の設定ファイルでシンプルなLinuxユーザー管理を可能にします。サーバとクライアントモジュールで構成され、少ない手順で導入可能です。またLDAPなどの既存の仕組みと併用することが出来ます。おすすめの構成としてはnginxと組み合わせ、letsencryptなどの証明証を利用し、SSL通信で暗号化して利用することをおすすめします。

## インストール
rhel,debian系共にリポジトリを提供しています。

### redhat/centos
```
$ curl -fsSL https://repo.stns.jp/scripts/yum-repo.sh | sh
```
### debian/ubuntu
```
$ curl -fsSL https://repo.stns.jp/scripts/apt-repo.sh | sh
```

サーバとクライントをインストールします。サーバは通常は専用のホストを設けて運用すると思いますが、手順書の便宜上同じホストと言う前提で記載しています。
クライアントは`libnss-stns`と`libpam-stns`に分離しています。また本手順はrhel系のコマンドを利用しますが、debian系でも同等の作業可能です。

STNSに加えて名前解決のキャッシュを行うためにnscdをインストールします。
```
$ yum install stns libnss-stns libpam-stns nscd
```

## 設定ファイル
### サーバ
インストールが完了したらサーバの設定から行います。

* /etc/stns/stns.conf

```toml
port = 1104
include = "/etc/stns/conf.d/*"

user = "test_user"
password = "test_password"

[users.example]
id = 1001
group_id = 1001
keys = ["ssh-rsa XXXXX…"]

[groups.example]
id = 1001
users = ["example"]

```

ベーシック認証を利用し、1104ポートで起動、`/etc/stns/conf.d/*`配下の設定ファイルを追加で読み込むように設定する例です。部署やチームごとに設定ファイルを分離し、運用するのが良いでしょう。
例ではユーザーexample、グループexampleを定義しています。
設定を記述したらreloadしてください。

```
$ service stns reload
```

### クライアント
まずnscdの設定を行います。ユーザーとグループ名の名前解決をキャッシュし、それ以外はキャッシュしない設定にしています。

* /etc/nscd.conf

```
enable-cache            passwd          yes
positive-time-to-live   passwd          180
negative-time-to-live   passwd          300
check-files             passwd          yes
shared                  group           yes

enable-cache            group           yes
positive-time-to-live   group           180
negative-time-to-live   group           300
check-files             group           yes
shared                  group           yes

enable-cache            hosts           no
enable-cache            services        no
enable-cache            netgroup        no
```

設定が完了したらnscdをreloadしてください。

```
$ service nscd reload
```

次にstnsの設定を行います。

* /etc/stns/libnss_stns.conf


```toml
api_end_point = ["http://<server-ip>:1104"]

user = "test_user"
password = "test_password"

wrapper_path = "/usr/local/bin/stns-query-wrapper"

chain_ssh_wrapper = "/usr/libexec/openssh/ssh-ldap-wrapper"

ssl_verify = true
```

設定としてはサーバのエンドポイント、ベーシック認証のID、パスワードを定義しています。`chain_ssh_wrapper`についてはSSHログイン時の公開鍵を取得する際にSTNSに加えて取得先がある場合に取得コマンドを定義します。定義されたコマンドに`ユーザー名`を引数に渡して取得を試みます。
また`ssl_verify`についてはSTNSサーバをnginxなどと組み合わせてSSL対応した際に証明証の照合エラーを無視するか否かの設定です。`false`に設定した場合に証明証のエラーを無視します。

* /etc/nssswitch.conf

```
passwd:     files stns
shadow:     files stns
group:      files stns
```

nsswitch.confにstnsを追加し、stns経由での名前解決を有効にします。ldapを利用している場合は`passwd:     files stns ldap`のように記載することで併用可能です。
この時点で下記のように名前解決が出来ない場合はnscdがネガティブキャッシュしている可能性があるのでキャッシュを削除してください。

```
$  id example                                                                                                                                                                
uid=1001(example) gid=1001(example) groups=1001(example)
```

キャッシュの削除
```
$ /usr/sbin/nscd -i passwd
```

最後にSSHログインを可能にするため、sshdの設定を行います。

* /etc/sshd/sshd_config

```
PubkeyAuthentication yes
AuthorizedKeysCommand /usr/local/bin/stns-key-wrapper
AuthorizedKeysCommandUser root
```

公開鍵認証を許可し、公開鍵取得コマンドに`/usr/local/bin/stns-key-wrapper`を設定してください。設定後sshdを再起動しましょう。

```
service sshd restart
```

以上でSTNSのインストールは完了です。
