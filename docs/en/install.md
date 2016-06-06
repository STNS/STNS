# STNS Installation Guide

## What's SNTS?

STNS allows you to easily manage Linux users with simple TOML-based configuration. It consists of server and client implementation, which requires only a few steps to install. Moreover, you can use it with existing user management systems such as LDAP.

It's highly recommended that you adopt such a system architecture like that STNS server sits behind Nginx which handles and terminates TLS session using a free certificate provided by [Let's Encrypt](https://letsencrypt.org/).

## Installation

We provides both YUM and APT repositories.

You can install YUM/APT repository by executing the command below. Although you must want to install SNTS server and client into separate hosts, we here install them into the same host for the sake of simplicity of explanation.

### Install YUM Repository

```
$ curl -fsSL https://repo.stns.jp/scripts/yum-repo.sh | sh
```

### Install APT Repository

```
$ curl -fsSL https://repo.stns.jp/scripts/apt-repo.sh | sh
```

### Install STNS Server and Client

Alhough we use commands for RHEL in what follows, the equivalent commands will work also for Debian family.

You can install STNS server by installing the stns package. The STNS client consists of two programs: `libnss-stns` and `libpam-stns`. Additionally, you have to install nscd to cache result of name resolution.

```
$ yum install stns libnss-stns libpam-stns nscd
```

## Configuration

### STNS Server

After successfully installing packages, let's configure STNS server.

`/etc/stns/stns.conf`:

```toml
port     = 1104
include  = "/etc/stns/conf.d/*"

user     = "test_user"
password = "test_password"

[users.example]
id       = 1001
group_id = 1001
keys     = ["ssh-rsa XXXXX…"]

[groups.example]
id       = 1001
users    = ["example"]
```

This configuration means that the STNS server:

* Utilizes basic authentication to restrict requests.
* Listens to 1104 port.
* Includes additional configurations placed under `/etc/stns/conf.d`.
* Defines example user and group.

We encourage you to set configurations for each teams into separate files.

Reload the server right after modifying the file to activate the new configuration.

```
$ service stns reload
```

### STNS Client

**Firstly**, configure nscd to cache result of user and group names resolution. This configuration means that the system caches only user and group names.

`/etc/nscd.conf`:

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

Reload nscd right after modifying the file to activate the new configuration.

```
$ service nscd reload
```

**Secondly**, configure STNS client.

`/etc/stns/libnss_stns.conf`:

```toml
api_end_point     = ["http://<server-ip>:1104"]

user              = "test_user"
password          = "test_password"

wrapper_path      = "/usr/local/bin/stns-query-wrapper"

chain_ssh_wrapper = "/usr/libexec/openssh/ssh-ldap-wrapper"

ssl_verify        = true
```

This file conigures the location of SNTS server, and the combination of id and password for basic authentication.

You can set a script path by `chain_ssh_wrapper` to retrieve SSH public key from other place except for STNS server. STNS client executes the script with a user name as an argument.

`ssl_verify` tells if the client must verify or not the TLS certificate in the negotiation process with STNS server. If `false` is set, the client ignores the verification error of TLS certificate.

**Thirdly**, configure the name resolution order like below.

`/etc/nssswitch.conf`:

```
passwd:     files stns
shadow:     files stns
group:      files stns
```

Add snts into `nsswitch.conf` to enable name resolution using STNS. To use LDAP concurrently, you can configure like: `passwd: files stns ldap`.

If name resolution fails like below, purge caches in case that nscd has negative caches.

```
$  id example
uid=1001(example) gid=1001(example) groups=1001(example)
```

You can purge negative caches like below:

```
$ /usr/sbin/nscd -i passwd
```

**Lastly**, configure sshd to enable SSH login using STNS.

`/etc/sshd/sshd_config`:

```
PubkeyAuthentication yes
AuthorizedKeysCommand /usr/local/bin/stns-key-wrapper
AuthorizedKeysCommandUser root
```

This configuration means that tha SSH server:

* Activates SSH Public key authentication.
* Use `/usr/local/bin/stns-key-wrapper` to retrieve the public key for the login user.

Reload sshd right after the modifiying the coniguration file.

```
service sshd restart
```

Installation of STNS has been finally completed!
