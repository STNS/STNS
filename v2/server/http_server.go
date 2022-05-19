package server

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	stdLog "log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/STNS/STNS/v2/api"
	"github.com/STNS/STNS/v2/middleware"
	"github.com/STNS/STNS/v2/model"
	"github.com/STNS/STNS/v2/stns"
	"github.com/facebookgo/pidfile"
	"github.com/labstack/echo"
	emiddleware "github.com/labstack/echo/middleware"

	"github.com/labstack/gommon/log"

	"github.com/lestrrat/go-server-starter/listener"
)

type httpServer struct {
	baseServer
}

func newHTTPServer(conf *stns.Config, backend model.Backend, logger *log.Logger) (*httpServer, error) {
	s := &httpServer{
		baseServer{
			config:  conf,
			backend: backend,
			logger:  logger,
		},
	}
	return s, nil
}

// Run サーバの起動
func (s *httpServer) Run() error {
	e := echo.New()
	e.Logger = s.logger
	e.StdLogger = stdLog.New(s.logger.Output(), "", stdLog.Ldate|stdLog.Ltime|stdLog.Llongfile)
	e.GET("/status", status)

	if err := pidfile.Write(); err != nil {
		return err
	}

	defer func() {
		if err := os.Remove(pidfile.GetPidfilePath()); err != nil {
			e.Logger.Errorf("Error removing %s: %s", pidfile.GetPidfilePath(), err)
		}
	}()

	e.Use(middleware.Backend(s.backend))
	e.Use(middleware.AddHeader(s.backend))

	e.Use(emiddleware.Recover())
	e.Use(emiddleware.LoggerWithConfig(emiddleware.LoggerConfig{
		Format: `{"time":"${time_rfc3339_nano}","remote_ip":"${remote_ip}","host":"${host}",` +
			`"method":"${method}","uri":"${uri}","status":${status}}` + "\n",
		Output: s.logger.Output(),
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
					if c.Path() == "/" || strings.HasSuffix(c.Path(), "/status") || len(os.Getenv("CI")) > 0 {
						return true
					}
					return false
				},
			}))
	}

	if s.config.TokenAuth != nil {
		e.Use(middleware.TokenAuthWithConfig(middleware.TokenAuthConfig{
			Skipper: func(c echo.Context) bool {
				if c.Path() == "/" || strings.HasSuffix(c.Path(), "/status") || len(os.Getenv("CI")) > 0 {
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

	if s.config.AllowIPs != nil {
		e.Use(middleware.IPFilterWithConfig(middleware.IPFilterConfig{
			AllowIPs: s.config.AllowIPs,
			Logger:   s.logger,
		}))
	}

	if s.config.UseServerStarter {
		listeners, err := listener.ListenAll()
		if listeners == nil || err != nil {
			return err
		}
		e.Listener = listeners[0]
	}
	v1 := e.Group("/v1")
	v1.GET("/status", status)
	api.UserEndpoints(v1)
	api.GroupEndpoints(v1)

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello! STNS!!1")
	})

	customServer := &http.Server{
		WriteTimeout: 1 * time.Minute,
	}

	if e.Listener == nil {
		p := strconv.Itoa(s.config.Port)
		customServer.Addr = ":" + p
		if os.Getenv("STNS_LISTEN") != "" {
			customServer.Addr = os.Getenv("STNS_LISTEN")
		}
	}

	// tls client authentication
	if s.config.TLS != nil {
		if _, err := os.Stat(s.config.TLS.CA); err == nil {
			ca, err := ioutil.ReadFile(s.config.TLS.CA)
			if err != nil {
				return err
			}
			caPool := x509.NewCertPool()
			caPool.AppendCertsFromPEM(ca)

			tlsConfig := &tls.Config{
				ClientCAs:              caPool,
				SessionTicketsDisabled: true,
				ClientAuth:             tls.RequireAndVerifyClientCert,
			}

			customServer.TLSConfig = tlsConfig
		}
	}

	if s.config.TLS != nil && s.config.TLS.Cert != "" && s.config.TLS.Key != "" {
		if customServer.TLSConfig == nil {
			customServer.TLSConfig = new(tls.Config)
		}
		cert, err := tls.LoadX509KeyPair(s.config.TLS.Cert, s.config.TLS.Key)
		if err != nil {
			return err
		}
		customServer.TLSConfig.Certificates = make([]tls.Certificate, 1)
		customServer.TLSConfig.Certificates[0] = cert
	}

	go func() {
		if err := e.StartServer(customServer); err != nil {
			e.Logger.Error("shutting down the server: %s", err)
			return
		}

	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		return err
	}
	return nil
}

func status(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}
