package server

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"sort"
	"strconv"
	"strings"
	"syscall"

	"github.com/STNS/STNS/v2/model"
	"github.com/STNS/STNS/v2/stns"
	"github.com/facebookgo/pidfile"
	"github.com/labstack/gommon/log"
	"github.com/lestrrat/go-server-starter/listener"
	"github.com/nmcclain/ldap"
	"github.com/tredoe/osutil/user/crypt"
	_ "github.com/tredoe/osutil/user/crypt/md5_crypt"
	_ "github.com/tredoe/osutil/user/crypt/sha256_crypt"
	_ "github.com/tredoe/osutil/user/crypt/sha512_crypt"
)

type ldapServer struct {
	baseServer
}

type ldapHandler struct {
	backend model.Backend
	logger  *log.Logger
	config  *stns.Config
}

func (h ldapHandler) Bind(bindDN, rawPassword string, conn net.Conn) (ldap.LDAPResultCode, error) {
	bindDN = strings.ToLower(bindDN)
	baseDN := strings.ToLower("," + h.config.LDAP.BaseDN)

	if !strings.HasSuffix(bindDN, baseDN) {
		h.logger.Warn(fmt.Sprintf("Bind Error: BindDN %s not our BaseDN %s", bindDN, h.config.LDAP.BaseDN))
		return ldap.LDAPResultInvalidCredentials, nil
	}

	parts := strings.Split(strings.TrimSuffix(bindDN, baseDN), ",")
	groupName := ""
	userName := ""
	if len(parts) == 1 {
		userName = strings.TrimPrefix(parts[0], "cn=")
	} else if len(parts) == 2 {
		userName = strings.TrimPrefix(parts[0], "cn=")
		groupName = strings.TrimPrefix(parts[1], "ou=")
	} else {
		h.logger.Warn(fmt.Sprintf("Bind Error: BindDN %s should have only one or two parts (has %d)", bindDN, len(parts)))
		return ldap.LDAPResultInvalidCredentials, nil
	}

	var user *model.User
	var group *model.Group

	users, err := h.backend.FindUserByName(userName)
	if err != nil {
		h.logger.Warn(fmt.Sprintf("Bind Error: User %s not found.", userName))
		return ldap.LDAPResultInvalidCredentials, nil
	}
	user = users[userName].(*model.User)

	if groupName != "" {
		groups, err := h.backend.FindGroupByName(groupName)
		if err != nil {
			h.logger.Warn(fmt.Sprintf("Bind Error: Group %s not found.", groupName))
			return ldap.LDAPResultInvalidCredentials, nil
		}

		group = groups[groupName].(*model.Group)

		if user.GroupID != group.ID {
			h.logger.Warn(fmt.Sprintf("Bind Error: User %s primary group is not %s.", userName, groupName))
			return ldap.LDAPResultInvalidCredentials, nil
		}

	}
	if user.Password != "" {
		c := crypt.NewFromHash(user.Password)
		if c.Verify(user.Password, []byte(rawPassword)) != nil {
			h.logger.Warn(fmt.Sprintf("Bind Error: invalid credentials as %s from %s", bindDN, conn.RemoteAddr().String()))
			return ldap.LDAPResultInvalidCredentials, nil
		}
	}
	return ldap.LDAPResultSuccess, nil

}

func (h ldapHandler) Search(bindDN string, searchReq ldap.SearchRequest, conn net.Conn) (ldap.ServerSearchResult, error) {
	bindDN = strings.ToLower(bindDN)
	baseDN := strings.ToLower("," + h.config.LDAP.BaseDN)
	searchBaseDN := strings.ToLower(searchReq.BaseDN)

	// validate the user is authenticated and has appropriate access
	if len(bindDN) < 1 {
		return ldap.ServerSearchResult{ResultCode: ldap.LDAPResultInsufficientAccessRights}, fmt.Errorf("Search Error: Anonymous BindDN not allowed %s", bindDN)
	}
	if !strings.HasSuffix(bindDN, baseDN) {
		return ldap.ServerSearchResult{ResultCode: ldap.LDAPResultInsufficientAccessRights}, fmt.Errorf("Search Error: BindDN %s not in our BaseDN %s", bindDN, h.config.LDAP.BaseDN)
	}
	if !strings.HasSuffix(searchBaseDN, baseDN) {
		return ldap.ServerSearchResult{ResultCode: ldap.LDAPResultInsufficientAccessRights}, fmt.Errorf("Search Error: search BaseDN %s is not in our BaseDN %s", searchBaseDN, baseDN)
	}

	entries := []*ldap.Entry{}
	filterEntity, err := ldap.GetFilterObjectClass(searchReq.Filter)
	if err != nil {
		return ldap.ServerSearchResult{
			ResultCode: ldap.LDAPResultOperationsError,
		}, fmt.Errorf("Search Error: error parsing filter: %s", searchReq.Filter)
	}

	groups, err := h.backend.Groups()
	if err != nil {
		return ldap.ServerSearchResult{
			ResultCode: ldap.LDAPResultOperationsError,
		}, fmt.Errorf("Search Error: can't fetch groups: %s [%s]", filterEntity, searchReq.Filter)
	}

	users, err := h.backend.Users()
	if err != nil {
		return ldap.ServerSearchResult{
			ResultCode: ldap.LDAPResultOperationsError,
		}, fmt.Errorf("Search Error: can't fetch users: %s [%s]", filterEntity, searchReq.Filter)
	}

	switch filterEntity {
	default:
		return ldap.ServerSearchResult{
			ResultCode: ldap.LDAPResultOperationsError,
		}, fmt.Errorf("Search Error: unhandled filter type: %s [%s]", filterEntity, searchReq.Filter)
	case "posixgroup":
		for _, g := range groups {
			attrs := []*ldap.EntryAttribute{}
			attrs = append(attrs, &ldap.EntryAttribute{"cn", []string{g.GetName()}})
			attrs = append(attrs, &ldap.EntryAttribute{"description", []string{fmt.Sprintf("%s via LDAP", g.GetName())}})
			attrs = append(attrs, &ldap.EntryAttribute{"gidNumber", []string{fmt.Sprintf("%d", g.GetID())}})
			attrs = append(attrs, &ldap.EntryAttribute{"objectClass", []string{"posixGroup"}})
			attrs = append(attrs, &ldap.EntryAttribute{"memberUid", h.getGroupMemberNames(g.(*model.Group), users, h.config.LDAP.BaseDN)})
			dn := fmt.Sprintf("cn=%s,ou=groups,%s", g.GetName(), h.config.LDAP.BaseDN)

			entries = append(entries, &ldap.Entry{dn, attrs})
		}
	case "posixaccount", "":
		for _, u := range users {
			userGroup, err := h.backend.FindGroupByID(u.(*model.User).GroupID)
			if err != nil {
				return ldap.ServerSearchResult{
					ResultCode: ldap.LDAPResultOperationsError,
				}, fmt.Errorf("Search Error: can't fetch primary group : %s [%s]", filterEntity, searchReq.Filter)

			}

			var group *model.Group
			for _, g := range userGroup {
				group = g.(*model.Group)
			}

			memberOf := []int{}
			for _, g := range groups {
				// find other group which the user belongs
				for _, gm := range g.(*model.Group).Users {
					if gm == u.GetName() {
						memberOf = append(memberOf, g.GetID())
					}
				}
			}
			memberOf = append(memberOf, u.(*model.User).GroupID)

			attrs := []*ldap.EntryAttribute{}
			attrs = append(attrs, &ldap.EntryAttribute{"cn", []string{u.GetName()}})
			attrs = append(attrs, &ldap.EntryAttribute{"uid", []string{u.GetName()}})
			attrs = append(attrs, &ldap.EntryAttribute{"givenName", []string{u.GetName()}})
			attrs = append(attrs, &ldap.EntryAttribute{"ou", []string{group.Name}})
			attrs = append(attrs, &ldap.EntryAttribute{"uidNumber", []string{fmt.Sprintf("%d", u.GetID())}})
			attrs = append(attrs, &ldap.EntryAttribute{"accountStatus", []string{"active"}})
			attrs = append(attrs, &ldap.EntryAttribute{"objectClass", []string{"posixAccount"}})
			attrs = append(attrs, &ldap.EntryAttribute{"loginShell", []string{u.(*model.User).Shell}})
			attrs = append(attrs, &ldap.EntryAttribute{"homeDirectory", []string{u.(*model.User).Directory}})
			attrs = append(attrs, &ldap.EntryAttribute{"description", []string{fmt.Sprintf("%s via LDAP", u.GetName())}})
			attrs = append(attrs, &ldap.EntryAttribute{"gecos", []string{u.(*model.User).Gecos}})
			attrs = append(attrs, &ldap.EntryAttribute{"gidNumber", []string{strconv.Itoa(u.(*model.User).GroupID)}})
			attrs = append(attrs, &ldap.EntryAttribute{"memberOf", h.getGroupDNs(memberOf, groups, h.config.LDAP.BaseDN)})
			attrs = append(attrs, &ldap.EntryAttribute{"sshPublicKey", u.(*model.User).Keys})
			dn := fmt.Sprintf("cn=%s,ou=%s,%s", u.GetName(), group.Name, h.config.LDAP.BaseDN)
			entries = append(entries, &ldap.Entry{dn, attrs})
		}
	}

	return ldap.ServerSearchResult{entries, []string{}, []ldap.Control{}, ldap.LDAPResultSuccess}, nil
}

func newLDAPServer(conf *stns.Config, backend model.Backend, logger *log.Logger) (*ldapServer, error) {
	s := &ldapServer{
		baseServer{
			config:  conf,
			backend: backend,
			logger:  logger,
		},
	}
	return s, nil
}

func (s *ldapServer) Run() error {
	quit := make(chan bool)
	ld := ldap.NewServer()
	ld.QuitChannel(quit)
	ld.EnforceLDAP = true
	h := ldapHandler{
		backend: s.backend,
		config:  s.config,
	}
	ld.BindFunc("", h)
	ld.SearchFunc("", h)

	if err := pidfile.Write(); err != nil {
		return err
	}
	defer func() {
		if err := os.Remove(pidfile.GetPidfilePath()); err != nil {
			s.logger.Errorf("Error removing %s: %s", pidfile.GetPidfilePath(), err)
		}
	}()

	p := strconv.Itoa(s.config.Port)

	lnstr := ":" + p
	if os.Getenv("STNS_LISTEN") != "" {
		lnstr = os.Getenv("STNS_LISTEN")
	}
	if s.config.UseServerStarter {
		listeners, err := listener.ListenAll()
		if listeners == nil || err != nil {
			return err
		}
		lnstr = listeners[0].Addr().String()
	}

	s.logger.Info("start ldap server")

	go func() {
		sigQuit := make(chan os.Signal)
		signal.Notify(sigQuit, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)
		<-sigQuit
		quit <- true
	}()
	if s.config.TLS != nil && s.config.TLS.Cert != "" && s.config.TLS.Key != "" {
		if err := ld.ListenAndServeTLS(lnstr, s.config.TLS.Cert, s.config.TLS.Key); err != nil {
			return err
		}
	} else {
		if err := ld.ListenAndServe(lnstr); err != nil {
			return err
		}
	}
	return nil
}

func (h ldapHandler) getGroupMemberNames(group *model.Group, users map[string]model.UserGroup, baseName string) []string {
	groupMemberNames := map[string]bool{}

	// this group belongs user(primary)
	for _, u := range users {
		if u.(*model.User).GroupID == group.ID {
			groupMemberNames[u.GetName()] = true
		}
	}

	// this group belongs user(other)
	for _, memberName := range group.Users {
		if memberName == "" {
			break
		}
		groupMemberNames[memberName] = true

	}
	g := []string{}
	for k, _ := range groupMemberNames {
		g = append(g, k)
	}

	sort.Strings(g)

	return g
}
func (h ldapHandler) getGroupDNs(gids []int, usergroups map[string]model.UserGroup, baseName string) []string {
	groups := make(map[string]bool)
	for _, gid := range gids {
		for _, g := range usergroups {
			if g.GetID() == gid {
				dn := fmt.Sprintf("cn=%s,ou=groups,%s", g.GetName(), baseName)
				groups[dn] = true
			}
		}
	}

	g := []string{}
	for k, _ := range groups {
		g = append(g, k)
	}

	sort.Strings(g)

	return g
}
