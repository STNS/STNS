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

. /usr/include/shunit2/src/shunit2
