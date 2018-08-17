package main

import (
	"fmt"
	"os"

	"github.com/STNS/STNS/api"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var (
	version   string
	revision  string
	goversion string
	builddate string
	builduser string
)

func init() {
	formatter := new(log.JSONFormatter)
	formatter.TimestampFormat = "2006-01-02 15:04:05"
	log.SetFormatter(formatter)
}

var flags = []cli.Flag{
	cli.StringFlag{
		Name:   "logfile",
		Value:  "/var/log/stns.log",
		Usage:  "Log file path",
		EnvVar: "STNS_LOG",
	},
	cli.StringFlag{
		Name:   "config",
		Value:  "/etc/stns/stns.conf",
		Usage:  "Server config",
		EnvVar: "STNS_CONFIG",
	},
	cli.StringFlag{
		Name:   "pidfile",
		Value:  "/var/run/stns.pid",
		Usage:  "pid file path",
		EnvVar: "STNS_PID",
	},
}

var commands = []cli.Command{
	{
		Name:    "server",
		Aliases: []string{"s"},
		Usage:   "Launch core api server",
		Action:  api.LaunchServer,
	},
}

func printVersion(c *cli.Context) {
	fmt.Printf("stns version: %s (%s)\n", version, revision)
	fmt.Printf("build at %s (with %s) by %s\n", builddate, goversion, builduser)
}

func appBefore(c *cli.Context) error {
	if c.GlobalString("logfile") != "" {
		os.Setenv("STNS_LOG", c.GlobalString("logfile"))
	}

	if c.GlobalString("config") != "" {
		os.Setenv("STNS_CONFIG", c.GlobalString("config"))
	}

	if c.GlobalString("pidfile") != "" {
		os.Setenv("STNS_PID", c.GlobalString("pidfile"))
	}
	return nil
}

func main() {
	cli.VersionPrinter = printVersion

	app := cli.NewApp()
	app.Name = "stns"
	app.Usage = "Simple Toml Name Service"
	app.Flags = flags

	if len(os.Args) <= 1 {
		app.Action = api.LaunchServer
	} else {
		app.Commands = commands
	}

	app.Before = appBefore

	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}
