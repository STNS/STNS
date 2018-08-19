#include "stns.h"
#include "toml.h"
#include <errno.h>
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#define GET_TOML_BYKEY(m, method, empty)                                                                               \
  if (0 != (raw = toml_raw_in(tab, #m))) {                                                                             \
    if (0 != method(raw, &c->m)) {                                                                                     \
      fprintf(stderr, "ERROR: cannot parse toml file %s key %s\n", filename, #m);                                      \
    }                                                                                                                  \
  } else {                                                                                                             \
    c->m = empty;                                                                                                      \
  }

void stns_load_config(char *filename, stns_conf_t *c)
{
  char errbuf[200];
  const char *key;
  const char *raw;
  toml_array_t *arr;

  FILE *fp = fopen(filename, "r");
  if (!fp) {
    fprintf(stderr, "ERROR: cannot open %s: %s\n", filename, strerror(errno));
    exit(1);
  }

  toml_table_t *tab = toml_parse_file(fp, errbuf, sizeof(errbuf));

  if (!tab) {
    fprintf(stderr, "ERROR: %s\n", errbuf);
    return;
  }

  GET_TOML_BYKEY(api_endpoint, toml_rtos, NULL);
  GET_TOML_BYKEY(auth_token, toml_rtos, NULL);
  GET_TOML_BYKEY(user, toml_rtos, NULL);
  GET_TOML_BYKEY(password, toml_rtos, NULL);
  GET_TOML_BYKEY(chain_ssh_wrapper, toml_rtos, NULL);
  GET_TOML_BYKEY(http_proxy, toml_rtos, NULL);

  GET_TOML_BYKEY(uid_shift, toml_rtoi, 0);
  GET_TOML_BYKEY(gid_shift, toml_rtoi, 0);
  GET_TOML_BYKEY(ssl_verify, toml_rtob, 1);
  GET_TOML_BYKEY(request_timeout, toml_rtoi, 10);
  GET_TOML_BYKEY(request_retry, toml_rtoi, 3);

  // 末尾の/を取り除く
  const int len = strlen(c->api_endpoint);
  if (len > 0) {
    if (c->api_endpoint[len - 1] == '/') {
      c->api_endpoint[len - 1] = '\0';
    }
  }

  fclose(fp);
  toml_free(tab);
}
