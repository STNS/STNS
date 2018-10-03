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

  assertEquals \
    "big_group:x:100:member0,member1,member2,member3,member4,member5,member6,member7,member8,member9,member10,member11,member12,member13,member14,member15,member16,member17,member18,member19,member20,member21,member22,member23,member24,member25,member26,member27,member28,member29,member30,member31,member32,member33,member34,member35,member36,member37,member38,member39,member40,member41,member42,member43,member44,member45,member46,member47,member48,member49,member50,member51,member52,member53,member54,member55,member56,member57,member58,member59,member60,member61,member62,member63,member64,member65,member66,member67,member68,member69,member70,member71,member72,member73,member74,member75,member76,member77,member78,member79,member80,member81,member82,member83,member84,member85,member86,member87,member88,member89,member90,member91,member92,member93,member94,member95,member96,member97,member98,member99,member100,member101,member102,member103,member104,member105,member106,member107,member108,member109,member110,member111,member112,member113,member114,member115,member116,member117,member118,member119,member120,member121,member122,member123,member124,member125,member126,member127,member128,member129,member130,member131,member132,member133,member134,member135,member136,member137,member138,member139,member140,member141,member142,member143,member144,member145,member146,member147,member148,member149,member150,member151,member152,member153,member154,member155,member156,member157,member158,member159,member160,member161,member162,member163,member164,member165,member166,member167,member168,member169,member170,member171,member172,member173,member174,member175,member176,member177,member178,member179,member180,member181,member182,member183,member184,member185,member186,member187,member188,member189,member190,member191,member192,member193,member194,member195,member196,member197,member198,member199" \
    "$(getent group big_group)"
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
  sudo -u test true
  assertTrue $?
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
