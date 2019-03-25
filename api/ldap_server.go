package api

import (
	"log"
	"net"
	"os"
	"strconv"

	"github.com/STNS/STNS/stns"
	"github.com/facebookgo/pidfile"
	"github.com/lestrrat/go-server-starter/listener"
	"github.com/nmcclain/ldap"
)

type ldapServer struct {
	baseServer
}

type ldapHandler struct {
}

func (h ldapHandler) Bind(bindDN, bindSimplePw string, conn net.Conn) (ldap.LDAPResultCode, error) {
	return ldap.LDAPResultInvalidCredentials, nil
}

func (h ldapHandler) Search(boundDN string, searchReq ldap.SearchRequest, conn net.Conn) (ldap.ServerSearchResult, error) {

	entries := []*ldap.Entry{
		&ldap.Entry{"cn=ned," + searchReq.BaseDN, []*ldap.EntryAttribute{
			&ldap.EntryAttribute{"cn", []string{"ned"}},
			&ldap.EntryAttribute{"uidNumber", []string{"5000"}},
			&ldap.EntryAttribute{"accountStatus", []string{"active"}},
			&ldap.EntryAttribute{"uid", []string{"ned"}},
			&ldap.EntryAttribute{"description", []string{"ned"}},
			&ldap.EntryAttribute{"objectClass", []string{"posixAccount"}},
		}},
		&ldap.Entry{"cn=trent," + searchReq.BaseDN, []*ldap.EntryAttribute{
			&ldap.EntryAttribute{"cn", []string{"trent"}},
			&ldap.EntryAttribute{"uidNumber", []string{"5005"}},
			&ldap.EntryAttribute{"accountStatus", []string{"active"}},
			&ldap.EntryAttribute{"uid", []string{"trent"}},
			&ldap.EntryAttribute{"description", []string{"trent"}},
			&ldap.EntryAttribute{"objectClass", []string{"posixAccount"}},
		}},
	}
	return ldap.ServerSearchResult{entries, []string{}, []ldap.Control{}, ldap.LDAPResultSuccess}, nil
}

func newLDAPServer(confPath string) (*ldapServer, error) {
	conf, err := stns.NewConfig(confPath)
	if err != nil {
		return nil, err
	}

	s := &ldapServer{baseServer{config: &conf}}
	return s, nil
}

func (s *ldapServer) Run() error {

	ld := ldap.NewServer()
	ld.BindFunc("", ldapHandler{})
	ld.SearchFunc("", ldapHandler{})

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

	if err := ld.ListenAndServe(lnstr); err != nil {
		return err
	}

	return nil
}
