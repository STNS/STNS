package backend

import (
	"errors"
	"flag"
	"fmt"

	"github.com/STNS/STNS/stns"
)

func NewBackend(config *stns.Config) Backend {
	switch config.Backend.Driver {
	case "mysql":
		return &Mysql{config}
	}
	return nil
}

type Backend interface {
	Migrate() error
}

func SubCommandRun(config *stns.Config) error {
	if len(flag.Args()) > 1 {
		switch flag.Args()[1] {
		case "init":
			b := NewBackend(config)
			if b == nil {
				return errors.New("unknown backend driver:" + config.Backend.Driver)
			}

			if err := b.Migrate(); err != nil {
				return err
			} else {
				fmt.Println("backend driver " + config.Backend.Driver + " init successful")
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

`
