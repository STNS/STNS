package config

import (
	"fmt"
	"path/filepath"
	"sync"

	"github.com/BurntSushi/toml"
	"github.com/STNS/STNS/attribute"
)

type Config struct {
	Port     int    `toml:"port"`
	Include  string `toml:"include"`
	User     string `toml:"user"`
	Password string `toml:"password"`
	Users    attribute.UserGroups
	Groups   attribute.UserGroups
}

var (
	All        *Config
	configLock = new(sync.RWMutex)
)

func Load(configFile *string) error {
	var config Config
	defaultConfig(&config)

	_, err := toml.DecodeFile(*configFile, &config)
	if err != nil {
		return err
	}

	if config.Include != "" {
		if err := includeConfigFile(&config, config.Include); err != nil {
			return err
		}
	}
	configLock.Lock()
	All = &config
	configLock.Unlock()
	return nil
}
func defaultConfig(config *Config) {
	config.Port = 1104
}

func includeConfigFile(config *Config, include string) error {
	files, err := filepath.Glob(include)
	if err != nil {
		return err
	}

	for _, file := range files {
		userSaved := config.Users
		groupSaved := config.Groups
		config.Users = attribute.UserGroups{}
		config.Groups = attribute.UserGroups{}

		_, err := toml.DecodeFile(file, &config)
		if err != nil {
			return fmt.Errorf("while loading included config file %s: %s", file, err)
		}
		config.Users.Merge(userSaved)
		config.Groups.Merge(groupSaved)
	}
	return nil
}
