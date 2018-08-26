package stns

import (
	"fmt"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/STNS/STNS/model"
)

func NewConfig(confPath string) (Config, error) {
	var conf Config
	defaultConfig(&conf)

	if _, err := toml.DecodeFile(confPath, &conf); err != nil {
		return conf, err
	}

	if conf.Include != "" {
		if err := includeConfigFile(&conf, conf.Include); err != nil {
			return Config{}, err
		}
	}

	return conf, nil
}

type Config struct {
	Port             int `toml:"port"`
	UseServerStarter bool
	Users            *model.Users
	Groups           *model.Groups
	Include          string `toml:"include"`
}

func defaultConfig(c *Config) {
	c.Port = 1104
}

func includeConfigFile(config *Config, include string) error {
	files, err := filepath.Glob(include)
	if err != nil {
		return err
	}

	for _, file := range files {
		_, err := toml.DecodeFile(file, &config)
		if err != nil {
			return fmt.Errorf("while loading included config file %s: %s", file, err)
		}
	}
	return nil
}
