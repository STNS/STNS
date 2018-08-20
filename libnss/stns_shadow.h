#ifndef STNS_SPWD_H
#define STNS_SPWD_H

extern enum nss_status ensure_spwd_by_name(char *, stns_conf_t *, const char *, struct spwd *, char *, size_t, int *);
extern enum nss_status ensure_spwd_by_uid(char *, stns_conf_t *, uid_t uid, struct spwd *, char *, size_t, int *);
extern enum nss_status inner_nss_stns_setspent(char *, stns_conf_t *);
extern enum nss_status inner_nss_stns_getspent_r(stns_conf_t *, struct spwd *, char *, size_t, int *);
extern enum nss_status _nss_stns_endspent(void);
#endif /* STNS_SPWD_H */
