package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/STNS/STNS/settings"
	"github.com/STNS/STNS/stns"
)

func main() {
	configFile := flag.String("conf", "/etc/stns/stns.conf", "config file path")
	pidFile := flag.String("pidfile", "/var/run/stns.pid", "File containing process PID")
	configCheck := flag.Bool("check-conf", false, "config check flag")
	version := flag.Bool("version", false, "Print version")
	verbose := flag.Bool("verbose", false, "verbose log")
	flag.Parse()

	if *version {
		fmt.Println("STNS version " + settings.VERSION)
		os.Exit(0)
	}

	config, err := stns.LoadConfig(*configFile)
	if err != nil {
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

	server := stns.Create(config, *configFile, *pidFile, *verbose)
	server.Start()
}
