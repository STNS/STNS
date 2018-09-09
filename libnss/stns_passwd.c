#include "stns.h"

static json_t *entries = NULL;
static int entry_idx   = 0;

#define PASSWD_ENSURE(entry)                                                                                           \
  const json_int_t id       = json_integer_value(json_object_get(entry, "id"));                                        \
  const json_int_t group_id = json_integer_value(json_object_get(entry, "group_id"));                                  \
  const char *name          = json_string_value(json_object_get(entry, "name"));                                       \
  const char *gecos         = json_string_value(json_object_get(entry, "gecos"));                                      \
  const char *dir           = json_string_value(json_object_get(entry, "directory"));                                  \
  const char *shell         = json_string_value(json_object_get(entry, "shell"));                                      \
  char passwd[]             = "x";                                                                                     \
  rbuf->pw_uid              = c->uid_shift + id;                                                                       \
  rbuf->pw_gid              = c->gid_shift + group_id;                                                                 \
                                                                                                                       \
  STNS_SET_DEFAULT_VALUE(sh, shell, "/bin/bash");                                                                      \
  char b[MAXBUF];                                                                                                      \
  sprintf(b, "/home/%s", name);                                                                                        \
  STNS_SET_DEFAULT_VALUE(d, dir, b);                                                                                   \
  SET_ATTRBUTE(pw, name, name)                                                                                         \
  SET_ATTRBUTE(pw, passwd, passwd)                                                                                     \
  SET_ATTRBUTE(pw, gecos, gecos)                                                                                       \
  SET_ATTRBUTE(pw, dir, dir)                                                                                           \
  SET_ATTRBUTE(pw, shell, shell)

STNS_ENSURE_BY(name, const char *, user_name, string, name, (strcmp(current, user_name) == 0), passwd, PASSWD)
STNS_ENSURE_BY(uid, uid_t, uid, integer, id, current == uid, passwd, PASSWD)

STNS_GET_SINGLE_VALUE_METHOD(getpwnam_r, const char *name, "users?name=%s", name, passwd, )
STNS_GET_SINGLE_VALUE_METHOD(getpwuid_r, uid_t uid, "users?id=%d", uid, passwd, USER_ID_QUERY_AVAILABLE)
STNS_SET_ENTRIES(pw, PASSWD, passwd, users)
