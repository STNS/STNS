#include "stns.h"

static json_t *entries = NULL;
static int entry_idx   = 0;

#define GROUP_GET_SINGLE(method, first, format, value)                                                                 \
  enum nss_status _nss_stns_##method(first, struct group *rbuf, char *buf, size_t buflen, int *errnop)                 \
  {                                                                                                                    \
    int curl_result;                                                                                                   \
    stns_http_response_t r;                                                                                            \
    stns_conf_t c;                                                                                                     \
    char url[MAXBUF];                                                                                                  \
                                                                                                                       \
    stns_load_config(STNS_CONFIG_FILE, &c);                                                                            \
                                                                                                                       \
    sprintf(url, format, value);                                                                                       \
    curl_result = stns_request(&c, url, &r);                                                                           \
                                                                                                                       \
    if (curl_result != CURLE_OK) {                                                                                     \
      return NSS_STATUS_UNAVAIL;                                                                                       \
    }                                                                                                                  \
                                                                                                                       \
    return ensure_group_by_##value(r.data, &c, value, rbuf, buf, buflen, errnop);                                      \
  }

#define GROUP_SET_ATTRBUTE(name)                                                                                       \
  int name##_length = strlen(name) + 1;                                                                                \
                                                                                                                       \
  if (buflen < name##_length) {                                                                                        \
    *errnop = ERANGE;                                                                                                  \
    return NSS_STATUS_TRYAGAIN;                                                                                        \
  }                                                                                                                    \
                                                                                                                       \
  strcpy(buf, name);                                                                                                   \
  rbuf->gr_##name = buf;                                                                                               \
  buf += name##_length;                                                                                                \
  buflen -= name##_length;

#define GROUP_ENSURE(group)                                                                                            \
  const json_int_t id = json_integer_value(json_object_get(group, "id"));                                              \
  const char *name    = json_string_value(json_object_get(group, "name"));                                             \
  char passwd[]       = "x";                                                                                           \
                                                                                                                       \
  rbuf->gr_gid = c->gid_shift + id;                                                                                    \
                                                                                                                       \
  GROUP_SET_ATTRBUTE(name)                                                                                             \
  GROUP_SET_ATTRBUTE(passwd)                                                                                           \
  rbuf->gr_mem = (char **)buf;                                                                                         \
                                                                                                                       \
  json_t *members = json_object_get(group, "users");                                                                   \
  int i;                                                                                                               \
  for (i = 0; i < json_array_size(members); i++) {                                                                     \
    json_t *member = json_array_get(members, i);                                                                       \
    if (!json_is_string(member)) {                                                                                     \
      return NSS_STATUS_UNAVAIL;                                                                                       \
    }                                                                                                                  \
    const char *user = json_string_value(member);                                                                      \
    if (buflen <= strlen(user)) {                                                                                      \
      return NSS_STATUS_TRYAGAIN;                                                                                      \
    }                                                                                                                  \
    rbuf->gr_mem[i] = strdup(user);                                                                                    \
    buf += strlen(rbuf->gr_mem[i]) + 1;                                                                                \
    buflen -= strlen(rbuf->gr_mem[i]) + 1;                                                                             \
  }                                                                                                                    \
  rbuf->gr_mem[json_array_size(members)] = NULL;

#define GROUP_ENSURE_BY(method_key, key_type, key_name, json_type, json_key, match_method)                             \
  enum nss_status ensure_group_by_##method_key(char *data, stns_conf_t *c, key_type key_name, struct group *rbuf,      \
                                               char *buf, size_t buflen, int *errnop)                                  \
  {                                                                                                                    \
    int i;                                                                                                             \
    json_error_t error;                                                                                                \
    json_t *group;                                                                                                     \
    json_t *groups = json_loads(data, 0, &error);                                                                      \
                                                                                                                       \
    if (groups == NULL) {                                                                                              \
      syslog(LOG_ERR, "%s[L%d] json parse error: %s", __func__, __LINE__, error.text);                                 \
      goto leave;                                                                                                      \
    }                                                                                                                  \
                                                                                                                       \
    json_array_foreach(groups, i, group)                                                                               \
    {                                                                                                                  \
      if (!json_is_object(group)) {                                                                                    \
        continue;                                                                                                      \
      }                                                                                                                \
      key_type current = json_##json_type##_value(json_object_get(group, #json_key));                                  \
                                                                                                                       \
      if (match_method) {                                                                                              \
        GROUP_ENSURE(group)                                                                                            \
        free(data);                                                                                                    \
        json_decref(groups);                                                                                           \
        return NSS_STATUS_SUCCESS;                                                                                     \
      }                                                                                                                \
    }                                                                                                                  \
                                                                                                                       \
    free(data);                                                                                                        \
    json_decref(groups);                                                                                               \
    return NSS_STATUS_NOTFOUND;                                                                                        \
  leave:                                                                                                               \
    return NSS_STATUS_UNAVAIL;                                                                                         \
  }

GROUP_ENSURE_BY(name, const char *, group_name, string, name, (strcmp(current, group_name) == 0))
GROUP_ENSURE_BY(gid, gid_t, gid, integer, id, current == gid)

GROUP_GET_SINGLE(getgrnam_r, const char *name, "groups?name=%s", name)
GROUP_GET_SINGLE(getgrgid_r, gid_t gid, "groups?id=%d", gid)
STNS_SET_ENTRIES(gr, GROUP, group, groups)
