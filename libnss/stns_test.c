#include "stns.h"
#include "toml.h"
#include <criterion/criterion.h>

Test(misc, failing)
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
