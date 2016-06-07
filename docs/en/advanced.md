# STNS Advanced Guide

## Create home directory for login user

If you want the system to create home directory for login user after login in via SSH, add the configuration below into `/etc/pam.d/sshd`:

```
$ echo 'session    required     pam_mkhomedir.so skel=/etc/skel/ umask=0022' >> /etc/pam.d/sshd
```

## Add deploy users using `link_users`

We often create a special user account whose privilege is restricted for a certain purpose that we deploy application to remote server by the user account. To realize such a way to deploy, we can add real user accounts' SSH public keys to the user account's `~/.ssh/authorized_keys`.

STNS allows us to easily bundle deploy users using the `link_users` feature.

```toml
[users.deploy]
id         = 1
group_id   = 1
link_users = ["example1","example2"]

[users.example1]
id         = 2
group_id   = 1
keys       = ["ssh-rsa aaa"]

[users.example2]
id         = 3
group_id   = 1
keys       = ["ssh-rsa bbb"]
```

With the configuration above, the users shown as `example1` and `example2` can login to a server as `deploy` user.

When the `deploy` user attempts to login to a server via SSH, STNS server returns to STNS client keys for both `example1` (`ssh-rsa aaa`) and `example2` (`ssh-rsa bbb`). That is, the users are linked.

## Express organizational structure using `link_groups`

It's a common requirement to control who can access to a server with respect to organizational structure. Imagine such a case like that we want to give read and write privilege to ones in a division, and, at the same time, we want to add read-only privilege to ones in a department to which the division belong.

STNS meets such a complicated requiment.

```toml
[groups.department]
users       = ["user1"]
link_groups = ["division"]

[groups.division]
users       = ["user2"]
```

In the configuration above, `user1` belongs to `department` group and `user2` belongs to `division` group. `link_group` is set to link `department` group to `division` group.

As a result:

* `department` group includes `division` group
* `user2` belongs to both `department` group and `division` group

You can confirm it by executing `id` command like below:

```
$ id user2
uid=1001(user1) gid=1002(division) groups=1001(department),1002(division)
```

## Manage passwords for `sudo`

STNS privides two ways to manage password for `sudo`:

1. Use a common password for a special user account for `sudo`
2. Authenticate users by their passwords as usual

### 1. Use a common password for a special user account for `sudo`

#### Server

You can create a special user account for `sudo` like below. It can be considered like a secondary root password.

`/etc/stns/stns.conf`:

```toml
salt_enable       = true
stretching_number = 100000

[sudoers.example]
password          = "a3b20fc634ac4bad5be8a40566acb00adcd2e5bf2fb9be4750150553d529b799"
hash_type         = "sha256"
```

The configuration above:

* Enables salt and does 100 thousands of stretching to secure the password
* Utilizes SHA256 algorithm to hash the password (`hash_type` can be set to `sha256` or `sha512`)

You can use [stns-passwd](https://github.com/STNS/stns-passwd) to get such a hashed password.

#### Client

Configuration for STNS client is needed in turn.

`/etc/pam.d/sudo`:

```
#%PAM-1.0
auth       sufficient   libpam_stns.so sudo example
auth       include      system-auth
account    include      system-auth
password   include      system-auth
session    optional     pam_keyinit.so revoke
session    required     pam_limits.so
```

You have to notice the two arguments which follows `libpam_stns.so`. The first one is to enable `sudo` authorization and the second one is to set the user name for `sudo` (`example` user is set at above).

If you configure `/etc/sudoers` correctly, you will successfully be authorized by `sudo` using the password which is set in `[sodoers.example]` section in the STNS server configuration.

### 2. Authenticate users by their passwords as usual

#### Server

You can also set hashed passwords for each users like below:

`/etc/stns/stns.conf`:

```toml
[users.example]
id        = 1000
group_id  = 1000
directory = "/home/example"
password  = "f2ca1bb6c7e907d06dafe4687e579fce76b37e4e93b7605022da52e6ccc26fd2"
hash_type = "sha256"
```

#### Client

Client can be configured as below:

`/etc/pam.d/sudo`:

```
#%PAM-1.0
auth       sufficient   libpam_stns.so
auth       include      system-auth
account    include      system-auth
password   include      system-auth
session    optional     pam_keyinit.so revoke
session    required     pam_limits.so
```

You have to notice, at this time, that there's no arguments which follows `libpam_stns.so`. If no arguments set, STNS uses Linux user account and process the authorization using the hashed password set in the server configuration.

In other way, you can use the hashed password for each users set in STNS server configuration for login authentication.

`/etc/pam.d/system-auth` (RHEL) or `/etc/pam.d/common-auth` (Debian Family):

```
#%PAM-1.0
#%PAM-1.0
# This file is auto-generated.
# User changes will be destroyed the next time authconfig is run.
auth        required      pam_env.so
auth        sufficient    libpam_stns.so
…
```
