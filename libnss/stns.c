#include "stns.h"
#include "toml.h"
#include <errno.h>
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <dirent.h>

pthread_mutex_t user_mutex   = PTHREAD_MUTEX_INITIALIZER;
pthread_mutex_t group_mutex  = PTHREAD_MUTEX_INITIALIZER;
pthread_mutex_t delete_mutex = PTHREAD_MUTEX_INITIALIZER;
int highest_user_id          = 0;
int lowest_user_id           = 0;
int highest_group_id         = 0;
int lowest_group_id          = 0;

SET_GET_HIGH_LOW_ID(highest, user);
SET_GET_HIGH_LOW_ID(lowest, user);
SET_GET_HIGH_LOW_ID(highest, group);
SET_GET_HIGH_LOW_ID(lowest, group);

ID_QUERY_AVAILABLE(user, high, <)
ID_QUERY_AVAILABLE(user, low, >)
ID_QUERY_AVAILABLE(group, high, <)
ID_QUERY_AVAILABLE(group, low, >)

#define TRIM_SLASH(key)                                                                                                \
  if (c->key != NULL) {                                                                                                \
    const int key##_len = strlen(c->key);                                                                              \
    if (key##_len > 0) {                                                                                               \
      if (c->key[key##_len - 1] == '/') {                                                                              \
        c->key = strndup(c->key, key##_len - 1);                                                                       \
      }                                                                                                                \
    }                                                                                                                  \
  }

static void stns_force_create_cache_dir(stns_conf_t *c)
{
  if (c->cache && geteuid() == 0) {
    struct stat statBuf;

    char path[MAXBUF];
    sprintf(path, "%s", c->cache_dir);
    if (stat(path, &statBuf) != 0) {
      mode_t um = {0};
      um        = umask(0);
      mkdir(path, S_IRUSR | S_IWUSR | S_IXUSR | S_IRGRP | S_IWGRP | S_IXGRP | S_IROTH | S_IWOTH | S_IXOTH);
      umask(um);
    }
  }
}

void stns_load_config(char *filename, stns_conf_t *c)
{
  char errbuf[200];
  const char *raw;

  FILE *fp = fopen(filename, "r");
  if (!fp) {
    syslog(LOG_ERR, "%s(stns)[L%d] cannot open %s: %s", __func__, __LINE__, filename, strerror(errno));

    exit(1);
  }

  toml_table_t *tab = toml_parse_file(fp, errbuf, sizeof(errbuf));

  if (!tab) {
    syslog(LOG_ERR, "%s(stns)[L%d] %s", __func__, __LINE__, errbuf);
    exit(1);
  }

  GET_TOML_BYKEY(api_endpoint, toml_rtos, "http://localhost:1104/v1", TOML_STR);
  GET_TOML_BYKEY(cache_dir, toml_rtos, "/var/cache/stns", TOML_STR);
  GET_TOML_BYKEY(auth_token, toml_rtos, NULL, TOML_NULL_OR_INT);
  GET_TOML_BYKEY(user, toml_rtos, NULL, TOML_NULL_OR_INT);
  GET_TOML_BYKEY(password, toml_rtos, NULL, TOML_NULL_OR_INT);
  GET_TOML_BYKEY(query_wrapper, toml_rtos, NULL, TOML_NULL_OR_INT);
  GET_TOML_BYKEY(chain_ssh_wrapper, toml_rtos, NULL, TOML_NULL_OR_INT);
  GET_TOML_BYKEY(http_proxy, toml_rtos, NULL, TOML_NULL_OR_INT);

  GET_TOML_BYKEY(uid_shift, toml_rtoi, 0, TOML_NULL_OR_INT);
  GET_TOML_BYKEY(gid_shift, toml_rtoi, 0, TOML_NULL_OR_INT);
  GET_TOML_BYKEY(cache_ttl, toml_rtoi, 600, TOML_NULL_OR_INT);
  GET_TOML_BYKEY(negative_cache_ttl, toml_rtoi, 60, TOML_NULL_OR_INT);
  GET_TOML_BYKEY(ssl_verify, toml_rtob, 1, TOML_NULL_OR_INT);
  GET_TOML_BYKEY(cache, toml_rtob, 1, TOML_NULL_OR_INT);
  GET_TOML_BYKEY(request_timeout, toml_rtoi, 10, TOML_NULL_OR_INT);
  GET_TOML_BYKEY(request_retry, toml_rtoi, 3, TOML_NULL_OR_INT);
  GET_TOML_BYKEY(request_locktime, toml_rtoi, 60, TOML_NULL_OR_INT);

  TRIM_SLASH(api_endpoint)
  TRIM_SLASH(cache_dir)

  stns_force_create_cache_dir(c);
  fclose(fp);
  toml_free(tab);
}

void stns_unload_config(stns_conf_t *c)
{
  UNLOAD_TOML_BYKEY(api_endpoint);
  UNLOAD_TOML_BYKEY(cache_dir);
  UNLOAD_TOML_BYKEY(auth_token);
  UNLOAD_TOML_BYKEY(user);
  UNLOAD_TOML_BYKEY(password);
  UNLOAD_TOML_BYKEY(query_wrapper);
  UNLOAD_TOML_BYKEY(chain_ssh_wrapper);
  UNLOAD_TOML_BYKEY(http_proxy);
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

#define SET_TRIM_ID(high_or_low, user_or_group, short_name)                                                            \
  tp = strtok(NULL, ".");                                                                                              \
  trim(tp);                                                                                                            \
  set_##user_or_group##_##high_or_low##est_id(atoi(tp) + c->short_name##id_shift);

static size_t header_callback(char *buffer, size_t size, size_t nitems, void *userdata)
{

  stns_conf_t *c = (stns_conf_t *)userdata;
  char *tp;
  tp = strtok(buffer, ":");
  if (strcmp(tp, "User-Highest-Id") == 0) {
    SET_TRIM_ID(high, user, u)
  } else if (strcmp(tp, "User-Lowest-Id") == 0) {
    SET_TRIM_ID(low, user, u)
  } else if (strcmp(tp, "Group-Highest-Id") == 0) {
    SET_TRIM_ID(high, group, g)
  } else if (strcmp(tp, "Group-Lowest-Id") == 0) {
    SET_TRIM_ID(low, group, g)
  }

  return nitems * size;
}

// base https://github.com/linyows/octopass/blob/master/octopass.c
// size is always 1
static size_t response_callback(void *buffer, size_t size, size_t nmemb, void *userp)
{
  size_t segsize       = size * nmemb;
  stns_response_t *res = (stns_response_t *)userp;

  if (segsize > STNS_MAX_BUFFER_SIZE) {
    syslog(LOG_ERR, "%s(stns)[L%d] Response is too large", __func__, __LINE__);
    return 0;
  }

  if (!res->data) {
    res->data = (char *)malloc(segsize + 1);
  } else {
    res->data = (char *)realloc(res->data, res->size + segsize + 1);
  }

  memcpy(&(res->data[res->size]), buffer, segsize);
  res->size += segsize;
  res->data[res->size] = 0;

  return segsize;
}

// base https://github.com/linyows/octopass/blob/master/octopass.c
static CURLcode inner_http_request(stns_conf_t *c, char *path, stns_response_t *res)
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
  res->status_code = (long)200;

  if (auth != NULL) {
    headers = curl_slist_append(headers, auth);
  }

#ifdef DEBUG
  syslog(LOG_DEBUG, "%s(stns)[L%d] send http request: %s", __func__, __LINE__, url);
#endif

  curl = curl_easy_init();
  curl_easy_setopt(curl, CURLOPT_URL, url);
  curl_easy_setopt(curl, CURLOPT_NOPROGRESS, 1);
  curl_easy_setopt(curl, CURLOPT_USERAGENT, STNS_VERSION_WITH_NAME);
  curl_easy_setopt(curl, CURLOPT_HTTPHEADER, headers);
  curl_easy_setopt(curl, CURLOPT_SSL_VERIFYPEER, c->ssl_verify);
  curl_easy_setopt(curl, CURLOPT_SSL_VERIFYHOST, c->ssl_verify);
  curl_easy_setopt(curl, CURLOPT_TIMEOUT, c->request_timeout);
  curl_easy_setopt(curl, CURLOPT_WRITEFUNCTION, response_callback);
  curl_easy_setopt(curl, CURLOPT_HEADERFUNCTION, header_callback);
  curl_easy_setopt(curl, CURLOPT_WRITEDATA, res);
  curl_easy_setopt(curl, CURLOPT_HEADERDATA, c);
  curl_easy_setopt(curl, CURLOPT_NOSIGNAL, 1);
  curl_easy_setopt(curl, CURLOPT_FAILONERROR, 1);

  if (c->user != NULL) {
    curl_easy_setopt(curl, CURLOPT_USERNAME, c->user);
  }

  if (c->password != NULL) {
    curl_easy_setopt(curl, CURLOPT_PASSWORD, c->password);
  }

  if (c->http_proxy != NULL) {
    curl_easy_setopt(curl, CURLOPT_PROXY, c->http_proxy);
  }

  result = curl_easy_perform(curl);

  if (result != CURLE_OK) {
    syslog(LOG_ERR, "%s(stns)[L%d] http request failed: %s", __func__, __LINE__, curl_easy_strerror(result));
    if (result == CURLE_HTTP_RETURNED_ERROR) {
      long code;
      curl_easy_getinfo(curl, CURLINFO_RESPONSE_CODE, &code);
      res->status_code = code;
    }
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

// base: https://github.com/linyows/octopass/blob/master/octopass.c
void stns_export_file(char *file, char *data)
{
  struct stat statbuf;
  if (stat(file, &statbuf) != -1 && statbuf.st_uid != geteuid()) {
    return;
  }

  FILE *fp = fopen(file, "w");
  if (!fp) {
    syslog(LOG_ERR, "%s(stns)[L%d] cannot open %s", __func__, __LINE__, file);
    return;
  }
  if (data != NULL) {
    fprintf(fp, "%s", data);
  }
  fclose(fp);

  mode_t um = {0};
  um        = umask(0);
  chmod(file, S_IRUSR | S_IWUSR | S_IXUSR | S_IRGRP | S_IROTH);
  umask(um);
}

// base: https://github.com/linyows/octopass/blob/master/octopass.c
int stns_import_file(char *file, stns_response_t *res)
{
  FILE *fp = fopen(file, "r");
  if (!fp) {
    syslog(LOG_ERR, "%s(stns)[L%d] cannot open %s", __func__, __LINE__, file);
    return 0;
  }

  char buf[MAXBUF];
  int total_len = 0;
  int len       = 0;

  while (fgets(buf, sizeof(buf), fp) != NULL) {
    len = strlen(buf);
    if (!res->data) {
      res->data = (char *)malloc(len + 1);
    } else {
      res->data = (char *)realloc(res->data, total_len + len + 1);
    }
    strcpy(res->data + total_len, buf);
    total_len += len;
  }
  fclose(fp);

  return 1;
}

static void *delete_cache_files(void *data)
{
  stns_conf_t *c = (stns_conf_t *)data;
  DIR *dp;
  struct dirent *ent;
  struct stat statbuf;
  unsigned long now = time(NULL);
  char dir[MAXBUF];
  sprintf(dir, "%s/%d", c->cache_dir, geteuid());

  pthread_mutex_lock(&delete_mutex);
  if ((dp = opendir(dir)) == NULL) {
    syslog(LOG_ERR, "%s(stns)[L%d] cannot open %s: %s", __func__, __LINE__, dir, strerror(errno));
    pthread_mutex_unlock(&delete_mutex);
    return NULL;
  }

  char *buf = malloc(1);
  while ((ent = readdir(dp)) != NULL) {
    buf = (char *)realloc(buf, strlen(dir) + strlen(ent->d_name) + 2);
    sprintf(buf, "%s/%s", dir, ent->d_name);

    if (stat(buf, &statbuf) == 0 && (statbuf.st_uid == geteuid() || geteuid() == 0)) {
      unsigned long diff = now - statbuf.st_mtime;

      if (!S_ISDIR(statbuf.st_mode) &&
          ((diff > c->cache_ttl && statbuf.st_size > 0) || (diff > c->negative_cache_ttl && statbuf.st_size == 0))) {

        if (unlink(buf) == -1) {
          syslog(LOG_ERR, "%s(stns)[L%d] cannot delete %s: %s", __func__, __LINE__, buf, strerror(errno));
        }
      }
    }
  }
  free(buf);
  closedir(dp);
  pthread_mutex_unlock(&delete_mutex);
  return NULL;
}

int stns_request(stns_conf_t *c, char *path, stns_response_t *res)
{
  CURLcode result;
  pthread_t pthread;
  int retry_count  = c->request_retry;
  res->data        = NULL;
  res->size        = 0;
  res->status_code = (long)200;

  if (path == NULL) {
    return CURLE_HTTP_RETURNED_ERROR;
  }

  char *base = curl_escape(path, strlen(path));
  char fpath[MAXBUF];
  sprintf(fpath, "%s/%d/%s", c->cache_dir, geteuid(), base);
  free(base);

  if (c->cache) {
    FILE *fp = fopen(fpath, "r");
    if (fp != NULL) {
      fclose(fp);
      struct stat statbuf;
      if (stat(fpath, &statbuf) != -1 && statbuf.st_uid == geteuid()) {
        unsigned long now  = time(NULL);
        unsigned long diff = now - statbuf.st_mtime;

        // resource notfound
        if ((diff < c->cache_ttl && statbuf.st_size > 0) || (diff < c->negative_cache_ttl && statbuf.st_size == 0)) {
          if (statbuf.st_size == 0) {
            res->status_code = STNS_HTTP_NOTFOUND;
            return CURLE_HTTP_RETURNED_ERROR;
          }

          pthread_mutex_lock(&delete_mutex);
          if (!stns_import_file(fpath, res)) {
            pthread_mutex_unlock(&delete_mutex);
            goto request;
          }
          pthread_mutex_unlock(&delete_mutex);
          res->size = strlen(res->data);
          return CURLE_OK;
        }
      }
    }
  }
request:
  if (!stns_request_available(STNS_LOCK_FILE, c))
    return CURLE_COULDNT_CONNECT;

  if (c->cache) {
    pthread_create(&pthread, NULL, &delete_cache_files, (void *)c);
  }

  if (c->query_wrapper == NULL) {
    result = inner_http_request(c, path, res);
    while (1) {
      if (result != CURLE_OK && retry_count > 0) {
        if (result == CURLE_HTTP_RETURNED_ERROR) {
          break;
        }
        sleep(1);
        result = inner_http_request(c, path, res);
        retry_count--;
      } else {
        break;
      }
    }
  } else {
    result = stns_exec_cmd(c->query_wrapper, path, res);
  }

  if (result == CURLE_COULDNT_CONNECT) {
    stns_make_lockfile(STNS_LOCK_FILE);
  }

  if (c->cache) {
    pthread_join(pthread, NULL);
    pthread_mutex_lock(&delete_mutex);
    stns_export_file(fpath, res->data);
    pthread_mutex_unlock(&delete_mutex);
  }
  return result;
}

unsigned int match(char *pattern, char *text)
{
  regex_t regex;
  int rc;

  if (text == NULL) {
    return 0;
  }

  rc = regcomp(&regex, pattern, REG_EXTENDED | REG_NOSUB);
  if (rc == 0)
    rc = regexec(&regex, text, 0, 0, 0);
  regfree(&regex);
  return rc == 0;
}

int stns_exec_cmd(char *cmd, char *arg, stns_response_t *r)
{
  FILE *fp;
  char *c;

  r->data        = NULL;
  r->size        = 0;
  r->status_code = (long)200;

  if (!match("^[a-z0-9_.]+$", arg)) {
    return 0;
  }

  if (arg != NULL) {
    c = malloc(strlen(cmd) + strlen(arg) + 2);
    sprintf(c, "%s %s", cmd, arg);
  } else {
    c = cmd;
  }

  if ((fp = popen(c, "r")) == NULL) {
    goto err;
  }

  char buf[MAXBUF];
  int total_len = 0;
  int len       = 0;

  while (fgets(buf, sizeof(buf), fp) != NULL) {
    len = strlen(buf);
    if (r->data) {
      r->data = (char *)realloc(r->data, total_len + len + 1);
    } else {
      r->data = (char *)malloc(total_len + len + 1);
    }
    strcpy(r->data + total_len, buf);
    total_len += len;
  }
  pclose(fp);

  if (total_len == 0) {
    goto err;
  }
  if (arg != NULL)
    free(c);

  r->size = total_len;
  return 1;
err:
  if (arg != NULL)
    free(c);
  return 0;
}
