#ifndef STNS_H
#define STNS_H
#define DEBUG 1

#include <curl/curl.h>
#include <errno.h>
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include "toml.h"
#include <syslog.h>
#include <nss.h>
#include <grp.h>
#include <pwd.h>
#include <shadow.h>
#include <jansson.h>
#include <pthread.h>

#define STNS_VERSION "2.0.0"
#define STNS_VERSION_WITH_NAME "stns/" STNS_VERSION
// 10MB
#define STNS_MAX_BUFFER_SIZE (10 * 1024 * 1024)
#define STNS_CONFIG_FILE "/etc/stns/client/stns.conf"
#define MAXBUF 1024

typedef struct stns_http_response_t stns_http_response_t;
struct stns_http_response_t {
  char *data;
  size_t size;
  long *status_code;
};

typedef struct stns_conf_t stns_conf_t;
struct stns_conf_t {
  char *api_endpoint;
  char *auth_token;
  char *user;
  char *password;
  char *chain_ssh_wrapper;
  char *http_proxy;
  int uid_shift;
  int gid_shift;
  int ssl_verify;
  int request_timeout;
  int request_retry;
};

extern void stns_load_config(char *, stns_conf_t *);
extern int stns_request(stns_conf_t *, char *, stns_http_response_t *);

#define STNS_SET_ENTRIES(type, ltype, resource, query)                                                                 \
  pthread_mutex_t type##ent_mutex = PTHREAD_MUTEX_INITIALIZER;                                                         \
  enum nss_status inner_nss_stns_set##type##ent(char *data, stns_conf_t *c)                                            \
  {                                                                                                                    \
    pthread_mutex_lock(&type##ent_mutex);                                                                              \
    json_error_t error;                                                                                                \
                                                                                                                       \
    entries = json_loads(data, 0, &error);                                                                             \
                                                                                                                       \
    if (entries == NULL) {                                                                                             \
      syslog(LOG_ERR, "%s[L%d] json parse error: %s", __func__, __LINE__, error.text);                                 \
      free(data);                                                                                                      \
      pthread_mutex_unlock(&type##ent_mutex);                                                                          \
      return NSS_STATUS_UNAVAIL;                                                                                       \
    }                                                                                                                  \
    entry_idx = 0;                                                                                                     \
                                                                                                                       \
    pthread_mutex_unlock(&type##ent_mutex);                                                                            \
    return NSS_STATUS_SUCCESS;                                                                                         \
  }                                                                                                                    \
                                                                                                                       \
  enum nss_status _nss_stns_set##type##ent(void)                                                                       \
  {                                                                                                                    \
    int curl_result;                                                                                                   \
    stns_http_response_t r;                                                                                            \
    stns_conf_t c;                                                                                                     \
    stns_load_config(STNS_CONFIG_FILE, &c);                                                                            \
                                                                                                                       \
    curl_result = stns_request(&c, #query, &r);                                                                        \
    if (curl_result != CURLE_OK) {                                                                                     \
      return NSS_STATUS_UNAVAIL;                                                                                       \
    }                                                                                                                  \
                                                                                                                       \
    return inner_nss_stns_set##type##ent(r.data, &c);                                                                  \
  }                                                                                                                    \
                                                                                                                       \
  enum nss_status _nss_stns_end##type##ent(void)                                                                       \
  {                                                                                                                    \
    pthread_mutex_lock(&type##ent_mutex);                                                                              \
    json_decref(entries);                                                                                              \
    entry_idx = 0;                                                                                                     \
    pthread_mutex_unlock(&type##ent_mutex);                                                                            \
    return NSS_STATUS_SUCCESS;                                                                                         \
  }                                                                                                                    \
                                                                                                                       \
  enum nss_status inner_nss_stns_get##type##ent_r(stns_conf_t *c, struct resource *rbuf, char *buf, size_t buflen,     \
                                                  int *errnop)                                                         \
  {                                                                                                                    \
    enum nss_status ret = NSS_STATUS_SUCCESS;                                                                          \
    pthread_mutex_lock(&type##ent_mutex);                                                                              \
                                                                                                                       \
    if (entries == NULL) {                                                                                             \
      ret = _nss_stns_set##type##ent();                                                                                \
    }                                                                                                                  \
                                                                                                                       \
    if (ret != NSS_STATUS_SUCCESS) {                                                                                   \
      pthread_mutex_unlock(&type##ent_mutex);                                                                          \
      return ret;                                                                                                      \
    }                                                                                                                  \
                                                                                                                       \
    if (entry_idx >= json_array_size(entries)) {                                                                       \
      *errnop = ENOENT;                                                                                                \
      pthread_mutex_unlock(&type##ent_mutex);                                                                          \
      return NSS_STATUS_NOTFOUND;                                                                                      \
    }                                                                                                                  \
                                                                                                                       \
    json_t *user = json_array_get(entries, entry_idx);                                                                 \
                                                                                                                       \
    ltype##_ENSURE(user);                                                                                              \
    entry_idx++;                                                                                                       \
    pthread_mutex_unlock(&type##ent_mutex);                                                                            \
    return NSS_STATUS_SUCCESS;                                                                                         \
  }                                                                                                                    \
                                                                                                                       \
  enum nss_status _nss_stns_get##type##ent_r(struct resource *rbuf, char *buf, size_t buflen, int *errnop)             \
  {                                                                                                                    \
    stns_conf_t c;                                                                                                     \
    stns_load_config(STNS_CONFIG_FILE, &c);                                                                            \
    return inner_nss_stns_get##type##ent_r(&c, rbuf, buf, buflen, errnop);                                             \
  }

#endif /* STNS_H */
