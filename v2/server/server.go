package server

import (
	"errors"
	"os"
	"path/filepath"
	"plugin"
	"strings"

	"github.com/STNS/STNS/v2/model"
	"github.com/STNS/STNS/v2/stns"
	"github.com/facebookgo/pidfile"
	"github.com/iancoleman/strcase"

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
	logger := log.New("stns")

	// set conf
	conf, err := stns.NewConfig(os.Getenv("STNS_CONFIG"))
	if err != nil {
		return err
	}

	var backend model.Backend
	// set backend
	backend, err = model.NewBackendTomlFile(conf.Users, conf.Groups)
	if err != nil {
		return err
	}

	b, err := loadBackendModule(logger, &conf)
	if err != nil {
		return err
	}

	if b != nil {
		backend = b
	}

	if conf.Redis != nil && conf.Redis.Host != "" {
		r, err := model.NewBackendRedis(backend, logger, conf.Redis.Host, conf.Redis.Port, conf.Redis.Password, conf.Redis.TTL, conf.Redis.DB)
		if err != nil {
			return err
		}
		backend = r
	}

	// set log output
	if os.Getenv("STNS_LOG") != "" {
		f, err := os.OpenFile(os.Getenv("STNS_LOG"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			return errors.New("error opening file :" + err.Error())
		}
		logger.SetOutput(f)
	} else {
		logger.SetLevel(log.DEBUG)
	}

	switch os.Getenv("STNS_PROTOCOL") {
	case "ldap":
		serv, err = newLDAPServer(&conf, backend, logger)
	default:
		serv, err = newHTTPServer(&conf, backend, logger)
	}

	if err != nil {
		return errors.New("server init:" + err.Error())
	}
	return serv.Run()
}

type baseServer struct {
	config  *stns.Config
	logger  *log.Logger
	backend model.Backend
}

func loadBackendModule(logger *log.Logger, conf *stns.Config) (model.Backend, error) {
	if conf.LoadModule == "" {
		return nil, nil
	}
	p, err := plugin.Open(filepath.Join(conf.ModulePath, conf.LoadModule))
	if err != nil {
		return nil, err
	}

	name := conf.LoadModule
	name = strings.Replace(name, ".so", "", 1)
	name = strcase.ToCamel(strings.Replace(name, "mod_stns_", "", 1))
	b, err := p.Lookup("NewBackend" + name)
	if err != nil {
		return nil, err
	}

	backend, err := b.(func(*stns.Config) (model.Backend, error))(conf)
	if err != nil {
		return nil, err
	}
	logger.Infof("load module %s", name)
	return backend, err
}
