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
	verbose := flag.Bool("verbose", false, "verbose log")
	logFile := flag.String("logfile", "/var/log/stns.log", "log file path")
	flag.Parse()

	config, err := stns.LoadConfig(*configFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if len(flag.Args()) > 0 {
		switch flag.Args()[0] {
		case "version":
			fmt.Println("STNS version " + settings.VERSION)
			os.Exit(0)
		case "check-conf":
			fmt.Println("configuration file " + *configFile + " test is successful")
			os.Exit(0)
		default:
			fmt.Fprintln(os.Stderr, "unknown command:"+flag.Args()[0])
			os.Exit(1)
		}
	}

	// set log
	if *logFile != "" {
		f, err := os.OpenFile(*logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			log.Fatal("error opening file :", err.Error())
		}
		log.SetOutput(f)
	}

	server := stns.Create(config, *configFile, *pidFile, *verbose)
	server.Start()
}
