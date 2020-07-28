package main

import (
	"fmt"
	"os"

	"github.com/STNS/STNS/v2/model"
	"github.com/STNS/STNS/v2/server"
	"github.com/STNS/STNS/v2/stns"
	"github.com/urfave/cli"
)

var (
	version   string
	revision  string
	goversion string
	builddate string
	builduser string
)

var flags = []cli.Flag{
	cli.StringFlag{
		Name:   "logfile",
		Value:  "",
		Usage:  "Log file path",
		EnvVar: "STNS_LOG",
	},
	cli.StringFlag{
		Name:   "config",
		Value:  "/etc/stns/server/stns.conf",
		Usage:  "Server config",
		EnvVar: "STNS_CONFIG",
	},
	cli.StringFlag{
		Name:   "pidfile",
		Value:  "/var/run/stns.pid",
		Usage:  "pid file path",
		EnvVar: "STNS_PID",
	},
	cli.StringFlag{
		Name:   "protocol",
		Value:  "http",
		Usage:  "interface protocol",
		EnvVar: "STNS_PROTOCOL",
	},
	cli.StringFlag{
		Name:   "listen",
		Value:  "",
		Usage:  "listern addrand port(xxx.xxx.xxx.xxx:yyy)",
		EnvVar: "STNS_LISTEN",
	},
}

var commands = []cli.Command{
	{
		Name:    "server",
		Aliases: []string{"s"},
		Usage:   "Launch core api server",
		Action:  server.LaunchServer,
	},
	{
		Name:    "checkconf",
		Aliases: []string{"c"},
		Usage:   "Check Config",
		Action:  checkConfig,
	},
}

func printVersion(c *cli.Context) {
	fmt.Printf("stns version: %s (%s)\n", version, revision)
	fmt.Printf("build at %s (with %s) by %s\n", builddate, goversion, builduser)
}

func appBefore(c *cli.Context) error {
	// I want to quit this implementation
	if c.GlobalString("logfile") != "" {
		os.Setenv("STNS_LOG", c.GlobalString("logfile"))
	}

	if c.GlobalString("config") != "" {
		os.Setenv("STNS_CONFIG", c.GlobalString("config"))
	}

	if c.GlobalString("pidfile") != "" {
		os.Setenv("STNS_PID", c.GlobalString("pidfile"))
	}

	if c.GlobalString("protocol") != "" {
		os.Setenv("STNS_PROTOCOL", c.GlobalString("protocol"))
	}

	if c.GlobalString("listen") != "" {
		os.Setenv("STNS_LISTEN", c.GlobalString("listen"))
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
		app.Action = server.LaunchServer
	} else {
		app.Commands = commands
	}

	app.Before = appBefore

	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}

func checkConfig(c *cli.Context) error {
	conf, err := stns.NewConfig(os.Getenv("STNS_CONFIG"))
	if err != nil {
		return err
	}
	_, err = model.NewBackendTomlFile(conf.Users, conf.Groups)
	if err == nil {
		fmt.Println("config is good!!1")
	}
	return err
}
