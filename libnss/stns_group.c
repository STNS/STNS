#include "stns.h"

static json_t *entries = NULL;
static int entry_idx   = 0;

#define GROUP_ENSURE(group)                                                                                            \
  const json_int_t id = json_integer_value(json_object_get(group, "id"));                                              \
  const char *name    = json_string_value(json_object_get(group, "name"));                                             \
  char passwd[]       = "x";                                                                                           \
                                                                                                                       \
  rbuf->gr_gid = c->gid_shift + id;                                                                                    \
                                                                                                                       \
  SET_ATTRBUTE(gr, name, name)                                                                                         \
  SET_ATTRBUTE(gr, passwd, passwd)                                                                                     \
  rbuf->gr_mem = (char **)buf;                                                                                         \
                                                                                                                       \
  json_t *members = json_object_get(group, "users");                                                                   \
  int i;                                                                                                               \
  int ptr_area_size = (json_array_size(members) + 1) * sizeof(char *);                                                 \
  char *next_member;                                                                                                   \
                                                                                                                       \
  if (buflen < ptr_area_size) {                                                                                        \
    (*errnop) = ERANGE;                                                                                                \
    return NSS_STATUS_TRYAGAIN;                                                                                        \
  }                                                                                                                    \
                                                                                                                       \
  next_member = buf + ptr_area_size;                                                                                   \
  for (i = 0; i < json_array_size(members); i++) {                                                                     \
    json_t *member = json_array_get(members, i);                                                                       \
    if (!json_is_string(member)) {                                                                                     \
      return NSS_STATUS_UNAVAIL;                                                                                       \
    }                                                                                                                  \
    const char *user = json_string_value(member);                                                                      \
    int user_length  = strlen(user) + 1;                                                                               \
    if (buflen < user_length) {                                                                                        \
      *errnop = ERANGE;                                                                                                \
      return NSS_STATUS_TRYAGAIN;                                                                                      \
    }                                                                                                                  \
    strcpy(next_member, user);                                                                                         \
    rbuf->gr_mem[i] = next_member;                                                                                     \
    next_member += user_length + 1;                                                                                    \
    buflen -= user_length + 1;                                                                                         \
  }                                                                                                                    \
  rbuf->gr_mem[json_array_size(members)] = NULL;

STNS_ENSURE_BY(name, const char *, group_name, string, name, (strcmp(current, group_name) == 0), group, GROUP)
STNS_ENSURE_BY(gid, gid_t, gid, integer, id, current == gid, group, GROUP)

STNS_GET_SINGLE_VALUE_METHOD(getgrnam_r, const char *name, "groups?name=%s", name, group, )
STNS_GET_SINGLE_VALUE_METHOD(getgrgid_r, gid_t gid, "groups?id=%d", gid, group, GROUP_ID_QUERY_AVAILABLE)
STNS_SET_ENTRIES(gr, GROUP, group, groups)
