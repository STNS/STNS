#ifndef STNS_GROUP_H
#define STNS_GROUP_H

extern enum nss_status ensure_group_by_name(char *, stns_conf_t *, const char *, struct group *, char *, size_t, int *);
extern enum nss_status ensure_group_by_gid(char *, stns_conf_t *, gid_t gid, struct group *, char *, size_t, int *);
extern enum nss_status inner_nss_stns_setgrent(char *, stns_conf_t *);
extern enum nss_status inner_nss_stns_getgrent_r(stns_conf_t *, struct group *, char *, size_t, int *);
extern enum nss_status _nss_stns_endgrent(void);
#endif /* STNS_GROUP_H */
