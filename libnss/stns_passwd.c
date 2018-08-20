#include "stns.h"

pthread_mutex_t pwent_mutex = PTHREAD_MUTEX_INITIALIZER;
static json_t *entries      = NULL;
static int entry_idx        = 0;

#define PASSWD_DEFAULT(buf, name, def)                                                                                 \
  char buf[MAXBUF];                                                                                                    \
  if (name != NULL && strlen(name) > 0) {                                                                              \
    strcpy(buf, name);                                                                                                 \
  } else {                                                                                                             \
    strcpy(buf, def);                                                                                                  \
  }                                                                                                                    \
  name = buf;

#define PASSWD_GET_SINGLE(method, first, format, value)                                                                \
  enum nss_status _nss_stns_##method(first, struct passwd *pwd, char *buf, size_t buflen, int *errnop)                 \
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
    return ensure_passwd_by_##value(r.data, &c, value, pwd, buf, buflen, errnop);                                      \
  }

#define PASSWD_SET_ATTRBUTE(name)                                                                                      \
  int name##_length = strlen(name) + 1;                                                                                \
                                                                                                                       \
  if (buflen < name##_length) {                                                                                        \
    *errnop = ERANGE;                                                                                                  \
    return NSS_STATUS_TRYAGAIN;                                                                                        \
  }                                                                                                                    \
                                                                                                                       \
  strcpy(buf, name);                                                                                                   \
  pwd->pw_##name = buf;                                                                                                \
  buf += name##_length;                                                                                                \
  buflen -= name##_length;

#define PASSWD_ENSURE(user)                                                                                            \
  const json_int_t id       = json_integer_value(json_object_get(user, "id"));                                         \
  const json_int_t group_id = json_integer_value(json_object_get(user, "group_id"));                                   \
  const char *name          = json_string_value(json_object_get(user, "name"));                                        \
  const char *gecos         = json_string_value(json_object_get(user, "gecos"));                                       \
  const char *dir           = json_string_value(json_object_get(user, "directory"));                                   \
  const char *shell         = json_string_value(json_object_get(user, "shell"));                                       \
  char passwd[]             = "x";                                                                                     \
                                                                                                                       \
  pwd->pw_uid = c->uid_shift + id;                                                                                     \
  pwd->pw_gid = c->gid_shift + group_id;                                                                               \
                                                                                                                       \
  PASSWD_DEFAULT(sh, shell, "/bin/bash");                                                                              \
                                                                                                                       \
  char b[MAXBUF];                                                                                                      \
  sprintf(b, "/home/%s", name);                                                                                        \
  PASSWD_DEFAULT(d, dir, b);                                                                                           \
  PASSWD_SET_ATTRBUTE(name)                                                                                            \
  PASSWD_SET_ATTRBUTE(passwd)                                                                                          \
  PASSWD_SET_ATTRBUTE(gecos)                                                                                           \
  PASSWD_SET_ATTRBUTE(dir)                                                                                             \
  PASSWD_SET_ATTRBUTE(shell)

#define PASSWD_ENSURE_BY(method_key, key_type, key_name, json_type, json_key, match_method)                            \
  enum nss_status ensure_passwd_by_##method_key(char *data, stns_conf_t *c, key_type key_name, struct passwd *pwd,     \
                                                char *buf, size_t buflen, int *errnop)                                 \
  {                                                                                                                    \
    int i;                                                                                                             \
    json_error_t error;                                                                                                \
    json_t *user;                                                                                                      \
    json_t *users = json_loads(data, 0, &error);                                                                       \
                                                                                                                       \
    if (users == NULL) {                                                                                               \
      syslog(LOG_ERR, "%s[L%d] json parse error: %s", __func__, __LINE__, error.text);                                 \
      goto leave;                                                                                                      \
    }                                                                                                                  \
                                                                                                                       \
    json_array_foreach(users, i, user)                                                                                 \
    {                                                                                                                  \
      if (!json_is_object(user)) {                                                                                     \
        continue;                                                                                                      \
      }                                                                                                                \
      key_type current = json_##json_type##_value(json_object_get(user, #json_key));                                   \
                                                                                                                       \
      if (match_method) {                                                                                              \
        PASSWD_ENSURE(user)                                                                                            \
        free(data);                                                                                                    \
        json_decref(users);                                                                                            \
        return NSS_STATUS_SUCCESS;                                                                                     \
      }                                                                                                                \
    }                                                                                                                  \
                                                                                                                       \
    free(data);                                                                                                        \
    json_decref(users);                                                                                                \
    return NSS_STATUS_NOTFOUND;                                                                                        \
  leave:                                                                                                               \
    return NSS_STATUS_UNAVAIL;                                                                                         \
  }

PASSWD_ENSURE_BY(name, const char *, user_name, string, name, (strcmp(current, user_name) == 0))
PASSWD_ENSURE_BY(uid, uid_t, uid, integer, id, current == uid)

PASSWD_GET_SINGLE(getpwnam_r, const char *name, "users?name=%s", name)
PASSWD_GET_SINGLE(getpwuid_r, uid_t uid, "users?id=%d", uid)

enum nss_status inner_nss_stns_setpwent(char *data, stns_conf_t *c)
{
  pthread_mutex_lock(&pwent_mutex);
  json_error_t error;

  entries = json_loads(data, 0, &error);

  if (entries == NULL) {
    syslog(LOG_ERR, "%s[L%d] json parse error: %s", __func__, __LINE__, error.text);
    free(data);
    pthread_mutex_unlock(&pwent_mutex);
    return NSS_STATUS_UNAVAIL;
  }
  entry_idx = 0;

  pthread_mutex_unlock(&pwent_mutex);
  return NSS_STATUS_SUCCESS;
}

enum nss_status _nss_stns_setpwent(void)
{
  int curl_result;
  stns_http_response_t r;
  json_error_t error;
  stns_conf_t c;
  stns_load_config(STNS_CONFIG_FILE, &c);

  curl_result = stns_request(&c, "users", &r);
  if (curl_result != CURLE_OK) {
    return NSS_STATUS_UNAVAIL;
  }

  return inner_nss_stns_setpwent(r.data, &c);
}

enum nss_status _nss_stns_endpwent(void)
{
  pthread_mutex_lock(&pwent_mutex);
  json_decref(entries);
  entry_idx = 0;
  pthread_mutex_unlock(&pwent_mutex);
  return NSS_STATUS_SUCCESS;
}

enum nss_status inner_nss_stns_getpwent_r(stns_conf_t *c, struct passwd *pwd, char *buf, size_t buflen, int *errnop)
{
  enum nss_status ret = NSS_STATUS_SUCCESS;
  pthread_mutex_lock(&pwent_mutex);

  if (entries == NULL) {
    ret = _nss_stns_setpwent();
  }

  if (ret != NSS_STATUS_SUCCESS) {
    pthread_mutex_unlock(&pwent_mutex);
    return ret;
  }

  if (entry_idx >= json_array_size(entries)) {
    *errnop = ENOENT;
    pthread_mutex_unlock(&pwent_mutex);
    return NSS_STATUS_NOTFOUND;
  }

  json_t *user = json_array_get(entries, entry_idx);

  PASSWD_ENSURE(user);
  entry_idx++;
  pthread_mutex_unlock(&pwent_mutex);
  return NSS_STATUS_SUCCESS;
}

enum nss_status _nss_stns_getpwent_r(struct passwd *pwd, char *buf, size_t buflen, int *errnop)
{
  enum nss_status ret = NSS_STATUS_SUCCESS;
  stns_conf_t c;
  stns_load_config(STNS_CONFIG_FILE, &c);
  return inner_nss_stns_getpwent_r(&c, pwd, buf, buflen, errnop);
}
