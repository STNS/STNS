#!/bin/bash

test_id()
{
  assertEquals \
    "uid=10001(test) gid=0(root) groups=0(root),10001(test)" \
    "$(id test)"
  assertEquals \
    "" \
    "$(id notfound)"
}

test_getent_passwd()
{
  assertEquals \
    "test:x:10001:0::/home/test:/bin/bash" \
    "$(getent passwd | grep test)"

  assertEquals \
    "foo:x:10002:0::/home/foo:/bin/bash" \
    "$(getent passwd | grep foo)"

  assertEquals \
    "test:x:10001:0::/home/test:/bin/bash" \
    "$(getent passwd test)"
}

test_getent_group()
{
  assertEquals \
    "test:x:10001:test" \
    "$(getent group | grep test)"

  assertEquals \
    "bar:x:10002:foo" \
    "$(getent group | grep bar)"

  assertEquals \
    "test:x:10001:test" \
    "$(getent group test)"
}

test_getent_shadow()
{
  assertEquals \
    "test:test:::::::" \
    "$(getent shadow | grep 'test:test')"

  assertEquals \
    "foo:test:::::::" \
    "$(getent shadow | grep foo)"

  assertEquals \
    "test:test:::::::" \
    "$(getent shadow test)"
}

test_sudo()
{
  assertTrue `sudo -u test ls`
}

test_key_wrapper()
{
  assertEquals \
    "key1
key2
aaabbbccc
ddd" \
    "$(tmp/libs/stns-key-wrapper test test)"

  assertEquals \
    "http request failed user: fuga" \
    "$((tmp/libs/stns-key-wrapper fuga)2>&1)"

  assertEquals \
    "User name is a required parameter" \
    "$((tmp/libs/stns-key-wrapper)2>&1)"
}


. /usr/include/shunit2/src/shunit2
