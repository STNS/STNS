#include "stns.h"
#include "stns_test.h"

void readfile(char *file, char **result)
{
  FILE *fp;
  char buff[MAXBUF];
  int len;
  int total;

  fp = fopen(file, "r");
  if (fp == NULL) {
    return;
  }

  total   = 0;
  *result = NULL;
  for (;;) {
    if (fgets(buff, MAXBUF, fp) == NULL) {
      break;
    }
    len = strlen(buff);

    if (!*result) {
      *result = (char *)malloc(total + len + 1);
    } else {
      *result = realloc(*result, total + len + 1);
    }
    if (*result == NULL) {
      break;
    }
    strcpy(*result + total, buff);
    total += len;
  }
  fclose(fp);
}

Test(stns_load_config, load_ok)
{
  char *f = "test/stns.conf";
  stns_conf_t c;
  stns_load_config(f, &c);

  cr_assert_str_eq(c.api_endpoint, "http://<server-ip>:1104/v2");
  cr_assert_str_eq(c.auth_token, "xxxxxxxxxxxxxxx");
  cr_assert_str_eq(c.user, "test_user");
  cr_assert_str_eq(c.password, "test_password");
  cr_assert_str_eq(c.chain_ssh_wrapper, "/usr/libexec/openssh/ssh-ldap-wrapper");
  cr_assert_str_eq(c.query_wrapper, "/usr/local/bin/stns-wrapper");
  cr_assert_str_eq(c.http_proxy, "http://your.proxy.com");
  cr_assert_eq(c.ssl_verify, 1);
  cr_assert_eq(c.uid_shift, 1000);
  cr_assert_eq(c.gid_shift, 2000);
  cr_assert_eq(c.request_timeout, 3);
  cr_assert_eq(c.request_retry, 3);
  cr_assert_eq(c.negative_cache_ttl, 10);
}

Test(stns_request, http_request)
{
  char expect_body[1024];
  stns_conf_t c;
  stns_response_t r;

  c.api_endpoint    = "https://httpbin.org";
  c.cache_dir       = "/var/cache/stns";
  c.cache           = 0;
  c.user            = NULL;
  c.password        = NULL;
  c.query_wrapper   = NULL;
  c.request_timeout = 3;
  c.request_retry   = 3;
  c.auth_token      = NULL;
  stns_request(&c, "user-agent", &r);

  sprintf(expect_body, "{\n  \"user-agent\": \"%s\"\n}\n", STNS_VERSION_WITH_NAME);
  cr_assert_str_eq(r.data, expect_body);
}

Test(stns_request, http_cache)
{
  struct stat st;
  stns_conf_t c;
  stns_response_t r;
  char *path = "/var/cache/stns/get%3Fexample";

  c.api_endpoint    = "https://httpbin.org";
  c.cache_dir       = "/var/cache/stns";
  c.cache           = 1;
  c.cache_ttl       = 2;
  c.user            = NULL;
  c.password        = NULL;
  c.query_wrapper   = NULL;
  c.request_timeout = 3;
  c.request_retry   = 3;
  c.auth_token      = NULL;

  stns_request(&c, "get?example", &r);
  cr_assert_eq(stat(path, &st), 0);
  sleep(5);
  // deleted by thread
  stns_request(&c, "get?notfound", &r);
  cr_assert_eq(stat(path, &st), -1);
}

Test(stns_request, wrapper_request_ok)
{
  stns_conf_t c;
  stns_response_t r;
  int res;

  c.query_wrapper = "test/dummy_arg.sh";
  c.cache_dir     = "/var/cache/stns";
  c.cache         = 0;

  res = stns_request(&c, "test", &r);
  cr_assert_str_eq(r.data, "ok\n");
  cr_assert_eq(res, 1);
}

Test(stns_request, wrapper_request_ng)
{
  stns_conf_t c;
  stns_response_t r;
  int res;

  c.cache         = 0;
  c.query_wrapper = "test/dummy_arg.sh";

  res = stns_request(&c, NULL, &r);
  cr_assert_eq(res, 22);
}

Test(stns_request_available, ok)
{
  char expect_body[1024];
  stns_conf_t c;
  stns_response_t r;

  c.request_locktime = 1;
  stns_make_lockfile(STNS_LOCK_FILE);
  cr_assert_eq(stns_request_available(STNS_LOCK_FILE, &c), 0);
  sleep(3);
  cr_assert_eq(stns_request_available(STNS_LOCK_FILE, &c), 1);
}

Test(stns_exec_cmd, ok)
{
  char expect_body[1024];
  stns_response_t result;
  int r = stns_exec_cmd("test/dummy.sh", "test", &result);

  cr_assert_eq(r, 1);
  cr_expect_str_eq(result.data, "aaabbbccc\nddd\n");
  free(result.data);
}

Test(query_available, ok)
{
  set_user_highest_id(10);
  set_user_lowest_id(3);
  set_group_highest_id(10);
  set_group_lowest_id(3);

  cr_assert_eq(stns_user_highest_query_available(1), 1);
  cr_assert_eq(stns_user_highest_query_available(11), 0);
  cr_assert_eq(stns_user_lowest_query_available(4), 1);
  cr_assert_eq(stns_user_lowest_query_available(2), 0);

  cr_assert_eq(stns_group_highest_query_available(1), 1);
  cr_assert_eq(stns_group_highest_query_available(11), 0);
  cr_assert_eq(stns_group_lowest_query_available(4), 1);
  cr_assert_eq(stns_group_lowest_query_available(2), 0);
}
