package api

import (
	"context"
	"errors"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"plugin"
	"strconv"
	"time"

	emiddleware "github.com/labstack/echo/middleware"

	"github.com/STNS/STNS/middleware"
	"github.com/STNS/STNS/model"
	"github.com/STNS/STNS/stns"
	"github.com/facebookgo/pidfile"
	"github.com/labstack/echo"

	"github.com/urfave/cli"

	// PostgreSQL driver
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/lestrrat/go-server-starter/listener"
)

type server struct {
	config *stns.Config
}

func newServer(confPath string) (*server, error) {
	conf, err := stns.NewConfig(confPath)
	if err != nil {
		return nil, err
	}

	s := &server{config: &conf}
	return s, nil
}
func status(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}

func (s *server) loadModules(logger echo.Logger, getterBackends model.GetterBackends) error {
	for _, v := range s.config.LoadModules {
		p, err := plugin.Open(filepath.Join(s.config.ModulePath, v))
		if err != nil {
			return err
		}

		name, err := p.Lookup("ModuleName")
		if err != nil {
			return err
		}

		b, err := p.Lookup("NewBackend" + name.(string))
		if err != nil {
			return err
		}

		backend, err := b.(func(*stns.Config) (model.Backend, error))(s.config)
		if err != nil {
			return err
		}
		getterBackends = append(getterBackends, backend)
		logger.Info("load modules %s", name.(string))
	}
	return nil
}

// Run サーバの起動
func (s *server) Run() error {
	var getterBackends model.GetterBackends
	e := echo.New()
	e.GET("/status", status)

	if err := pidfile.Write(); err != nil {
		return err
	}
	defer removePidFile(e)

	b, err := model.NewBackendTomlFile(s.config.Users, s.config.Groups)
	if err != nil {
		return err
	}
	getterBackends = append(getterBackends, b)

	err = s.loadModules(e.Logger, getterBackends)
	if err != nil {
		return err
	}

	if os.Getenv("STNS_LOG") != "" {
		f, err := os.OpenFile(os.Getenv("STNS_LOG"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			return errors.New("error opening file :" + err.Error())
		}
		e.Logger.SetOutput(f)
	}

	e.Use(middleware.GetterBackends(getterBackends))
	e.Use(middleware.AddHeader(getterBackends))
	e.Use(emiddleware.Recover())
	e.Use(emiddleware.LoggerWithConfig(emiddleware.LoggerConfig{
		Format: `{"time":"${time_rfc3339_nano}","remote_ip":"${remote_ip}","host":"${host}",` +
			`"method":"${method}","uri":"${uri}","status":${status}}` + "\n",
	}))

	if s.config.BasicAuth != nil {
		e.Use(emiddleware.BasicAuthWithConfig(
			emiddleware.BasicAuthConfig{
				Validator: func(username, password string, c echo.Context) (bool, error) {
					if username == s.config.BasicAuth.User && password == s.config.BasicAuth.Password {
						return true, nil
					}
					return false, nil
				},
				Skipper: func(c echo.Context) bool {
					if c.Path() == "/" || c.Path() == "/status" || len(os.Getenv("CI")) > 0 {
						return true
					}
					return false
				},
			}))
	}

	if s.config.TokenAuth != nil {
		e.Use(middleware.TokenAuthWithConfig(middleware.TokenAuthConfig{
			Skipper: func(c echo.Context) bool {

				if c.Path() == "/" || c.Path() == "/status" || len(os.Getenv("CI")) > 0 {
					return true
				}

				return false
			},
			Validator: func(token string) bool {
				for _, a := range s.config.TokenAuth.Tokens {
					if a == token {
						return true
					}
				}
				return false
			},
		}))
	}

	if s.config.UseServerStarter {
		listeners, err := listener.ListenAll()
		if listeners == nil || err != nil {
			return err
		}
		e.Listener = listeners[0]
	} else {
		p := strconv.Itoa(s.config.Port)
		l, err := net.Listen("tcp", ":"+p)
		if err != nil {
			return err
		}
		e.Listener = l
	}

	go func() {
		customServer := &http.Server{
			WriteTimeout: 1 * time.Minute,
		}
		if err := e.StartServer(customServer); err != nil {
			e.Logger.Info("shutting down the server")
		}
	}()

	v1 := e.Group("/v1")
	UserEndpoints(v1)
	GroupEndpoints(v1)

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello! STNS!!1")
	})

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		return err
	}
	return nil
}

func removePidFile(e *echo.Echo) {
	if err := os.Remove(pidfile.GetPidfilePath()); err != nil {
		e.Logger.Fatalf("Error removing %s: %s", pidfile.GetPidfilePath(), err)
	}
}

func LaunchServer(c *cli.Context) error {

	pidfile.SetPidfilePath(os.Getenv("STNS_PID"))
	serv, err := newServer(os.Getenv("STNS_CONFIG"))
	if err != nil {
		return errors.New("server init:" + err.Error())
	}
	return serv.Run()
}
