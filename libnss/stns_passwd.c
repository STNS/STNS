#include "stns.h"

#define PASSWD_DEFAULT(buf, name, def)                                                                                 \
  char buf[MAXBUF];                                                                                                    \
  if (strlen(name) > 0) {                                                                                              \
    strcpy(buf, name);                                                                                                 \
  } else {                                                                                                             \
    strcpy(buf, def);                                                                                                  \
  }

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
    return ensure_passwd(r.data, &c, value, pwd, buf, buflen, errnop);                                                 \
  }

static enum nss_status ensure_passwd(char *data, stns_conf_t *c, const char *name, struct passwd *pwd, char *buf,
                                     size_t buflen, int *errnop)
{
  int i;
  json_error_t error;
  json_t *user;
  json_t *users = json_loads(data, 0, &error);

  if (users == NULL) {
    syslog(LOG_ERR, "%s[L%d] json parse error: %s", __func__, __LINE__, error.text);
    goto leave;
  }

  json_array_foreach(users, i, user)
  {
    if (!json_is_object(user)) {
      continue;
    }
    const char *user_name = json_string_value(json_object_get(user, "name"));

    if (name != NULL && strcmp(user_name, name) == 0) {
      const json_int_t id       = json_integer_value(json_object_get(user, "id"));
      const json_int_t group_id = json_integer_value(json_object_get(user, "group_id"));
      const char *gecos         = json_string_value(json_object_get(user, "gecos"));
      const char *home_dir      = json_string_value(json_object_get(user, "directory"));
      const char *shell         = json_string_value(json_object_get(user, "shell"));

      PASSWD_DEFAULT(sh, shell, "/bin/bash");

      char b[MAXBUF];
      sprintf(b, "/home/%s", name);
      PASSWD_DEFAULT(dir, home_dir, b);

      int name_length    = strlen(name) + 1;
      int pw_length      = strlen("x") + 1;
      int gecos_length   = strlen(gecos) + 1;
      int homedir_length = strlen(home_dir) + 1;
      int shell_length   = strlen(shell) + 1;

      int total_length = name_length + pw_length + gecos_length + shell_length + homedir_length;

      if (buflen < total_length) {
        *errnop = ERANGE;
        return NSS_STATUS_TRYAGAIN;
      }

      pwd->pw_uid = c->uid_shift + id;
      pwd->pw_gid = c->gid_shift + group_id;

      strcpy(buf, name);
      pwd->pw_name = buf;
      buf += name_length;
      strcpy(buf, "x");
      pwd->pw_passwd = buf;
      buf += pw_length;

      strcpy(buf, gecos);
      pwd->pw_name = buf;
      buf += gecos_length;

      strcpy(buf, home_dir);
      pwd->pw_dir = buf;
      buf += homedir_length;

      strcpy(buf, shell);
      pwd->pw_shell = buf;
      buf += shell_length;
      return NSS_STATUS_SUCCESS;
    }
  }

  free(data);
  json_decref(users);
  return NSS_STATUS_NOTFOUND;
leave:
  free(data);
  return NSS_STATUS_UNAVAIL;
}

PASSWD_GET_SINGLE(getpwnam_r, const char *name, "users?name=%s", name)
PASSWD_GET_SINGLE(getpwuid_r, uid_t uid, "users?id=%d", uid)
