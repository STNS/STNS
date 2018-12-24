package api

import (
	"fmt"
	"os"

	"github.com/STNS/STNS/model"
	"github.com/STNS/STNS/stns"
	"github.com/urfave/cli"
)

func CheckConfig(c *cli.Context) error {
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
