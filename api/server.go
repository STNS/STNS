package api

import (
	"errors"
	"os"
	"path/filepath"
	"plugin"

	"github.com/STNS/STNS/model"
	"github.com/STNS/STNS/stns"
	"github.com/facebookgo/pidfile"

	"github.com/labstack/gommon/log"
	_ "github.com/tredoe/osutil/user/crypt/md5_crypt"
	_ "github.com/tredoe/osutil/user/crypt/sha256_crypt"
	_ "github.com/tredoe/osutil/user/crypt/sha512_crypt"
	"github.com/urfave/cli"
)

type server interface {
	Run() error
}

func LaunchServer(c *cli.Context) error {
	var serv server
	var err error
	pidfile.SetPidfilePath(os.Getenv("STNS_PID"))

	switch os.Getenv("STNS_PROTOCOL") {
	case "ldap":
		serv, err = newLDAPServer(os.Getenv("STNS_CONFIG"))
	default:
		serv, err = newHTTPServer(os.Getenv("STNS_CONFIG"))
	}

	if err != nil {
		return errors.New("server init:" + err.Error())
	}
	return serv.Run()
}

type baseServer struct {
	config *stns.Config
}

func (s *baseServer) loadModules(logger *log.Logger, backends *model.Backends) error {
	for _, v := range s.config.LoadModules {
		p, err := plugin.Open(filepath.Join(s.config.ModulePath, v))
		if err != nil {
			return err
		}

		n, err := p.Lookup("ModuleName")
		if err != nil {
			return err
		}
		name := *(n.(*string))
		b, err := p.Lookup("NewBackend" + name)
		if err != nil {
			return err
		}

		backend, err := b.(func(*stns.Config) (model.Backend, error))(s.config)
		if err != nil {
			return err
		}
		*backends = append(*backends, backend)
		logger.Infof("load modules %s", name)
	}
	return nil
}
