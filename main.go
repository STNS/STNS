package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/STNS/STNS/api"
	"github.com/STNS/STNS/config"
	"github.com/STNS/STNS/pid"
	"github.com/ant0ine/go-json-rest/rest"
)

func getHandler() http.Handler {
	server := rest.NewApi()
	server.Use(rest.DefaultDevStack...)
	// using basic auth
	if config.All.User != "" && config.All.Password != "" {

		var basicAuthMiddleware = &rest.AuthBasicMiddleware{
			Realm: "stns",
			Authenticator: func(user string, password string) bool {
				return user == config.All.User && password == config.All.Password
			},
		}
		server.Use(&rest.IfMiddleware{
			Condition: func(request *rest.Request) bool {
				return request.URL.Path != "/healthcheck"
			},
			IfTrue: basicAuthMiddleware,
		})
	}

	router, err := rest.MakeRouter(
		rest.Get("/:resource_name/list", api.GetList),
		rest.Get("/:resource_name/:column/:value", api.Get),
		rest.Get("/healthcheck", api.HealthChech),
	)
	if err != nil {
		log.Fatal(err)
	}

	server.SetApp(router)
	return server.MakeHandler()
}

func main() {
	configFile := flag.String("conf", "/etc/stns/stns.conf", "config file path")
	pidFile := flag.String("pidfile", "/var/run/stns.pid", "File containing process PID")
	configCheck := flag.Bool("check-conf", false, "config check flag")

	flag.Parse()

	if err := config.Load(configFile); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if *configCheck {
		fmt.Println("check config success!")
		os.Exit(0)
	}

	// set log
	f, err := os.OpenFile("/var/log/stns.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal("error opening file :", err.Error())
	}
	log.SetOutput(f)

	// wait signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, os.Signal(syscall.SIGTERM), os.Signal(syscall.SIGUSR1))
	go func() {
		for {
		Loop:
			sig := <-c
			if sig == os.Signal(syscall.SIGUSR1) {
				log.Print("Received signal reload config")
				if err := config.Load(configFile); err != nil {
					log.Fatal(err)
					os.Exit(1)
				}
				log.Printf("Complete reload config\n")
				goto Loop
			} else {
				log.Printf("Received signal '%v', exiting\n", sig)
				pid.RemovePidFile(pidFile)
				os.Exit(0)
			}
		}
	}()

	// make pid
	if err := pid.CreatePidFile(pidFile); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	defer pid.RemovePidFile(pidFile)
	log.Printf("Start Server pid:%d", os.Getpid())

	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(config.All.Port), getHandler()))
}
