#include "stns.h"
#include "stns_test.h"

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
  cr_assert_str_eq(c.http_proxy, "http://your.proxy.com");
  cr_assert_eq(c.ssl_verify, 1);
  cr_assert_eq(c.uid_shift, 1000);
  cr_assert_eq(c.gid_shift, 2000);
  cr_assert_eq(c.request_timeout, 3);
  cr_assert_eq(c.request_retry, 3);
}

Test(stns_request, request_ok)
{
  char *f = "test/stns.conf";
  char expect_body[1024];
  stns_conf_t c;
  stns_http_response_t r;

  c.api_endpoint    = "https://httpbin.org";
  c.request_timeout = 3;
  c.request_retry   = 3;
  c.auth_token      = NULL;
  stns_request(&c, "user-agent", &r);

  cr_assert_eq(r.status_code, (long *)200);
  sprintf(expect_body, "{\n  \"user-agent\": \"%s\"\n}\n", STNS_VERSION_WITH_NAME);
  cr_assert_str_eq(r.data, expect_body);
}

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
