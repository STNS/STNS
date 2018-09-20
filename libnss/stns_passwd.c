#include "stns.h"

static JSON_Value *entries  = NULL;
static int entry_idx        = 0;
pthread_mutex_t pwent_mutex = PTHREAD_MUTEX_INITIALIZER;

#define PASSWD_ENSURE(entry)                                                                                           \
  int id            = (int)json_value_get_number(json_object_get_value(entry, "id"));                                  \
  int group_id      = (int)json_value_get_number(json_object_get_value(entry, "group_id"));                            \
  const char *name  = json_value_get_string(json_object_get_value(entry, "name"));                                     \
  const char *gecos = json_value_get_string(json_object_get_value(entry, "gecos"));                                    \
  const char *dir   = json_value_get_string(json_object_get_value(entry, "directory"));                                \
  const char *shell = json_value_get_string(json_object_get_value(entry, "shell"));                                    \
  char passwd[]     = "x";                                                                                             \
  rbuf->pw_uid      = c->uid_shift + id;                                                                               \
  rbuf->pw_gid      = c->gid_shift + group_id;                                                                         \
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
STNS_ENSURE_BY(uid, uid_t, uid, number, id, current + (c->uid_shift) == uid, passwd, PASSWD)

STNS_GET_SINGLE_VALUE_METHOD(getpwnam_r, const char *name, "users?name=%s", name, passwd, , )
STNS_GET_SINGLE_VALUE_METHOD(getpwuid_r, uid_t uid, "users?id=%d", uid, passwd, USER_ID_QUERY_AVAILABLE, -(c.uid_shift))
STNS_SET_ENTRIES(pw, PASSWD, passwd, users)
