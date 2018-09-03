#include "stns_test.h"

Test(ensure_group_by_name, ok)
{
  char *f = "test/example2.json";
  char *json;
  int code;
  struct group grd;
  char buffer[MAXBUF];
  stns_conf_t c;
  stns_response_t r;
  c.gid_shift = 0;

  readfile(f, &json);
  code = ensure_group_by_name(json, &c, "group1", &grd, buffer, MAXBUF, 0);
  cr_assert_eq(code, NSS_STATUS_SUCCESS);
  cr_assert_str_eq(grd.gr_name, "group1");
  cr_assert_eq(grd.gr_gid, 1);
  cr_assert_str_eq(grd.gr_passwd, "x");
  cr_assert_str_eq(grd.gr_mem[0], "test");

  readfile(f, &json);
  code = ensure_group_by_name(json, &c, "group2", &grd, buffer, MAXBUF, 0);
  cr_assert_eq(code, NSS_STATUS_SUCCESS);
  cr_assert_str_eq(grd.gr_name, "group2");
  cr_assert_eq(grd.gr_gid, 2);
  cr_assert_str_eq(grd.gr_passwd, "x");
  cr_assert_str_eq(grd.gr_mem[0], "foo");
  cr_assert_str_eq(grd.gr_mem[1], "bar");

  // id shift
  readfile(f, &json);
  c.gid_shift = 200;
  ensure_group_by_name(json, &c, "group1", &grd, buffer, MAXBUF, 0);
  cr_assert_eq(code, NSS_STATUS_SUCCESS);
  cr_assert_eq(grd.gr_gid, 201);

  readfile(f, &json);
  c.gid_shift = 0;
  code        = ensure_group_by_name(json, &c, "group3", &grd, buffer, MAXBUF, 0);
  cr_assert_eq(code, NSS_STATUS_NOTFOUND);

  char *n = malloc(1);
  strcpy(n, "");
  code = ensure_group_by_name(n, &c, "group3", &grd, buffer, MAXBUF, 0);
  cr_assert_eq(code, NSS_STATUS_UNAVAIL);
}

Test(ensure_group_by_gid, ok)
{
  char *f = "test/example2.json";
  char *json;
  int code;
  struct group grd;
  char buffer[MAXBUF];
  stns_conf_t c;
  stns_response_t r;
  c.gid_shift = 0;

  readfile(f, &json);
  code = ensure_group_by_gid(json, &c, 1, &grd, buffer, MAXBUF, 0);
  cr_assert_eq(code, NSS_STATUS_SUCCESS);
  cr_assert_str_eq(grd.gr_name, "group1");
  cr_assert_eq(grd.gr_gid, 1);
  cr_assert_str_eq(grd.gr_passwd, "x");
  cr_assert_str_eq(grd.gr_mem[0], "test");

  readfile(f, &json);
  code = ensure_group_by_gid(json, &c, 2, &grd, buffer, MAXBUF, 0);
  cr_assert_eq(code, NSS_STATUS_SUCCESS);
  cr_assert_str_eq(grd.gr_name, "group2");
  cr_assert_eq(grd.gr_gid, 2);
  cr_assert_str_eq(grd.gr_passwd, "x");
  cr_assert_str_eq(grd.gr_mem[0], "foo");
  cr_assert_str_eq(grd.gr_mem[1], "bar");

  // id shift
  readfile(f, &json);
  c.gid_shift = 200;
  ensure_group_by_gid(json, &c, 1, &grd, buffer, MAXBUF, 0);
  cr_assert_eq(code, NSS_STATUS_SUCCESS);
  cr_assert_eq(grd.gr_gid, 201);

  readfile(f, &json);
  c.gid_shift = 0;
  code        = ensure_group_by_gid(json, &c, 3, &grd, buffer, MAXBUF, 0);
  cr_assert_eq(code, NSS_STATUS_NOTFOUND);

  char *n = malloc(1);
  strcpy(n, "");
  code = ensure_group_by_gid(n, &c, 3, &grd, buffer, MAXBUF, 0);
  cr_assert_eq(code, NSS_STATUS_UNAVAIL);
}

Test(inner_nss_stns_setgrent, ok)
{
  char *f = "test/example2.json";
  char *json;
  int code;
  struct group grd;
  char buffer[MAXBUF];
  stns_conf_t c;
  stns_response_t r;

  c.gid_shift = 0;
  readfile(f, &json);
  code = inner_nss_stns_setgrent(json, &c);
  cr_assert_eq(code, NSS_STATUS_SUCCESS);

  char *n = malloc(1);
  strcpy(n, "");
  code = inner_nss_stns_setgrent(n, &c);
  cr_assert_eq(code, NSS_STATUS_UNAVAIL);
  _nss_stns_endgrent();
}

Test(inner_nss_stns_getgrent_r, ok)
{
  char *f = "test/example2.json";
  char *json;
  int code;
  int errnop = 0;
  struct group grd;
  char buffer[MAXBUF];
  stns_conf_t c;
  stns_response_t r;

  c.gid_shift = 0;
  readfile(f, &json);
  code = inner_nss_stns_setgrent(json, &c);
  cr_assert_eq(code, NSS_STATUS_SUCCESS);

  code = inner_nss_stns_getgrent_r(&c, &grd, buffer, MAXBUF, &errnop);
  cr_assert_eq(code, NSS_STATUS_SUCCESS);
  cr_assert_eq(code, NSS_STATUS_SUCCESS);
  cr_assert_str_eq(grd.gr_name, "group1");
  cr_assert_eq(grd.gr_gid, 1);
  cr_assert_str_eq(grd.gr_passwd, "x");

  code = inner_nss_stns_getgrent_r(&c, &grd, buffer, MAXBUF, &errnop);
  cr_assert_eq(code, NSS_STATUS_SUCCESS);
  cr_assert_str_eq(grd.gr_name, "group2");
  cr_assert_eq(grd.gr_gid, 2);
  cr_assert_str_eq(grd.gr_passwd, "x");

  code = inner_nss_stns_getgrent_r(&c, &grd, buffer, MAXBUF, &errnop);
  cr_assert_eq(code, NSS_STATUS_NOTFOUND);
  _nss_stns_endgrent();
}
