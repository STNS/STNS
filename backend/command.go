package backend

import (
	"errors"
	"flag"
	"fmt"

	"github.com/STNS/STNS/stns"
)

var singletonB Backend

type Backend interface {
	Migrate() error
	Delete() error
}

func getInstance(config *stns.Config) Backend {
	if singletonB == nil {
		switch config.Backend.Driver {
		case "mysql":
			singletonB = &Mysql{config}
		}
	}
	return singletonB
}

func SubCommandRun(config *stns.Config) error {
	if len(flag.Args()) > 1 {
		b := getInstance(config)
		if b == nil {
			return errors.New("unknown backend driver:" + config.Backend.Driver)
		}

		switch flag.Args()[1] {
		case "init":
			if err := b.Migrate(); err != nil {
				return err
			} else {
				fmt.Println("backend driver " + config.Backend.Driver + " init successful")
				return nil
			}
		case "delete":
			if err := b.Delete(); err != nil {
				return err
			} else {
				fmt.Println("backend driver " + config.Backend.Driver + " delete successful")
				return nil
			}
		}
	}
	return errors.New(usageTemplate)
}

var usageTemplate = `
Usage:
	stns backend [arguments]

The commands are:
	init	initialize backend
	delete  remove all of the information
`
