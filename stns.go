package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/ant0ine/go-json-rest/rest"
	"github.com/pyama86/STNS/api"
	"github.com/pyama86/STNS/config"
	"github.com/pyama86/STNS/pid"
)

func startServer(pidFile string, configFile string) {
	if err := config.Load(configFile); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	if err := pid.CreatePidFile(pidFile); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	defer pid.RemovePidFile(pidFile)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, os.Signal(syscall.SIGTERM))
	go func() {
		sig := <-c
		log.Printf("Received signal '%v', exiting\n", sig)
		pid.RemovePidFile(pidFile)
		os.Exit(0)
	}()

	server := rest.NewApi()
	server.Use(rest.DefaultDevStack...)
	router, err := rest.MakeRouter(
		rest.Get("/:resource_name/list", api.GetList),
		rest.Get("/:resource_name/:column/:value", api.Get),
	)
	if err != nil {
		log.Fatal(err)
	}

	server.SetApp(router)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(config.All.Port), server.MakeHandler()))
}

func main() {

	f, err := os.OpenFile("/var/log/stns.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)

	if err != nil {
		log.Fatal("error opening file :", err.Error())
	}

	log.SetOutput(f)

	configFile := flag.String("conf", "/etc/stns/stns.conf", "config file path")
	pidFile := flag.String("pidfile", "/var/run/stns.pid", "File containing process PID")
	flag.Parse()

	startServer(*pidFile, *configFile)
}
