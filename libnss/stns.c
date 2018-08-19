#include "stns.h"
#include "toml.h"
#include <errno.h>
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>

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

// base https://github.com/linyows/octopass/blob/master/octopass.c
static size_t response_callback(void *buffer, size_t size, size_t nmemb, void *userp)
{
  size_t segsize            = size * nmemb;
  stns_http_response_t *res = (stns_http_response_t *)userp;

  if (segsize > STNS_MAX_BUFFER_SIZE) {
    fprintf(stderr, "Response is too large\n");
    return 0;
  }

  res->data = realloc(res->data, res->size + segsize + 1);

  if (!res->data) {
    res->data = (char *)malloc(segsize);
  } else {
    res->data = (char *)realloc(res->data, res->size + segsize);
  }

  memcpy(&(res->data[res->size]), buffer, (size_t)segsize);
  res->size += segsize;
  res->data[res->size] = '\0';

  return segsize;
}

// base https://github.com/linyows/octopass/blob/master/octopass.c
static CURLcode _stns_request(stns_conf_t *c, char *path, stns_http_response_t *res)
{
  char *auth;
  char *url;
  CURL *hnd;
  CURLcode result;
  struct curl_slist *headers = NULL;

  if (c->auth_token != NULL) {
    auth = (char *)malloc(strlen(c->auth_token) + 22);
    sprintf(auth, "Authorization: token %s", c->auth_token);
  } else {
    auth = NULL;
  }

  url = (char *)malloc(strlen(c->api_endpoint) + strlen(path) + 2);
  sprintf(url, "%s/%s", c->api_endpoint, path);

  res->data        = NULL;
  res->size        = 0;
  res->status_code = (long *)0;

  if (auth != NULL) {
    headers = curl_slist_append(headers, auth);
  }

#ifdef DEBUG
  printf("send http request: %s\n", url);
#endif

  hnd = curl_easy_init();
  curl_easy_setopt(hnd, CURLOPT_URL, url);
  curl_easy_setopt(hnd, CURLOPT_NOPROGRESS, 1);
  curl_easy_setopt(hnd, CURLOPT_USERAGENT, STNS_VERSION_WITH_NAME);
  curl_easy_setopt(hnd, CURLOPT_HTTPHEADER, headers);
  curl_easy_setopt(hnd, CURLOPT_SSL_VERIFYPEER, c->ssl_verify);
  curl_easy_setopt(hnd, CURLOPT_TIMEOUT, c->request_timeout);
  curl_easy_setopt(hnd, CURLOPT_WRITEFUNCTION, response_callback);
  curl_easy_setopt(hnd, CURLOPT_WRITEDATA, res);
  curl_easy_setopt(hnd, CURLOPT_NOSIGNAL, 1);
  curl_easy_setopt(hnd, CURLOPT_FAILONERROR, 1);

  result = curl_easy_perform(hnd);

  if (result != CURLE_OK) {
    fprintf(stderr, "http request failed: %s\n", curl_easy_strerror(result));
  } else {
    long *code;
    curl_easy_getinfo(hnd, CURLINFO_RESPONSE_CODE, &code);
    res->status_code = code;
  }

  free(auth);
  free(url);
  curl_easy_cleanup(hnd);
  curl_slist_free_all(headers);
  return result;
}

void stns_request(stns_conf_t *c, char *path, stns_http_response_t *res)
{
  CURLcode result;
  int retry_count = c->request_retry;

  result = _stns_request(c, path, res);
  while (1) {
    if (result != CURLE_OK && retry_count > 0) {
      sleep(1);
      result = _stns_request(c, path, res);
      retry_count--;
    } else {
      break;
    }
  }
}
