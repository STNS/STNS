package api

import (
	"os"

	"github.com/facebookgo/pidfile"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

// LaunchServer
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
	return nil
}
