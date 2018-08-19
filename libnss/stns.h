#ifndef STNS_H
#define STNS_H

#include <errno.h>
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include "toml.h"

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
#endif /* STNS_H */
