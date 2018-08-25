#include "stns.h"
#include "toml.h"
#include <errno.h>
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>

pthread_mutex_t user_mutex  = PTHREAD_MUTEX_INITIALIZER;
pthread_mutex_t group_mutex = PTHREAD_MUTEX_INITIALIZER;
int highest_user_id         = 0;
int lowest_user_id          = 0;
int highest_group_id        = 0;
int lowest_group_id         = 0;

SET_GET_HIGH_LOW_ID(highest, user);
SET_GET_HIGH_LOW_ID(lowest, user);
SET_GET_HIGH_LOW_ID(highest, group);
SET_GET_HIGH_LOW_ID(lowest, group);

ID_QUERY_AVAILABLE(user, high, <)
ID_QUERY_AVAILABLE(user, low, >)
ID_QUERY_AVAILABLE(group, high, <)
ID_QUERY_AVAILABLE(group, low, >)

void stns_load_config(char *filename, stns_conf_t *c)
{
  char errbuf[200];
  const char *raw;

  FILE *fp = fopen(filename, "r");
  if (!fp) {
    syslog(LOG_ERR, "%s[L%d] cannot open %s: %s", __func__, __LINE__, filename, strerror(errno));

    exit(1);
  }

  toml_table_t *tab = toml_parse_file(fp, errbuf, sizeof(errbuf));

  if (!tab) {
    syslog(LOG_ERR, "%s[L%d] %s", __func__, __LINE__, errbuf);
    exit(1);
  }

  GET_TOML_BYKEY(api_endpoint, toml_rtos, "http://localhost:1104/v1");
  GET_TOML_BYKEY(auth_token, toml_rtos, NULL);
  GET_TOML_BYKEY(user, toml_rtos, NULL);
  GET_TOML_BYKEY(password, toml_rtos, NULL);
  GET_TOML_BYKEY(query_wrapper, toml_rtos, NULL);
  GET_TOML_BYKEY(chain_ssh_wrapper, toml_rtos, NULL);
  GET_TOML_BYKEY(http_proxy, toml_rtos, NULL);

  GET_TOML_BYKEY(uid_shift, toml_rtoi, 0);
  GET_TOML_BYKEY(gid_shift, toml_rtoi, 0);
  GET_TOML_BYKEY(ssl_verify, toml_rtob, 1);
  GET_TOML_BYKEY(request_timeout, toml_rtoi, 10);
  GET_TOML_BYKEY(request_retry, toml_rtoi, 3);
  GET_TOML_BYKEY(request_locktime, toml_rtoi, 60);

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

static void trim(char *s)
{
  int i, j;

  for (i = strlen(s) - 1; i >= 0 && isspace(s[i]); i--)
    ;
  s[i + 1] = '\0';
  for (i = 0; isspace(s[i]); i++)
    ;
  if (i > 0) {
    j = 0;
    while (s[i])
      s[j++] = s[i++];
    s[j] = '\0';
  }
}

#define SET_TRIM_ID(high_or_low, user_or_group)                                                                        \
  tp = strtok(NULL, ".");                                                                                              \
  trim(tp);                                                                                                            \
  set_##high_or_low##est_##user_or_group##_id(atoi(tp));

static size_t header_callback(char *buffer, size_t size, size_t nitems, void *userdata)
{
  char *tp;
  tp = strtok(buffer, ":");
  if (strcmp(tp, "User-Highest-Id") == 0) {
    SET_TRIM_ID(high, user)
  } else if (strcmp(tp, "User-Lowest-Id") == 0) {
    SET_TRIM_ID(low, user)
  } else if (strcmp(tp, "Group-Highest-Id") == 0) {
    SET_TRIM_ID(high, group)
  } else if (strcmp(tp, "Group-Lowest-Id") == 0) {
    SET_TRIM_ID(low, group)
  }

  return nitems * size;
}

// base https://github.com/linyows/octopass/blob/master/octopass.c
// size is always 1
static size_t response_callback(void *buffer, size_t size, size_t nmemb, void *userp)
{
  size_t segsize            = size * nmemb;
  stns_http_response_t *res = (stns_http_response_t *)userp;

  if (segsize > STNS_MAX_BUFFER_SIZE) {
    syslog(LOG_ERR, "%s[L%d] Response is too large", __func__, __LINE__);
    return 0;
  }

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

static int _stns_wrapper_request(stns_conf_t *c, char *path, stns_http_response_t *res)
{
  int rsize    = 0;
  char *result = malloc(1);
  res->data    = NULL;
  res->size    = 0;

  if (stns_exec_cmd(c->query_wrapper, path, result)) {
    rsize = strlen(result);

    if (res->data) {
      res->data = (char *)realloc(res->data, rsize + 1);
    } else {
      res->data = (char *)malloc(rsize + 1);
    }

    memcpy(&(res->data[0]), result, (size_t)rsize);
    res->data[rsize] = '\0';
  } else {
    free(result);
    return 0;
  }

  free(result);
  return 1;
}

// base https://github.com/linyows/octopass/blob/master/octopass.c
static CURLcode _stns_http_request(stns_conf_t *c, char *path, stns_http_response_t *res)
{
  char *auth;
  char *url;
  CURL *curl;
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
  syslog(LOG_DEBUG, "%s[L%d] send http request: %s", __func__, __LINE__, url);
#endif

  curl = curl_easy_init();
  curl_easy_setopt(curl, CURLOPT_URL, url);
  curl_easy_setopt(curl, CURLOPT_NOPROGRESS, 1);
  curl_easy_setopt(curl, CURLOPT_USERAGENT, STNS_VERSION_WITH_NAME);
  curl_easy_setopt(curl, CURLOPT_HTTPHEADER, headers);
  curl_easy_setopt(curl, CURLOPT_SSL_VERIFYPEER, c->ssl_verify);
  curl_easy_setopt(curl, CURLOPT_TIMEOUT, c->request_timeout);
  curl_easy_setopt(curl, CURLOPT_WRITEFUNCTION, response_callback);
  curl_easy_setopt(curl, CURLOPT_HEADERFUNCTION, header_callback);
  curl_easy_setopt(curl, CURLOPT_WRITEDATA, res);
  curl_easy_setopt(curl, CURLOPT_NOSIGNAL, 1);
  curl_easy_setopt(curl, CURLOPT_FAILONERROR, 1);

  if (c->user != NULL) {
    curl_easy_setopt(curl, CURLOPT_USERNAME, c->user);
  }

  if (c->password != NULL) {
    curl_easy_setopt(curl, CURLOPT_PASSWORD, c->password);
  }

  result = curl_easy_perform(curl);

  if (result != CURLE_OK) {
    syslog(LOG_ERR, "%s[L%d] http request failed: %s", __func__, __LINE__, curl_easy_strerror(result));
  } else {
    long *code;
    curl_easy_getinfo(curl, CURLINFO_RESPONSE_CODE, &code);
    res->status_code = code;
  }

  free(auth);
  free(url);
  curl_easy_cleanup(curl);
  curl_slist_free_all(headers);
  return result;
}

int stns_request_available(char *path, stns_conf_t *c)
{
  struct stat st;
  if (stat(path, &st) != 0) {
    return 1;
  }

  unsigned long now  = time(NULL);
  unsigned long diff = now - st.st_ctime;
  if (diff > c->request_locktime) {
    remove(path);
    return 1;
  }
  return 0;
}

void stns_make_lockfile(char *path)
{
  FILE *fp;
  fp = fopen(path, "w");
  if (fp) {
    fclose(fp);
  }
}

int stns_request(stns_conf_t *c, char *path, stns_http_response_t *res)
{
  CURLcode result;
  int retry_count = c->request_retry;

  if (!stns_request_available(STNS_LOCK_FILE, c))
    return CURLE_COULDNT_CONNECT;

  if (c->query_wrapper == NULL) {
    result = _stns_http_request(c, path, res);
    while (1) {
      if (result != CURLE_OK && retry_count > 0) {
        if (result == CURLE_HTTP_RETURNED_ERROR) {
          break;
        }
        sleep(1);
        result = _stns_http_request(c, path, res);
        retry_count--;
      } else {
        break;
      }
    }
  } else {
    result = _stns_wrapper_request(c, path, res);
  }

  if (result == CURLE_COULDNT_CONNECT) {
    stns_make_lockfile(STNS_LOCK_FILE);
  }

  return result;
}

int stns_exec_cmd(char *cmd, char *arg, char *result)
{
  FILE *fp;
  char *c;

  if (arg != NULL) {
    c = malloc(strlen(cmd) + strlen(arg) + 2);
    sprintf(c, "%s %s", cmd, arg);
  } else {
    c = cmd;
  }

  if ((fp = popen(c, "r")) == NULL) {
    return 0;
  }

  char buf[MAXBUF];
  int total_len = 0;
  int len       = 0;

  while (fgets(buf, sizeof(buf), fp) != NULL) {
    len = strlen(buf);
    if (result) {
      result = (char *)realloc(result, total_len + len + 1);
    } else {
      return 0;
    }
    strcpy(result + total_len, buf);
    total_len += len;
  }
  result[total_len] = '\0';
  pclose(fp);

  if (total_len == 0) {
    return 0;
  }
  return 1;
}
