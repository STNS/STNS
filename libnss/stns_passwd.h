#ifndef STNS_PWD_H
#define STNS_PWD_H

extern enum nss_status ensure_passwd_by_name(char *, stns_conf_t *, const char *, struct passwd *, char *, size_t,
                                             int *);
extern enum nss_status ensure_passwd_by_uid(char *, stns_conf_t *, uid_t uid, struct passwd *, char *, size_t, int *);
#endif /* STNS_PWD_H */
