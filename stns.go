package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/ant0ine/go-json-rest/rest"
	"github.com/pyama86/STNS/api"
	"github.com/pyama86/STNS/config"
)

func createPidFile(pidFile string) error {
	if pidString, err := ioutil.ReadFile(pidFile); err == nil {
		pid, err := strconv.Atoi(string(pidString))
		if err == nil {
			if _, err := os.Stat(fmt.Sprintf("/proc/%d/", pid)); err == nil {
				return fmt.Errorf("pid file found, ensure stns  is not running or delete %s", pidFile)
			}
		}
	}

	file, err := os.Create(pidFile)
	if err != nil {
		return err
	}

	defer file.Close()

	_, err = fmt.Fprintf(file, "%d", os.Getpid())
	return err
}

func removePidFile(pidFile string) {
	if err := os.Remove(pidFile); err != nil {
		log.Fatal("Error removing %s: %s", pidFile, err)
	}
}

func startServer(pidFile string, configFile string) {
	if err := config.Load(configFile); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	if err := createPidFile(pidFile); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	defer removePidFile(pidFile)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, os.Signal(syscall.SIGTERM))
	go func() {
		sig := <-c
		log.Printf("Received signal '%v', exiting\n", sig)
		removePidFile(pidFile)
		os.Exit(0)
	}()

	server := rest.NewApi()
	server.Use(rest.DefaultDevStack...)
	router, err := rest.MakeRouter(
		rest.Get("/:resource_name/:column/:value", api.GetAttribute),
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
