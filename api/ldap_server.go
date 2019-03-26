package api

// refs: https://github.com/glauth/glauth/blob/master/configbackend.go
/*
	MIT License

	Copyright (c) 2018 Ned McClain & Ben Yanke

	Permission is hereby granted, free of charge, to any person obtaining a copy
	of this software and associated documentation files (the "Software"), to deal
	in the Software without restriction, including without limitation the rights
	to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
	copies of the Software, and to permit persons to whom the Software is
	furnished to do so, subject to the following conditions:

	The above copyright notice and this permission notice shall be included in all
	copies or substantial portions of the Software.

	THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
	IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
	FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
	AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
	LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
	OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
	SOFTWARE.
*/
import (
	"fmt"
	"net"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/STNS/STNS/model"
	"github.com/STNS/STNS/stns"
	"github.com/facebookgo/pidfile"
	"github.com/labstack/gommon/log"
	"github.com/lestrrat/go-server-starter/listener"
	"github.com/nmcclain/ldap"
	"github.com/tredoe/osutil/user/crypt"
)

type ldapServer struct {
	baseServer
	logger *log.Logger
}

type ldapHandler struct {
	backends model.Backends
	logger   *log.Logger
	config   *stns.Config
}

func (h ldapHandler) Bind(bindDN, rawPassword string, conn net.Conn) (ldap.LDAPResultCode, error) {
	bindDN = strings.ToLower(bindDN)
	baseDN := strings.ToLower(h.config.LDAP.BaseDN)

	if !strings.HasSuffix(bindDN, baseDN) {
		h.logger.Warn(fmt.Sprintf("Bind Error: BindDN %s not our BaseDN %s", bindDN, baseDN))
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

	users, err := h.backends.FindUserByName(userName)
	if err != nil {
		log.Warn(fmt.Sprintf("Bind Error: User %s not found.", userName))
		return ldap.LDAPResultInvalidCredentials, nil
	}
	user = users[userName].(*model.User)

	if groupName != "" {
		groups, err := h.backends.FindGroupByName(groupName)
		if err != nil {
			log.Warn(fmt.Sprintf("Bind Error: Group %s not found.", groupName))
			return ldap.LDAPResultInvalidCredentials, nil
		}

		group = groups[groupName].(*model.Group)

		if user.GroupID != group.ID {
			log.Warn(fmt.Sprintf("Bind Error: User %s primary group is not %s.", userName, groupName))
			return ldap.LDAPResultInvalidCredentials, nil
		}

	}
	c := crypt.NewFromHash(user.Password)
	if c.Verify(user.Password, []byte(rawPassword)) != nil {
		log.Warn(fmt.Sprintf("Bind Error: invalid credentials as %s from %s", bindDN, conn.RemoteAddr().String()))
		return ldap.LDAPResultInvalidCredentials, nil
	}
	return ldap.LDAPResultSuccess, nil

}

func (h ldapHandler) Search(bindDN string, searchReq ldap.SearchRequest, conn net.Conn) (ldap.ServerSearchResult, error) {
	bindDN = strings.ToLower(bindDN)
	baseDN := strings.ToLower(h.config.LDAP.BaseDN)
	searchBaseDN := strings.ToLower(searchReq.BaseDN)

	// validate the user is authenticated and has appropriate access
	if len(bindDN) < 1 {
		return ldap.ServerSearchResult{ResultCode: ldap.LDAPResultInsufficientAccessRights}, fmt.Errorf("Search Error: Anonymous BindDN not allowed %s", bindDN)
	}
	if !strings.HasSuffix(bindDN, baseDN) {
		return ldap.ServerSearchResult{ResultCode: ldap.LDAPResultInsufficientAccessRights}, fmt.Errorf("Search Error: BindDN %s not in our BaseDN %s", bindDN, baseDN)
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

	switch filterEntity {
	default:
		return ldap.ServerSearchResult{
			ResultCode: ldap.LDAPResultOperationsError,
		}, fmt.Errorf("Search Error: unhandled filter type: %s [%s]", filterEntity, searchReq.Filter)
	case "posixgroup":
		groups, err := h.backends.Groups()
		if err != nil {
			return ldap.ServerSearchResult{
				ResultCode: ldap.LDAPResultOperationsError,
			}, fmt.Errorf("Search Error: can't fetch groups: %s [%s]", filterEntity, searchReq.Filter)
		}

		users, err := h.backends.Users()
		if err != nil {
			return ldap.ServerSearchResult{
				ResultCode: ldap.LDAPResultOperationsError,
			}, fmt.Errorf("Search Error: can't fetch users: %s [%s]", filterEntity, searchReq.Filter)
		}

		for _, g := range groups {
			attrs := []*ldap.EntryAttribute{}
			attrs = append(attrs, &ldap.EntryAttribute{"cn", []string{g.GetName()}})
			attrs = append(attrs, &ldap.EntryAttribute{"description", []string{fmt.Sprintf("%s via LDAP", g.GetName())}})
			attrs = append(attrs, &ldap.EntryAttribute{"gidNumber", []string{fmt.Sprintf("%d", g.GetID())}})
			attrs = append(attrs, &ldap.EntryAttribute{"objectClass", []string{"posixGroup"}})
			attrs = append(attrs, &ldap.EntryAttribute{"memberUid", h.getGroupMemberNames(g.(*model.Group), users, baseDN)})
			dn := fmt.Sprintf("cn=%s,ou=groups,%s", g.GetName(), baseDN)

			entries = append(entries, &ldap.Entry{dn, attrs})
		}
	case "posixaccount", "":
		users, err := h.backends.Users()
		if err != nil {
			return ldap.ServerSearchResult{
				ResultCode: ldap.LDAPResultOperationsError,
			}, fmt.Errorf("Search Error: can't fetch users: %s [%s]", filterEntity, searchReq.Filter)
		}
		groups, err := h.backends.Groups()
		if err != nil {
			return ldap.ServerSearchResult{
				ResultCode: ldap.LDAPResultOperationsError,
			}, fmt.Errorf("Search Error: can't fetch group: %s [%s]", filterEntity, searchReq.Filter)
		}

		for _, u := range users {
			var group *model.Group
			memberOf := []int{}

			for _, g := range groups {
				// find primary group which the user belongs
				if g.GetID() == u.(*model.User).GroupID {
					group = g.(*model.Group)
				}

				// find other group which the user belongs
				for _, gm := range g.(*model.Group).Users {
					if gm == u.GetName() {
						memberOf = append(memberOf, g.GetID())
					}
				}
			}

			if group == nil {
				return ldap.ServerSearchResult{
					ResultCode: ldap.LDAPResultOperationsError,
				}, fmt.Errorf("Search Error: primary group id is required : %s [%s]", filterEntity, searchReq.Filter)
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
			attrs = append(attrs, &ldap.EntryAttribute{"memberOf", h.getGroupDNs(memberOf, groups, baseDN)})
			attrs = append(attrs, &ldap.EntryAttribute{"sshPublicKey", u.(*model.User).Keys})
			dn := fmt.Sprintf("cn=%s,ou=%s,%s", u.GetName(), group.Name, baseDN)
			entries = append(entries, &ldap.Entry{dn, attrs})
		}
	}

	return ldap.ServerSearchResult{entries, []string{}, []ldap.Control{}, ldap.LDAPResultSuccess}, nil
}

func newLDAPServer(confPath string) (*ldapServer, error) {
	conf, err := stns.NewConfig(confPath)
	if err != nil {
		return nil, err
	}

	s := &ldapServer{
		baseServer: baseServer{config: &conf},
		logger:     log.New("stns-ldap"),
	}
	return s, nil
}

func (s *ldapServer) Run() error {
	var backends model.Backends
	b, err := model.NewBackendTomlFile(s.config.Users, s.config.Groups)
	if err != nil {
		return err
	}
	backends = append(backends, b)

	err = s.loadModules(s.logger, &backends)
	if err != nil {
		return err
	}

	ld := ldap.NewServer()
	ld.EnforceLDAP = true
	h := ldapHandler{
		backends: backends,
		config:   s.config,
	}
	ld.BindFunc("", h)
	ld.SearchFunc("", h)

	if err := pidfile.Write(); err != nil {
		return err
	}
	defer func() {
		if err := os.Remove(pidfile.GetPidfilePath()); err != nil {
			log.Fatalf("Error removing %s: %s", pidfile.GetPidfilePath(), err)
		}
	}()

	p := strconv.Itoa(s.config.Port)
	lnstr := ":" + p
	if s.config.UseServerStarter {
		listeners, err := listener.ListenAll()
		if listeners == nil || err != nil {
			return err
		}
		lnstr = listeners[0].Addr().String()
	}

	log.Info("start ldap server")
	if err := ld.ListenAndServe(lnstr); err != nil {
		return err
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
