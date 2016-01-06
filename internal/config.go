package stns

import (
	"fmt"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Port    int
	Include string
	Users   map[string]*Attr
	Groups  map[string]*Attr
}

type UserAttr struct {
	Group_Id  int
	Directory string
	Shell     string
	Gecos     string
	Keys      []string
}
type GroupAttr struct {
	Users []string
}
type Attr struct {
	Id   int
	Name string
	*UserAttr
	*GroupAttr
}

var AllConfig *Config

func LoadConfig(configFile string) error {
	var config Config
	defaultConfig(&config)

	_, err := toml.DecodeFile(configFile, &config)
	if err != nil {
		return err
	}

	if config.Include != "" {
		if err := includeConfigFile(&config, config.Include); err != nil {
			return err
		}
	}

	AllConfig = &config
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
		config.Users = nil
		config.Groups = nil

		_, err := toml.DecodeFile(file, &config)
		if err != nil {
			return fmt.Errorf("while loading included config file %s: %s", file, err)
		}
		config.Users = merge(config.Users, userSaved)
		config.Groups = merge(config.Groups, groupSaved)
	}
	return nil
}

func merge(m1 map[string]*Attr, m2 map[string]*Attr) map[string]*Attr {
	m := map[string]*Attr{}

	for i, v := range m1 {
		m[i] = v
	}
	for i, v := range m2 {
		m[i] = v
	}
	return m
}
