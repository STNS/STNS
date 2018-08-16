package api

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/STNS/STNS/middleware"
	"github.com/STNS/STNS/stns"
	"github.com/facebookgo/pidfile"
	"github.com/labstack/echo"

	"github.com/sirupsen/logrus"
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

// Run サーバの起動
func (s *server) Run() error {
	e := echo.New()
	e.GET("/status", status)

	if err := pidfile.Write(); err != nil {
		return err
	}
	defer removePidFile()

	e.Use(middleware.Config(s.config))

	if s.config.UseServerStarter {
		listeners, err := listener.ListenAll()
		if listeners == nil || err != nil {
			return err
		}
		e.Listener = listeners[0]
	} else {
		l, err := net.Listen("tcp", ":8050")
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
			logrus.Info("shutting down the server")
		}
	}()

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

func removePidFile() {
	if err := os.Remove(pidfile.GetPidfilePath()); err != nil {
		logrus.Fatalf("Error removing %s: %s", pidfile.GetPidfilePath(), err)
	}
}

func LaunchServer(c *cli.Context) error {
	logrus.SetLevel(logrus.WarnLevel)
	if os.Getenv("STNS_LOG") != "" {
		f, err := os.OpenFile(os.Getenv("STNS_LOG"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			logrus.Fatal("error opening file :", err.Error())
		}
		logrus.SetOutput(f)
	}

	pidfile.SetPidfilePath(c.GlobalString("pidfile"))
	serv, err := newServer(os.Getenv("STNS_CONFIG"))
	if err != nil {
		return err
	}
	return serv.Run()
}
