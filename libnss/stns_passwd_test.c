#include "stns.h"
#include "stns_passwd.h"
#include "toml.h"
#include <criterion/criterion.h>

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

Test(ensure_passwd_by_name, ok)
{
  char *f = "test/example1.json";
  char *json;
  int code;
  struct passwd pwd;
  char buffer[MAXBUF];
  stns_conf_t c;
  stns_http_response_t r;
  readfile(f, &json);
  c.uid_shift = 0;
  c.gid_shift = 0;

  code = ensure_passwd_by_name(json, &c, "user1", &pwd, buffer, MAXBUF, 0);
  cr_assert_eq(code, NSS_STATUS_SUCCESS);
  cr_assert_str_eq(pwd.pw_name, "user1");
  cr_assert_eq(pwd.pw_uid, 1);
  cr_assert_eq(pwd.pw_gid, 1);
  cr_assert_str_eq(pwd.pw_passwd, "x");
  cr_assert_str_eq(pwd.pw_gecos, "test");
  cr_assert_str_eq(pwd.pw_shell, "/bin/sh");
  cr_assert_str_eq(pwd.pw_dir, "/home/admin/user1");

  // id shift
  readfile(f, &json);
  c.uid_shift = 100;
  c.gid_shift = 200;
  ensure_passwd_by_name(json, &c, "user1", &pwd, buffer, MAXBUF, 0);
  cr_assert_eq(code, NSS_STATUS_SUCCESS);
  cr_assert_eq(pwd.pw_uid, 101);
  cr_assert_eq(pwd.pw_gid, 201);

  // use default shell with dir
  readfile(f, &json);
  c.uid_shift = 0;
  c.gid_shift = 0;
  code        = ensure_passwd_by_name(json, &c, "user2", &pwd, buffer, MAXBUF, 0);
  cr_assert_eq(code, NSS_STATUS_SUCCESS);
  cr_assert_str_eq(pwd.pw_name, "user2");
  cr_assert_eq(pwd.pw_uid, 2);
  cr_assert_eq(pwd.pw_gid, 2);
  cr_assert_str_eq(pwd.pw_passwd, "x");
  cr_assert_str_eq(pwd.pw_gecos, "test");
  cr_assert_str_eq(pwd.pw_shell, "/bin/bash");
  cr_assert_str_eq(pwd.pw_dir, "/home/user2");

  readfile(f, &json);
  c.uid_shift = 0;
  c.gid_shift = 0;
  code        = ensure_passwd_by_name(json, &c, "user3", &pwd, buffer, MAXBUF, 0);
  cr_assert_eq(code, NSS_STATUS_NOTFOUND);

  char *n = malloc(1);
  strcpy(n, "");
  code = ensure_passwd_by_name(n, &c, "user3", &pwd, buffer, MAXBUF, 0);
  cr_assert_eq(code, NSS_STATUS_UNAVAIL);
}

Test(ensure_passwd_by_uid, ok)
{
  char *f = "test/example1.json";
  char *json;
  int code;
  struct passwd pwd;
  char buffer[MAXBUF];
  stns_conf_t c;
  stns_http_response_t r;
  readfile(f, &json);
  c.uid_shift = 0;
  c.gid_shift = 0;

  code = ensure_passwd_by_uid(json, &c, 1, &pwd, buffer, MAXBUF, 0);
  cr_assert_eq(code, NSS_STATUS_SUCCESS);
  cr_assert_str_eq(pwd.pw_name, "user1");
  cr_assert_eq(pwd.pw_uid, 1);
  cr_assert_eq(pwd.pw_gid, 1);
  cr_assert_str_eq(pwd.pw_passwd, "x");
  cr_assert_str_eq(pwd.pw_gecos, "test");
  cr_assert_str_eq(pwd.pw_shell, "/bin/sh");
  cr_assert_str_eq(pwd.pw_dir, "/home/admin/user1");

  // id shift
  readfile(f, &json);
  c.uid_shift = 100;
  c.gid_shift = 200;
  ensure_passwd_by_uid(json, &c, 1, &pwd, buffer, MAXBUF, 0);
  cr_assert_eq(code, NSS_STATUS_SUCCESS);
  cr_assert_eq(pwd.pw_uid, 101);
  cr_assert_eq(pwd.pw_gid, 201);

  // use default shell with dir
  readfile(f, &json);
  c.uid_shift = 0;
  c.gid_shift = 0;
  code        = ensure_passwd_by_uid(json, &c, 2, &pwd, buffer, MAXBUF, 0);
  cr_assert_eq(code, NSS_STATUS_SUCCESS);
  cr_assert_str_eq(pwd.pw_name, "user2");
  cr_assert_eq(pwd.pw_uid, 2);
  cr_assert_eq(pwd.pw_gid, 2);
  cr_assert_str_eq(pwd.pw_passwd, "x");
  cr_assert_str_eq(pwd.pw_gecos, "test");
  cr_assert_str_eq(pwd.pw_shell, "/bin/bash");
  cr_assert_str_eq(pwd.pw_dir, "/home/user2");

  readfile(f, &json);
  c.uid_shift = 0;
  c.gid_shift = 0;
  code        = ensure_passwd_by_uid(json, &c, 3, &pwd, buffer, MAXBUF, 0);
  cr_assert_eq(code, NSS_STATUS_NOTFOUND);

  char *n = malloc(1);
  strcpy(n, "");
  code = ensure_passwd_by_uid(n, &c, 3, &pwd, buffer, MAXBUF, 0);
  cr_assert_eq(code, NSS_STATUS_UNAVAIL);
}

Test(_nss_stns_setpwent, ok)
{
  char *f = "test/example1.json";
  char *json;
  int code;
  struct passwd pwd;
  char buffer[MAXBUF];
  stns_conf_t c;
  stns_http_response_t r;

  c.uid_shift = 0;
  c.gid_shift = 0;
  readfile(f, &json);
  code = inner_nss_stns_setpwent(json, &c);
  cr_assert_eq(code, NSS_STATUS_SUCCESS);

  char *n = malloc(1);
  strcpy(n, "");
  code = inner_nss_stns_setpwent(n, &c);
  cr_assert_eq(code, NSS_STATUS_UNAVAIL);
}
