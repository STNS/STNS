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

#define STNS_VERSION "2.0.0"
#define STNS_VERSION_WITH_NAME "stns/" STNS_VERSION
// 10MB
#define STNS_MAX_BUFFER_SIZE (10 * 1024 * 1024)

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
extern void stns_request(stns_conf_t *, char *, stns_http_response_t *);
#endif /* STNS_H */
