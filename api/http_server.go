package api

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/STNS/STNS/middleware"
	"github.com/STNS/STNS/model"
	"github.com/STNS/STNS/stns"
	"github.com/facebookgo/pidfile"
	"github.com/labstack/echo"
	emiddleware "github.com/labstack/echo/middleware"

	"github.com/labstack/gommon/log"

	"github.com/lestrrat/go-server-starter/listener"
)

type httpServer struct {
	baseServer
}

func newHTTPServer(confPath string) (*httpServer, error) {
	conf, err := stns.NewConfig(confPath)
	if err != nil {
		return nil, err
	}

	s := &httpServer{
		baseServer{config: &conf},
	}
	return s, nil
}

func status(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}

// Run サーバの起動
func (s *httpServer) Run() error {
	var backends model.Backends
	e := echo.New()
	if os.Getenv("STNS_LOG") != "" {
		f, err := os.OpenFile(os.Getenv("STNS_LOG"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			return errors.New("error opening file :" + err.Error())
		}
		e.Logger.SetOutput(f)
	} else {
		e.Logger.SetLevel(log.DEBUG)
	}
	e.GET("/status", status)

	if err := pidfile.Write(); err != nil {
		return err
	}

	defer func() {
		if err := os.Remove(pidfile.GetPidfilePath()); err != nil {
			e.Logger.Fatalf("Error removing %s: %s", pidfile.GetPidfilePath(), err)
		}
	}()

	b, err := model.NewBackendTomlFile(s.config.Users, s.config.Groups)
	if err != nil {
		return err
	}
	backends = append(backends, b)

	err = s.loadModules(e.Logger.(*log.Logger), &backends)
	if err != nil {
		return err
	}

	e.Use(middleware.Backends(backends))
	e.Use(middleware.AddHeader(backends))

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
	}

	go func() {
		customServer := &http.Server{
			WriteTimeout: 1 * time.Minute,
		}
		if e.Listener == nil {
			p := strconv.Itoa(s.config.Port)
			customServer.Addr = ":" + p
		}

		// tls client authentication
		if s.config.TLS != nil {
			if _, err := os.Stat(s.config.TLS.CA); err == nil {
				ca, err := ioutil.ReadFile(s.config.TLS.CA)
				if err != nil {
					e.Logger.Fatal(err)
				}
				caPool := x509.NewCertPool()
				caPool.AppendCertsFromPEM(ca)

				tlsConfig := &tls.Config{
					ClientCAs:              caPool,
					SessionTicketsDisabled: true,
					ClientAuth:             tls.RequireAndVerifyClientCert,
				}

				tlsConfig.BuildNameToCertificate()
				customServer.TLSConfig = tlsConfig
			}
		}

		if s.config.TLS != nil && s.config.TLS.Cert != "" && s.config.TLS.Key != "" {
			if customServer.TLSConfig == nil {
				customServer.TLSConfig = new(tls.Config)
			}
			customServer.TLSConfig.Certificates = make([]tls.Certificate, 1)
			customServer.TLSConfig.Certificates[0], err = tls.LoadX509KeyPair(s.config.TLS.Cert, s.config.TLS.Key)
			if err != nil {
				e.Logger.Fatal(err)
			}

		}

		if err := e.StartServer(customServer); err != nil {
			e.Logger.Fatal("shutting down the server")
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
