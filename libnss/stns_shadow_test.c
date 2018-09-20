#include "stns_test.h"
Test(ensure_spwd_by_name, ok)
{
  char *f = "test/example1.json";
  char *json;
  int code;
  struct spwd spbuf;
  char buffer[MAXBUF];
  stns_conf_t c;
  stns_response_t r;

  readfile(f, &json);
  code = ensure_spwd_by_name(json, &c, "user1", &spbuf, buffer, MAXBUF, 0);
  cr_assert_eq(code, NSS_STATUS_SUCCESS);
  cr_assert_str_eq(spbuf.sp_namp, "user1");
  cr_assert_str_eq(spbuf.sp_pwdp, "test");
  cr_assert_eq(spbuf.sp_lstchg, -1);
  cr_assert_eq(spbuf.sp_min, -1);
  cr_assert_eq(spbuf.sp_max, -1);
  cr_assert_eq(spbuf.sp_warn, -1);
  cr_assert_eq(spbuf.sp_inact, -1);
  cr_assert_eq(spbuf.sp_expire, -1);
  cr_assert_eq(spbuf.sp_flag, ~0ul);

  readfile(f, &json);
  code = ensure_spwd_by_name(json, &c, "user3", &spbuf, buffer, MAXBUF, 0);
  cr_assert_eq(code, NSS_STATUS_NOTFOUND);

  char *n = malloc(1);
  strcpy(n, "");
  code = ensure_spwd_by_name(n, &c, "user3", &spbuf, buffer, MAXBUF, 0);
  cr_assert_eq(code, NSS_STATUS_UNAVAIL);
}

Test(ensure_spwd_by_uid, ok)
{
  char *f = "test/example1.json";
  char *json;
  int code;
  struct spwd spbuf;
  char buffer[MAXBUF];
  stns_conf_t c;
  stns_response_t r;
  c.uid_shift = 0;

  readfile(f, &json);
  code = ensure_spwd_by_uid(json, &c, 1, &spbuf, buffer, MAXBUF, 0);
  cr_assert_eq(code, NSS_STATUS_SUCCESS);
  cr_assert_str_eq(spbuf.sp_namp, "user1");
  cr_assert_str_eq(spbuf.sp_pwdp, "test");
  cr_assert_eq(spbuf.sp_lstchg, -1);
  cr_assert_eq(spbuf.sp_min, -1);
  cr_assert_eq(spbuf.sp_max, -1);
  cr_assert_eq(spbuf.sp_warn, -1);
  cr_assert_eq(spbuf.sp_inact, -1);
  cr_assert_eq(spbuf.sp_expire, -1);
  cr_assert_eq(spbuf.sp_flag, ~0ul);

  readfile(f, &json);
  code = ensure_spwd_by_uid(json, &c, 3, &spbuf, buffer, MAXBUF, 0);
  cr_assert_eq(code, NSS_STATUS_NOTFOUND);

  char *n = malloc(1);
  strcpy(n, "");
  code = ensure_spwd_by_uid(n, &c, 3, &spbuf, buffer, MAXBUF, 0);
  cr_assert_eq(code, NSS_STATUS_UNAVAIL);
}

Test(inner_nss_stns_setspent, ok)
{
  char *f = "test/example1.json";
  char *json;
  int code;
  struct spwd spbuf;
  char buffer[MAXBUF];
  stns_conf_t c;
  stns_response_t r;

  c.uid_shift = 0;
  c.gid_shift = 0;
  readfile(f, &json);
  code = inner_nss_stns_setspent(json, &c);
  cr_assert_eq(code, NSS_STATUS_SUCCESS);

  char *n = malloc(1);
  strcpy(n, "");
  code = inner_nss_stns_setspent(n, &c);
  cr_assert_eq(code, NSS_STATUS_UNAVAIL);
  _nss_stns_endspent();
}

Test(inner_nss_stns_getspent_r, ok)
{
  char *f = "test/example1.json";
  char *json;
  int code;
  int errnop = 0;
  struct spwd spbuf;
  char buffer[MAXBUF];
  stns_conf_t c;
  stns_response_t r;

  readfile(f, &json);
  code = inner_nss_stns_setspent(json, &c);
  cr_assert_eq(code, NSS_STATUS_SUCCESS);

  code = inner_nss_stns_getspent_r(&c, &spbuf, buffer, MAXBUF, &errnop);
  cr_assert_eq(code, NSS_STATUS_SUCCESS);
  cr_assert_str_eq(spbuf.sp_namp, "user1");
  cr_assert_str_eq(spbuf.sp_pwdp, "test");
  cr_assert_eq(spbuf.sp_lstchg, -1);
  cr_assert_eq(spbuf.sp_min, -1);
  cr_assert_eq(spbuf.sp_max, -1);
  cr_assert_eq(spbuf.sp_warn, -1);
  cr_assert_eq(spbuf.sp_inact, -1);
  cr_assert_eq(spbuf.sp_expire, -1);
  cr_assert_eq(spbuf.sp_flag, ~0ul);

  code = inner_nss_stns_getspent_r(&c, &spbuf, buffer, MAXBUF, &errnop);
  cr_assert_eq(code, NSS_STATUS_SUCCESS);
  cr_assert_str_eq(spbuf.sp_namp, "user2");
  cr_assert_str_eq(spbuf.sp_pwdp, "!!");
  cr_assert_eq(spbuf.sp_lstchg, -1);
  cr_assert_eq(spbuf.sp_min, -1);
  cr_assert_eq(spbuf.sp_max, -1);
  cr_assert_eq(spbuf.sp_warn, -1);
  cr_assert_eq(spbuf.sp_inact, -1);
  cr_assert_eq(spbuf.sp_expire, -1);
  cr_assert_eq(spbuf.sp_flag, ~0ul);

  code = inner_nss_stns_getspent_r(&c, &spbuf, buffer, MAXBUF, &errnop);
  cr_assert_eq(code, NSS_STATUS_NOTFOUND);
  _nss_stns_endspent();
}
