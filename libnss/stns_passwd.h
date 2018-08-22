#ifndef STNS_PWD_H
#define STNS_PWD_H

extern enum nss_status ensure_passwd_by_name(char *, stns_conf_t *, const char *, struct passwd *, char *, size_t,
                                             int *);
extern enum nss_status ensure_passwd_by_uid(char *, stns_conf_t *, uid_t uid, struct passwd *, char *, size_t, int *);
extern enum nss_status inner_nss_stns_setpwent(char *, stns_conf_t *);
extern enum nss_status inner_nss_stns_getpwent_r(stns_conf_t *, struct passwd *, char *, size_t, int *);
extern enum nss_status _nss_stns_endpwent(void);
#endif /* STNS_PWD_H */
