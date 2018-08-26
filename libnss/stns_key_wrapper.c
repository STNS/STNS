#include "stns.h"
#include <getopt.h>
int main(int argc, char *argv[])
{

  int curl_result;
  stns_http_response_t r;
  stns_conf_t c;
  char url[MAXBUF];
  const char *tmpkey;
  char *keys      = NULL;
  char *conf_path = NULL;
  int ret;

  while ((ret = getopt(argc, argv, "c:")) != -1) {
    if (ret == -1)
      break;
    switch (ret) {
    case 'c':
      conf_path = optarg;
      break;
    default:
      break;
    }
  }

  if (argc == 1 || argc <= optind) {
    fprintf(stderr, "User name is a required parameter\n");
    return -1;
  }

  if (conf_path == NULL)
    stns_load_config(STNS_CONFIG_FILE, &c);
  else
    stns_load_config(conf_path, &c);

  sprintf(url, "users?name=%s", argv[optind]);
  curl_result = stns_request(&c, url, &r);
  if (curl_result != CURLE_OK) {
    fprintf(stderr, "http request failed user: %s\n", argv[optind]);
    return -1;
  }

  int i, k;
  int size = 0;
  json_error_t error;
  json_t *leaf;
  json_t *key;
  int key_size;
  json_t *root = json_loads(r.data, 0, &error);

  if (root == NULL) {
    free(r.data);
    syslog(LOG_ERR, "%s(stns)[L%d] json parse error: %s", __func__, __LINE__, error.text);
    return -1;
  }

  json_array_foreach(root, i, leaf)
  {
    if (!json_is_object(leaf))
      continue;

    json_array_foreach(json_object_get(leaf, "keys"), k, key)
    {
      if (size != 0) {
        keys[size] = '\n';
        size++;
        keys[size] = '\0';
      }
      tmpkey   = json_string_value(key);
      key_size = strlen(tmpkey);

      if (keys) {
        keys = (char *)realloc(keys, key_size + strlen(keys) + 2);
      } else {
        keys = (char *)malloc(strlen(tmpkey) + 2);
      }

      memcpy(&(keys[size]), tmpkey, (size_t)key_size);
      size += key_size;
    }
  }

  if (keys) {
    keys[size] = '\n';
    size++;
    keys[size] = '\0';
  }

  if (c.chain_ssh_wrapper != NULL) {
    char *result = malloc(1);
    if (stns_exec_cmd(c.chain_ssh_wrapper, argv[optind], result)) {
      key_size = strlen(result);
      keys     = (char *)realloc(keys, key_size + strlen(keys) + 1);
      memcpy(&(keys[size]), result, (size_t)key_size);
      size += key_size;
      keys[size] = '\0';
    }
  }

  fprintf(stdout, "%s\n", keys);
  free(keys);
  free(r.data);
  json_decref(root);
  return 0;
}
